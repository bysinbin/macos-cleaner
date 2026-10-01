package cleaner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

// AppComponent represents an individual file or directory belonging to an application
type AppComponent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"` // "bundle", "support", "cache", "preferences", "state", "container", "storage"
	Size    uint64 `json:"size"`
	SizeStr string `json:"sizeStr"`
}

// InstalledApp represents an installed macOS application with its total footprint
type InstalledApp struct {
	ID          string         `json:"id"` // Absolute path to .app bundle
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	BundleID    string         `json:"bundleId"`
	Path        string         `json:"path"`
	AppSize     uint64         `json:"appSize"`
	AppSizeStr  string         `json:"appSizeStr"`
	DataSize    uint64         `json:"dataSize"`
	DataSizeStr string         `json:"dataSizeStr"`
	TotalSize   uint64         `json:"totalSize"`
	TotalSizeStr string        `json:"totalSizeStr"`
	IsSystem    bool           `json:"isSystem"`
	Components  []AppComponent `json:"components"`
}

// UninstallerScanResult wraps the full installed application scan
type UninstallerScanResult struct {
	Apps         []InstalledApp `json:"apps"`
	TotalAppSize uint64         `json:"totalAppSize"`
	TotalAppSizeStr string      `json:"totalAppSizeStr"`
	TotalDataSize uint64        `json:"totalDataSize"`
	TotalDataSizeStr string     `json:"totalDataSizeStr"`
	TotalAppsCount int          `json:"totalAppsCount"`
}

// ScanInstalledApplications scans for applications and calculates their total footprint
func ScanInstalledApplications(ctx context.Context) (*UninstallerScanResult, error) {
	scanner := NewScanner(ctx)
	home, _ := os.UserHomeDir()

	roots := []string{
		"/Applications",
		"/System/Applications",
		filepath.Join(home, "Applications"),
	}

	bundleIDRegex := regexp.MustCompile(`(?s)<key>CFBundleIdentifier</key>\s*<string>([^<]+)</string>`)
	versionRegex := regexp.MustCompile(`(?s)<key>CFBundleShortVersionString</key>\s*<string>([^<]+)</string>`)
	nameRegex := regexp.MustCompile(`(?s)<key>CFBundleName</key>\s*<string>([^<]+)</string>`)
	fallbackBundleRegex := regexp.MustCompile(`(?:com|org|io|net|de|app)\.[a-zA-Z0-9_-]+(?:\.[a-zA-Z0-9_-]+)+`)

	var appPaths []string
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			rel, _ := filepath.Rel(root, p)
			depth := strings.Count(rel, string(filepath.Separator))
			if depth > 4 {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() && strings.HasSuffix(d.Name(), ".app") {
				appPaths = append(appPaths, p)
				return filepath.SkipDir
			}
			return nil
		})
	}

	var apps []InstalledApp
	var mu sync.Mutex
	var wg sync.WaitGroup

	var totalAppBytes uint64
	var totalDataBytes uint64

	// Concurrency limiter
	sem := make(chan struct{}, 8)

	for _, appPath := range appPaths {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			base := strings.TrimSuffix(filepath.Base(p), ".app")
			bundleID := ""
			version := ""
			displayName := base

			// Read Info.plist
			plistPath := filepath.Join(p, "Contents", "Info.plist")
			if data, err := os.ReadFile(plistPath); err == nil {
				if len(data) > 128*1024 {
					data = data[:128*1024]
				}
				if m := bundleIDRegex.FindSubmatch(data); len(m) > 1 {
					bundleID = string(m[1])
				} else if m := fallbackBundleRegex.Find(data); len(m) > 0 {
					bundleID = string(m)
				}
				if m := versionRegex.FindSubmatch(data); len(m) > 1 {
					version = string(m[1])
				}
				if m := nameRegex.FindSubmatch(data); len(m) > 1 && len(m[1]) > 0 {
					displayName = string(m[1])
				}
			}

			// Calculate App Bundle Size
			appSz, _, _ := scanner.CalculateDirSize(p)
			if appSz == 0 {
				return
			}

			isSystem := strings.HasPrefix(p, "/System")

			var components []AppComponent
			components = append(components, AppComponent{
				Name:    "Uygulama Paketi (.app)",
				Path:    p,
				Type:    "bundle",
				Size:    appSz,
				SizeStr: FormatBytes(appSz),
			})

			var dataSz uint64

			// Find associated user data folders if bundleID or name is known
			candidates := []struct {
				relDir  string
				typ     string
				nameFmt string
			}{
				{"Library/Application Support/" + base, "support", "Uygulama Destek Verileri"},
				{"Library/Caches/" + bundleID, "cache", "Kullanıcı Önbelleği"},
				{"Library/Preferences/" + bundleID + ".plist", "preferences", "Kullanıcı Tercihleri (.plist)"},
				{"Library/Saved Application State/" + bundleID + ".savedState", "state", "Oturum Durumu (Saved State)"},
				{"Library/Containers/" + bundleID, "container", "Sandbox Konteyneri"},
				{"Library/WebKit/" + bundleID, "storage", "WebKit Web Verileri"},
				{"Library/HTTPStorages/" + bundleID, "storage", "HTTP Ağ Depolaması"},
				{"Library/Logs/" + base, "support", "Uygulama Günlükleri (Logs)"},
			}

			if bundleID != "" && bundleID != base {
				candidates = append(candidates, struct {
					relDir  string
					typ     string
					nameFmt string
				}{"Library/Application Support/" + bundleID, "support", "Uygulama Destek Verileri"})
			}

			for _, c := range candidates {
				if c.relDir == "" || strings.HasSuffix(c.relDir, "/") {
					continue
				}
				full := filepath.Join(home, c.relDir)
				fi, err := os.Stat(full)
				if err != nil {
					continue
				}

				var sz uint64
				if fi.IsDir() {
					s, _, _ := scanner.CalculateDirSize(full)
					sz = s
				} else {
					sz = uint64(fi.Size())
				}

				if sz > 0 {
					dataSz += sz
					components = append(components, AppComponent{
						Name:    c.nameFmt,
						Path:    full,
						Type:    c.typ,
						Size:    sz,
						SizeStr: FormatBytes(sz),
					})
				}
			}

			totalApp := appSz + dataSz

			item := InstalledApp{
				ID:           p,
				Name:         displayName,
				Version:      version,
				BundleID:     bundleID,
				Path:         p,
				AppSize:      appSz,
				AppSizeStr:   FormatBytes(appSz),
				DataSize:     dataSz,
				DataSizeStr:  FormatBytes(dataSz),
				TotalSize:    totalApp,
				TotalSizeStr: FormatBytes(totalApp),
				IsSystem:     isSystem,
				Components:   components,
			}

			mu.Lock()
			apps = append(apps, item)
			atomic.AddUint64(&totalAppBytes, appSz)
			atomic.AddUint64(&totalDataBytes, dataSz)
			mu.Unlock()
		}(appPath)
	}

	wg.Wait()

	// Sort apps by TotalSize descending
	sort.Slice(apps, func(i, j int) bool {
		return apps[i].TotalSize > apps[j].TotalSize
	})

	return &UninstallerScanResult{
		Apps:             apps,
		TotalAppSize:     totalAppBytes,
		TotalAppSizeStr:  FormatBytes(totalAppBytes),
		TotalDataSize:    totalDataBytes,
		TotalDataSizeStr: FormatBytes(totalDataBytes),
		TotalAppsCount:   len(apps),
	}, nil
}

// MoveToTrash moves a file or folder to macOS user Trash using AppleScript or fallback
func MoveToTrash(path string) error {
	cleanPath := filepath.Clean(path)
	if _, err := os.Stat(cleanPath); err != nil {
		return err
	}

	// Try AppleScript Finder delete (native macOS Trash behavior)
	script := fmt.Sprintf(`tell application "Finder" to delete POSIX file %q`, cleanPath)
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback: move manually to ~/.Trash
	home, _ := os.UserHomeDir()
	trashDir := filepath.Join(home, ".Trash")
	base := filepath.Base(cleanPath)
	target := filepath.Join(trashDir, base)

	// If target exists, add unique timestamp
	if _, err := os.Stat(target); err == nil {
		target = filepath.Join(trashDir, fmt.Sprintf("%s_%d", base, os.Getpid()))
	}

	return os.Rename(cleanPath, target)
}

// UninstallApp performs complete uninstallation or data reset of an application
func UninstallApp(appID string, resetOnly bool, useTrash bool) (uint64, error) {
	cleanAppID := filepath.Clean(appID)
	if strings.HasPrefix(cleanAppID, "/System") {
		return 0, fmt.Errorf("sistem uygulamaları silinemez: %s", appID)
	}

	// Re-scan target app to find its exact components
	appsResult, err := ScanInstalledApplications(context.Background())
	if err != nil {
		return 0, err
	}

	var targetApp *InstalledApp
	for _, a := range appsResult.Apps {
		if a.ID == cleanAppID {
			targetApp = &a
			break
		}
	}

	if targetApp == nil {
		return 0, fmt.Errorf("uygulama bulunamadı: %s", appID)
	}

	var freedBytes uint64

	for _, comp := range targetApp.Components {
		if comp.Type == "bundle" && resetOnly {
			// In reset mode, do NOT touch the .app bundle!
			continue
		}

		if _, err := os.Stat(comp.Path); err != nil {
			continue
		}

		if useTrash {
			if err := MoveToTrash(comp.Path); err == nil {
				freedBytes += comp.Size
			} else {
				// Fallback to permanent delete if trash fails
				_ = os.RemoveAll(comp.Path)
				freedBytes += comp.Size
			}
		} else {
			if err := os.RemoveAll(comp.Path); err == nil {
				freedBytes += comp.Size
			}
		}
	}

	return freedBytes, nil
}

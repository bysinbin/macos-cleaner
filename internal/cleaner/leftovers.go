package cleaner

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
)

// LeftoverCategory classifies residual file types
type LeftoverCategory string

const (
	CatOrphanedApp LeftoverCategory = "orphaned_app" // Uninstalled application residual data
	CatSavedState  LeftoverCategory = "saved_state"  // Saved Application States of uninstalled apps
	CatDSStore     LeftoverCategory = "ds_store"     // .DS_Store files across projects
	CatBrokenLink  LeftoverCategory = "broken_link"  // Broken symlinks pointing nowhere
	CatTempFile    LeftoverCategory = "temp_file"    // .tmp, .swp, .bak, *~
)

// LeftoverItem represents a detected residual file or folder
type LeftoverItem struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Path        string           `json:"path"`
	Category    LeftoverCategory `json:"category"`
	CategoryStr string           `json:"categoryStr"`
	Description string           `json:"description"`
	Size        uint64           `json:"size"`
	SizeStr     string           `json:"sizeStr"`
	FileCount   int64            `json:"fileCount"`
	Selected    bool             `json:"selected"`
}

// LeftoverScanResult contains the full scan outcome
type LeftoverScanResult struct {
	Items             []LeftoverItem `json:"items"`
	TotalSize         uint64         `json:"totalSize"`
	TotalSizeStr      string         `json:"totalSizeStr"`
	TotalFiles        int64          `json:"totalFiles"`
	OrphanedAppsCount int            `json:"orphanedAppsCount"`
	DSStoreCount      int64          `json:"dsStoreCount"`
}

// Common generic application terms that should not be used as standalone matching tokens
var genericWords = map[string]bool{
	"desktop": true, "app": true, "apps": true, "application": true, "applications": true,
	"helper": true, "helpers": true, "client": true, "clients": true, "service": true,
	"services": true, "server": true, "servers": true, "editor": true, "viewer": true,
	"launcher": true, "tool": true, "tools": true, "suite": true, "manager": true,
	"utility": true, "utilities": true, "support": true, "cloud": true, "community": true,
	"edition": true, "updater": true, "installer": true, "uninstaller": true, "plugin": true,
	"framework": true, "frameworks": true, "shared": true, "common": true, "engine": true,
	"renderer": true, "agent": true, "daemon": true, "system": true, "macos": true,
	"package": true, "standard": true, "professional": true, "multi": true, "login": true,
	"open": true, "free": true, "plus": true, "ultra": true, "total": true, "smart": true,
}

// System/developer folders inside ~/Library/Application Support to NEVER touch
var ignoredAppSupport = map[string]bool{
	"apple": true, "appstore": true, "addressbook": true, "callhistorydb": true,
	"callhistorytransactions": true, "controlcenter": true, "clouddocs": true,
	"diskimages": true, "dock": true, "facetime": true, "fileprovider": true,
	"knowledge": true, "differentialprivacy": true, "crashreporter": true,
	"animoji": true, "bookmarks": true, "cef": true, "caches": true,
	"audiocomponents": true, "syncservices": true, "icdd": true,
	"diskcleaner": true, "icloud": true, "kpeople": true, "kpeoplevcard": true,
	"networkserviceproxy": true, "videosubscriptionsd": true, "stickersd": true,
	"tipsd": true, "familycircled": true, "homeenergyd": true, "defaultstore": true,
	"contactsd": true, "locationaccessstored": true, "privatecloudcomputed": true,
	"identityservicesd": true, "icloudmailagent": true, "appplaceholdersyncd": true,
	"ilifemediabrowser": true, "colors": true, "db": true, "go": true, "gopls": true,
	"kotlin": true, "cloudcode": true, "mobilesync": true, "feedbackassistant": true,
	"bugsnag": true, "sentry": true, "sesstorage": true,
}

func normalizeString(s string) string {
	s = strings.ToLower(s)
	var sb strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func stripPlatformSuffixes(s string) string {
	s = strings.TrimRight(s, "0123456789")
	suffixes := []string{"mac", "osx", "app", "desktop"}
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) && len(s) > len(suf)+2 {
			s = strings.TrimSuffix(s, suf)
			break
		}
	}
	return strings.TrimRight(s, "0123456789")
}

// splitIntoWords splits camelCase, kebab-case, snake_case, dot.case into lowercase words
func splitIntoWords(s string) []string {
	var words []string
	var current strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '.' || r == '-' || r == '_' || r == ' ' || r == '/' || r == '\\' {
			if current.Len() > 0 {
				words = append(words, current.String())
				current.Reset()
			}
			continue
		}

		if unicode.IsUpper(r) {
			if current.Len() > 0 && i > 0 && unicode.IsLower(runes[i-1]) {
				words = append(words, current.String())
				current.Reset()
			}
		}

		current.WriteRune(unicode.ToLower(r))
	}

	if current.Len() > 0 {
		words = append(words, current.String())
	}

	return words
}

// InstalledAppTracker indexes all installed applications across depths and resolves vendor/bundle mappings
type InstalledAppTracker struct {
	exactMatches      map[string]bool
	distinctiveTokens map[string]bool
}

// NewInstalledAppTracker recursively scans /Applications, /System/Applications, and ~/Applications
func NewInstalledAppTracker() *InstalledAppTracker {
	tracker := &InstalledAppTracker{
		exactMatches:      make(map[string]bool),
		distinctiveTokens: make(map[string]bool),
	}

	// Bi-directional vendor & suite mapping
	vendorMap := map[string][]string{
		"techsmith":     {"camtasia", "snagit"},
		"camtasia":      {"techsmith"},
		"snagit":        {"techsmith"},
		"adobe":         {"photoshop", "illustrator", "creativecloud", "aftereffects", "xd", "mediaencoder", "indesign", "acrobat", "lightroom", "dunamis", "ccxprocess", "premiere"},
		"photoshop":     {"adobe"},
		"illustrator":   {"adobe"},
		"creativecloud": {"adobe"},
		"aftereffects":  {"adobe"},
		"mediaencoder":  {"adobe"},
		"autodesk":      {"eagle", "fusion360", "maya", "3dsmax", "autocad"},
		"eagle":         {"autodesk"},
		"microsoft":     {"word", "excel", "powerpoint", "outlook", "teams", "onedrive", "vscode", "visualstudio", "code"},
		"jetbrains":     {"intellij", "webstorm", "pycharm", "goland", "clion", "rider", "datagrip", "rustrover"},
		"valve":         {"steam", "valvecorporation"},
		"steam":         {"valve", "valvecorporation"},
		"epicgames":     {"epic", "unreal"},
		"epic":          {"epicgames"},
		"asus":          {"armourycrate", "asusframework"},
		"armourycrate":  {"asus"},
		"oracle":        {"java", "jreinstaller"},
		"maxon":         {"cinema4d", "c4d"},
	}

	addToken := func(t string) {
		norm := normalizeString(t)
		if len(norm) >= 2 {
			tracker.exactMatches[norm] = true
			stripped := stripPlatformSuffixes(norm)
			if len(stripped) >= 3 {
				tracker.exactMatches[stripped] = true
			}

			if !genericWords[norm] && len(norm) >= 3 {
				tracker.distinctiveTokens[norm] = true
			}
			if !genericWords[stripped] && len(stripped) >= 3 {
				tracker.distinctiveTokens[stripped] = true
			}

			if aliases, ok := vendorMap[norm]; ok {
				for _, a := range aliases {
					na := normalizeString(a)
					tracker.exactMatches[na] = true
					if !genericWords[na] && len(na) >= 3 {
						tracker.distinctiveTokens[na] = true
					}
				}
			}
			if aliases, ok := vendorMap[stripped]; ok {
				for _, a := range aliases {
					na := normalizeString(a)
					tracker.exactMatches[na] = true
					if !genericWords[na] && len(na) >= 3 {
						tracker.distinctiveTokens[na] = true
					}
				}
			}
		}
	}

	roots := []string{
		"/Applications",
		"/System/Applications",
		ExpandPath("~/Applications"),
	}

	bundleIDRegex := regexp.MustCompile(`(?:com|org|io|net|de|app)\.[a-zA-Z0-9_-]+(?:\.[a-zA-Z0-9_-]+)+`)

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
				base := strings.TrimSuffix(d.Name(), ".app")
				addToken(base)
				for _, w := range splitIntoWords(base) {
					addToken(w)
				}

				// Check parent directory name (e.g. EAGLE-9.6.2, Adobe Photoshop 2026, PostgreSQL 17)
				parent := filepath.Base(filepath.Dir(p))
				if parent != "Applications" && parent != "System" && parent != "Utilities" {
					addToken(parent)
					for _, w := range splitIntoWords(parent) {
						addToken(w)
					}
				}

				// Inspect Info.plist for CFBundleIdentifier
				plistPath := filepath.Join(p, "Contents", "Info.plist")
				if data, err := os.ReadFile(plistPath); err == nil {
					if len(data) > 128*1024 {
						data = data[:128*1024]
					}
					matches := bundleIDRegex.FindAll(data, -1)
					for _, m := range matches {
						parts := bytes.Split(m, []byte("."))
						for _, seg := range parts {
							s := string(seg)
							if s == "com" || s == "org" || s == "io" || s == "net" || s == "de" || s == "app" || s == "apple" {
								continue
							}
							addToken(s)
							for _, w := range splitIntoWords(s) {
								addToken(w)
							}
						}
					}
				}

				// Skip traversing into the .app bundle itself
				return filepath.SkipDir
			}

			return nil
		})
	}

	return tracker
}

// IsInstalled checks if a given directory or state name belongs to an active installed application
func (t *InstalledAppTracker) IsInstalled(name string) bool {
	norm := normalizeString(name)
	if t.exactMatches[norm] {
		return true
	}

	stripped := stripPlatformSuffixes(norm)
	if t.exactMatches[stripped] {
		return true
	}

	words := splitIntoWords(name)
	for _, w := range words {
		if genericWords[w] || len(w) < 3 {
			continue
		}
		if t.distinctiveTokens[w] {
			return true
		}
		sw := stripPlatformSuffixes(w)
		if len(sw) >= 3 && !genericWords[sw] && t.distinctiveTokens[sw] {
			return true
		}
	}

	return false
}

// GetInstalledAppNames returns a set of normalized names and bundle stems of installed applications
func GetInstalledAppNames() map[string]bool {
	tracker := NewInstalledAppTracker()
	return tracker.exactMatches
}

func normalizeAppName(s string) string {
	return normalizeString(s)
}

// ScanLeftovers searches for orphaned app data, leftover states, and system debris
func ScanLeftovers(ctx context.Context) (*LeftoverScanResult, error) {
	scanner := NewScanner(ctx)
	tracker := NewInstalledAppTracker()

	var items []LeftoverItem
	var mu sync.Mutex
	var wg sync.WaitGroup

	var totalSize uint64
	var totalFiles int64
	var orphanedCount int
	var dsStoreCount int64

	// 1. Scan ~/Library/Application Support for orphaned app folders
	appSupportDir := ExpandPath("~/Library/Application Support")
	if entries, err := os.ReadDir(appSupportDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			name := e.Name()
			norm := normalizeString(name)

			if strings.HasPrefix(norm, "comapple") || ignoredAppSupport[norm] {
				continue
			}

			if tracker.IsInstalled(name) {
				continue
			}

			// Orphaned application folder found!
			fullPath := filepath.Join(appSupportDir, name)
			wg.Add(1)
			go func(folderName, path string) {
				defer wg.Done()
				sz, count, _ := scanner.CalculateDirSize(path)
				if sz < 10*1024 { // Skip tiny empty folders < 10KB
					return
				}

				item := LeftoverItem{
					ID:          path,
					Name:        folderName,
					Path:        path,
					Category:    CatOrphanedApp,
					CategoryStr: "Kaldırılmış Uygulama Verisi",
					Description: "Bu uygulama artık sistemde yüklü değil.",
					Size:        sz,
					SizeStr:     FormatBytes(sz),
					FileCount:   count,
					Selected:    true,
				}

				mu.Lock()
				items = append(items, item)
				atomic.AddUint64(&totalSize, sz)
				atomic.AddInt64(&totalFiles, count)
				orphanedCount++
				mu.Unlock()
			}(name, fullPath)
		}
	}

	// 2. Scan ~/Library/Saved Application State
	savedStateDir := ExpandPath("~/Library/Saved Application State")
	if entries, err := os.ReadDir(savedStateDir); err == nil {
		for _, e := range entries {
			name := e.Name()
			if !strings.HasSuffix(name, ".savedState") {
				continue
			}

			stem := strings.TrimSuffix(name, ".savedState")
			norm := normalizeString(stem)

			if strings.HasPrefix(norm, "comapple") {
				continue
			}

			if tracker.IsInstalled(stem) {
				continue
			}

			fullPath := filepath.Join(savedStateDir, name)
			wg.Add(1)
			go func(stateName, path string) {
				defer wg.Done()
				sz, count, _ := scanner.CalculateDirSize(path)

				item := LeftoverItem{
					ID:          path,
					Name:        stateName,
					Path:        path,
					Category:    CatSavedState,
					CategoryStr: "Eski Uygulama Durumu (State)",
					Description: "Kaldırılmış uygulamanın pencere ve oturum hafıza dosyaları.",
					Size:        sz,
					SizeStr:     FormatBytes(sz),
					FileCount:   count,
					Selected:    true,
				}

				mu.Lock()
				items = append(items, item)
				atomic.AddUint64(&totalSize, sz)
				atomic.AddInt64(&totalFiles, count)
				mu.Unlock()
			}(name, fullPath)
		}
	}

	// 3. Scan for .DS_Store, broken symlinks & temp files in user work directories
	home, _ := os.UserHomeDir()
	scanRoots := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Documents"),
	}

	var dsStoreFiles []string
	var dsStoreBytes uint64

	var brokenLinks []string
	var tempFiles []string
	var tempBytes uint64

	for _, root := range scanRoots {
		if _, err := os.Stat(root); err != nil {
			continue
		}

		_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			name := d.Name()
			if d.IsDir() {
				if name == ".git" || name == ".gemini" || name == "Library" || name == "node_modules" {
					return filepath.SkipDir
				}
				return nil
			}

			// Check .DS_Store
			if name == ".DS_Store" {
				fi, err := d.Info()
				if err == nil {
					dsStoreFiles = append(dsStoreFiles, p)
					dsStoreBytes += uint64(fi.Size())
				}
				return nil
			}

			// Check broken symlink
			if d.Type()&fs.ModeSymlink != 0 {
				if _, err := os.Stat(p); os.IsNotExist(err) {
					brokenLinks = append(brokenLinks, p)
				}
				return nil
			}

			// Check temporary files (*.tmp, *.swp, *~, *.bak)
			if strings.HasSuffix(name, ".tmp") || strings.HasSuffix(name, ".swp") ||
				strings.HasSuffix(name, "~") || strings.HasSuffix(name, ".bak") {
				fi, err := d.Info()
				if err == nil {
					tempFiles = append(tempFiles, p)
					tempBytes += uint64(fi.Size())
				}
			}

			return nil
		})
	}

	wg.Wait()

	// If .DS_Store files found, aggregate them as a cleanable bundle
	if len(dsStoreFiles) > 0 {
		dsStoreCount = int64(len(dsStoreFiles))
		item := LeftoverItem{
			ID:          "all-ds-store",
			Name:        "macOS .DS_Store Finder Artıkları",
			Path:        "Kullanıcı Çalışma Dizinleri (Desktop, Documents, Downloads)",
			Category:    CatDSStore,
			CategoryStr: "Finder Görünüm Artığı (.DS_Store)",
			Description: "Finder tarafından oluşturulmuş gereksiz pencere düzeni ve önbellek dosyaları.",
			Size:        dsStoreBytes,
			SizeStr:     FormatBytes(dsStoreBytes),
			FileCount:   dsStoreCount,
			Selected:    true,
		}
		items = append(items, item)
		totalSize += dsStoreBytes
		totalFiles += dsStoreCount
	}

	// If broken links found
	if len(brokenLinks) > 0 {
		for _, link := range brokenLinks {
			items = append(items, LeftoverItem{
				ID:          link,
				Name:        filepath.Base(link) + " (Kırık Link)",
				Path:        link,
				Category:    CatBrokenLink,
				CategoryStr: "Kırık Sembolik Link",
				Description: "Hedef dosyası veya dizini silinmiş geçersiz kısayol linki.",
				Size:        0,
				SizeStr:     "0 B",
				FileCount:   1,
				Selected:    true,
			})
		}
	}

	// If temp files found
	if len(tempFiles) > 0 {
		item := LeftoverItem{
			ID:          "all-temp-files",
			Name:        "Geçici ve Yedek Dosya Artıkları (*.tmp, *.swp, *~)",
			Path:        "Kullanıcı Çalışma Dizinleri",
			Category:    CatTempFile,
			CategoryStr: "Geçici Dosya Artıkları",
			Description: "Editörlerin veya programların geride bıraktığı geçici çalışma dosyaları.",
			Size:        tempBytes,
			SizeStr:     FormatBytes(tempBytes),
			FileCount:   int64(len(tempFiles)),
			Selected:    true,
		}
		items = append(items, item)
		totalSize += tempBytes
		totalFiles += int64(len(tempFiles))
	}

	// Sort items largest first
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	return &LeftoverScanResult{
		Items:             items,
		TotalSize:         totalSize,
		TotalSizeStr:      FormatBytes(totalSize),
		TotalFiles:        totalFiles,
		OrphanedAppsCount: orphanedCount,
		DSStoreCount:      dsStoreCount,
	}, nil
}

// CleanDSStoreBatch removes all .DS_Store files in user folders
func CleanDSStoreBatch() (uint64, int64, error) {
	home, _ := os.UserHomeDir()
	scanRoots := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Documents"),
	}

	var totalFreed uint64
	var deletedCount int64

	for _, root := range scanRoots {
		_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.Name() == ".DS_Store" {
				fi, err := d.Info()
				if err == nil {
					totalFreed += uint64(fi.Size())
				}
				_ = os.Remove(p)
				deletedCount++
			}
			return nil
		})
	}

	return totalFreed, deletedCount, nil
}

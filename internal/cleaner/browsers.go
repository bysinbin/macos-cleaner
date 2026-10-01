package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

// BrowserCacheItem represents a detected browser and its caches
type BrowserCacheItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Installed   bool     `json:"installed"`
	CachePaths  []string `json:"cachePaths"`
	Size        uint64   `json:"size"`
	SizeStr     string   `json:"sizeStr"`
	FileCount   int64    `json:"fileCount"`
	Selected    bool     `json:"selected"`
	Description string   `json:"description"`
}

// BrowsersScanResult represents full browser cache status
type BrowsersScanResult struct {
	Browsers     []BrowserCacheItem `json:"browsers"`
	TotalSize    uint64             `json:"totalSize"`
	TotalSizeStr string             `json:"totalSizeStr"`
	TotalFiles   int64              `json:"totalFiles"`
}

type browserDef struct {
	id          string
	name        string
	appNames    []string
	cacheRelDir []string
	desc        string
}

// ScanBrowserCaches scans all major browsers installed on macOS
func ScanBrowserCaches(ctx context.Context) (*BrowsersScanResult, error) {
	scanner := NewScanner(ctx)
	home, _ := os.UserHomeDir()

	defs := []browserDef{
		{
			id:       "chrome",
			name:     "Google Chrome",
			appNames: []string{"Google Chrome.app"},
			cacheRelDir: []string{
				"Library/Caches/Google/Chrome/Default/Cache",
				"Library/Caches/Google/Chrome/Default/Code Cache",
				"Library/Caches/Google/Chrome/Default/GPUCache",
				"Library/Caches/Google/Chrome/Default/Service Worker/CacheStorage",
				"Library/Caches/Google/Chrome/Profile 1/Cache",
				"Library/Caches/Google/Chrome/Profile 2/Cache",
			},
			desc: "Web sitesi önbellekleri, derlenmiş JS kodları ve GPU gölgelendirici önbellekleri.",
		},
		{
			id:       "safari",
			name:     "Apple Safari",
			appNames: []string{"Safari.app"},
			cacheRelDir: []string{
				"Library/Caches/com.apple.Safari",
				"Library/Caches/com.apple.Safari.SafeBrowsing",
				"Library/Containers/com.apple.Safari/Data/Library/Caches",
			},
			desc: "Safari web önbellekleri, ikon veritabanı ve güvenli tarama önbelleği.",
		},
		{
			id:       "firefox",
			name:     "Mozilla Firefox",
			appNames: []string{"Firefox.app"},
			cacheRelDir: []string{
				"Library/Caches/Firefox/Profiles/*/cache2",
				"Library/Caches/Firefox/Profiles/*/startupCache",
			},
			desc: "Firefox web sayfaları, resimler ve profil başlangıç önbelleği.",
		},
		{
			id:       "edge",
			name:     "Microsoft Edge",
			appNames: []string{"Microsoft Edge.app"},
			cacheRelDir: []string{
				"Library/Caches/Microsoft Edge/Default/Cache",
				"Library/Caches/Microsoft Edge/Default/Code Cache",
				"Library/Caches/Microsoft Edge/Default/GPUCache",
			},
			desc: "Edge web verileri, betik önbellekleri ve indirme geçmişi tamponları.",
		},
		{
			id:       "brave",
			name:     "Brave Browser",
			appNames: []string{"Brave Browser.app"},
			cacheRelDir: []string{
				"Library/Caches/BraveSoftware/Brave-Browser/Default/Cache",
				"Library/Caches/BraveSoftware/Brave-Browser/Default/Code Cache",
				"Library/Caches/BraveSoftware/Brave-Browser/Default/GPUCache",
			},
			desc: "Brave kalkan ve web tarayıcı önbellekleri.",
		},
		{
			id:       "arc",
			name:     "Arc Browser",
			appNames: []string{"Arc.app"},
			cacheRelDir: []string{
				"Library/Caches/company.thebrowser.Browser/Cache",
				"Library/Caches/company.thebrowser.Browser/Code Cache",
			},
			desc: "Arc Browser sekme ve web sayfası önbellekleri.",
		},
	}

	var items []BrowserCacheItem
	var mu sync.Mutex
	var wg sync.WaitGroup

	var totalSize uint64
	var totalFiles int64

	for _, def := range defs {
		wg.Add(1)
		go func(b browserDef) {
			defer wg.Done()

			// Check if installed
			isInstalled := false
			for _, app := range b.appNames {
				if _, err := os.Stat("/Applications/" + app); err == nil {
					isInstalled = true
					break
				}
				if _, err := os.Stat("/System/Applications/" + app); err == nil {
					isInstalled = true
					break
				}
			}

			var bSize uint64
			var bFiles int64
			var existingPaths []string

			for _, rel := range b.cacheRelDir {
				full := filepath.Join(home, rel)
				if strings.Contains(rel, "*") {
					matches, err := filepath.Glob(full)
					if err == nil {
						for _, m := range matches {
							if fi, err := os.Stat(m); err == nil {
								var sz uint64
								var cnt int64
								if fi.IsDir() {
									sz, cnt, _ = scanner.CalculateDirSize(m)
								} else {
									sz = uint64(fi.Size())
									cnt = 1
								}
								if sz > 0 {
									bSize += sz
									bFiles += cnt
									existingPaths = append(existingPaths, m)
								}
							}
						}
					}
					continue
				}

				if fi, err := os.Stat(full); err == nil {
					var sz uint64
					var cnt int64
					if fi.IsDir() {
						sz, cnt, _ = scanner.CalculateDirSize(full)
					} else {
						sz = uint64(fi.Size())
						cnt = 1
					}
					if sz > 0 {
						bSize += sz
						bFiles += cnt
						existingPaths = append(existingPaths, full)
					}
				}
			}

			// If browser is installed or cache has files
			if isInstalled || bSize > 0 {
				item := BrowserCacheItem{
					ID:          b.id,
					Name:        b.name,
					Installed:   isInstalled,
					CachePaths:  existingPaths,
					Size:        bSize,
					SizeStr:     FormatBytes(bSize),
					FileCount:   bFiles,
					Selected:    bSize > 0,
					Description: b.desc,
				}

				mu.Lock()
				items = append(items, item)
				atomic.AddUint64(&totalSize, bSize)
				atomic.AddInt64(&totalFiles, bFiles)
				mu.Unlock()
			}
		}(def)
	}

	wg.Wait()

	return &BrowsersScanResult{
		Browsers:     items,
		TotalSize:    totalSize,
		TotalSizeStr: FormatBytes(totalSize),
		TotalFiles:   totalFiles,
	}, nil
}

// CleanBrowserCaches removes cached files for specified browsers
func CleanBrowserCaches(browserIDs []string) (uint64, int64, error) {
	res, err := ScanBrowserCaches(context.Background())
	if err != nil {
		return 0, 0, err
	}

	idMap := make(map[string]bool)
	for _, id := range browserIDs {
		idMap[id] = true
	}

	var freed uint64
	var count int64

	for _, b := range res.Browsers {
		if len(idMap) > 0 && !idMap[b.ID] {
			continue
		}

		for _, p := range b.CachePaths {
			fi, err := os.Stat(p)
			if err != nil {
				continue
			}

			var sz uint64
			if fi.IsDir() {
				scanner := NewScanner(context.Background())
				s, c, _ := scanner.CalculateDirSize(p)
				sz = s
				count += c
			} else {
				sz = uint64(fi.Size())
				count++
			}

			if err := os.RemoveAll(p); err == nil {
				freed += sz
			}
		}
	}

	return freed, count, nil
}

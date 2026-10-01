package cleaner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
)

// ScanProgress communicates scanning progress
type ScanProgress struct {
	CurrentPath string `json:"currentPath"`
	ScannedDirs int64  `json:"scannedDirs"`
	ScannedFiles int64 `json:"scannedFiles"`
	FoundBytes  uint64 `json:"foundBytes"`
	Done        bool   `json:"done"`
}

// AppCacheItem represents a dynamically discovered cache directory in ~/Library/Caches
type AppCacheItem struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      uint64 `json:"size"`
	SizeStr   string `json:"sizeStr"`
	FileCount int64  `json:"fileCount"`
	Selected  bool   `json:"selected"`
}

// ScanResult holds all discovered cleanable data
type ScanResult struct {
	Targets        []CleanTarget  `json:"targets"`
	DynamicCaches  []AppCacheItem `json:"dynamicCaches"`
	TotalCleanable uint64         `json:"totalCleanable"`
	TotalCleanableStr string      `json:"totalCleanableStr"`
	DiskStats      *DiskStats     `json:"diskStats"`
}

// Scanner handles analyzing disk usage
type Scanner struct {
	ctx context.Context
}

// NewScanner creates a new Scanner instance
func NewScanner(ctx context.Context) *Scanner {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Scanner{ctx: ctx}
}

// CalculateDirSize computes the size and file count of a directory
func (s *Scanner) CalculateDirSize(path string) (uint64, int64, error) {
	var totalSize uint64
	var fileCount int64

	info, err := os.Lstat(path)
	if err != nil {
		return 0, 0, err
	}

	if !info.IsDir() {
		return uint64(info.Size()), 1, nil
	}

	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip unreadable entries gracefully
		}

		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}

		if !d.IsDir() {
			fi, err := d.Info()
			if err == nil {
				atomic.AddUint64(&totalSize, uint64(fi.Size()))
				atomic.AddInt64(&fileCount, 1)
			}
		}
		return nil
	})

	return totalSize, fileCount, err
}

// ScanTargets scans predefined targets and dynamic caches concurrently
func (s *Scanner) ScanTargets(progressChan chan<- ScanProgress) (*ScanResult, error) {
	targets := GetPredefinedTargets()
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 8) // Limit concurrent I/O

	var scannedDirs int64
	var scannedFiles int64
	var totalCleanable uint64

	var mapMu sync.Mutex
	knownPaths := make(map[string]bool)

	// Scan predefined targets
	for i := range targets {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			t := &targets[idx]
			resolved := ExpandPath(t.Path)
			t.Resolved = resolved

			info, err := os.Lstat(resolved)
			if err != nil || (!info.IsDir() && info.Mode()&os.ModeSymlink == 0) {
				t.Exists = false
				t.Size = 0
				t.SizeStr = "0 B"
				return
			}

			t.Exists = true
			mapMu.Lock()
			knownPaths[resolved] = true
			mapMu.Unlock()

			size, count, _ := s.CalculateDirSize(resolved)
			t.Size = size
			t.SizeStr = FormatBytes(size)
			t.ItemCount = count

			atomic.AddInt64(&scannedDirs, 1)
			atomic.AddInt64(&scannedFiles, count)
			atomic.AddUint64(&totalCleanable, size)

			if progressChan != nil {
				progressChan <- ScanProgress{
					CurrentPath:  t.Name,
					ScannedDirs:  atomic.LoadInt64(&scannedDirs),
					ScannedFiles: atomic.LoadInt64(&scannedFiles),
					FoundBytes:   atomic.LoadUint64(&totalCleanable),
				}
			}
		}(i)
	}

	wg.Wait()

	// Dynamic scan for other ~/Library/Caches items
	dynamicCaches := s.scanDynamicCaches(knownPaths, &scannedDirs, &scannedFiles, &totalCleanable, progressChan)

	diskStats, _ := GetDiskStats("/")

	result := &ScanResult{
		Targets:           targets,
		DynamicCaches:     dynamicCaches,
		TotalCleanable:    totalCleanable,
		TotalCleanableStr: FormatBytes(totalCleanable),
		DiskStats:         diskStats,
	}

	if progressChan != nil {
		progressChan <- ScanProgress{
			CurrentPath:  "Tarama tamamlandı",
			ScannedDirs:  atomic.LoadInt64(&scannedDirs),
			ScannedFiles: atomic.LoadInt64(&scannedFiles),
			FoundBytes:   totalCleanable,
			Done:         true,
		}
	}

	return result, nil
}

// scanDynamicCaches finds additional application caches not listed in predefined targets
func (s *Scanner) scanDynamicCaches(known map[string]bool, scannedDirs, scannedFiles *int64, totalCleanable *uint64, progressChan chan<- ScanProgress) []AppCacheItem {
	cachesDir := ExpandPath("~/Library/Caches")
	entries, err := os.ReadDir(cachesDir)
	if err != nil {
		return nil
	}

	var items []AppCacheItem
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 6)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(cachesDir, entry.Name())
		if known[fullPath] {
			continue // Already processed as predefined target
		}

		wg.Add(1)
		go func(name, path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			size, count, err := s.CalculateDirSize(path)
			if err != nil || size < 10*1024*1024 { // Filter out items < 10MB to keep list clean
				return
			}

			item := AppCacheItem{
				Name:      name,
				Path:      path,
				Size:      size,
				SizeStr:   FormatBytes(size),
				FileCount: count,
				Selected:  false, // User can opt-in
			}

			mu.Lock()
			items = append(items, item)
			mu.Unlock()

			atomic.AddInt64(scannedDirs, 1)
			atomic.AddInt64(scannedFiles, count)

			if progressChan != nil {
				progressChan <- ScanProgress{
					CurrentPath:  "Önbellek: " + name,
					ScannedDirs:  atomic.LoadInt64(scannedDirs),
					ScannedFiles: atomic.LoadInt64(scannedFiles),
					FoundBytes:   atomic.LoadUint64(totalCleanable),
				}
			}
		}(entry.Name(), fullPath)
	}

	wg.Wait()

	// Sort dynamic caches by size descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	return items
}

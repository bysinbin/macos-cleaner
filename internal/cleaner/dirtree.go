package cleaner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// DirItem represents a child directory or file in the breakdown
type DirItem struct {
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	IsDir        bool    `json:"isDir"`
	Size         uint64  `json:"size"`
	SizeStr      string  `json:"sizeStr"`
	Percentage   float64 `json:"percentage"` // % of parent directory
	ItemCount    int64   `json:"itemCount"`
	LastModified string  `json:"lastModified"`
}

// DirBreakdown holds the analysis of a specific directory
type DirBreakdown struct {
	CurrentPath string    `json:"currentPath"`
	ParentPath  string    `json:"parentPath"`
	TotalSize   uint64    `json:"totalSize"`
	TotalStr    string    `json:"totalStr"`
	Items       []DirItem `json:"items"`
}

// AnalyzeDirectory inspects a directory and computes immediate child sizes
func AnalyzeDirectory(ctx context.Context, targetPath string) (*DirBreakdown, error) {
	if targetPath == "" {
		targetPath = "~"
	}
	resolved := ExpandPath(targetPath)

	info, err := os.Lstat(resolved)
	if err != nil {
		return nil, fmt.Errorf("dizin bulunamadı: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("belirtilen yol bir dizin değil: %s", resolved)
	}

	entries, err := os.ReadDir(resolved)
	if err != nil {
		return nil, fmt.Errorf("dizin okunamadı: %v", err)
	}

	scanner := NewScanner(ctx)
	var items []DirItem
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 8) // Concurrent workers

	var totalDirSize uint64

	for _, entry := range entries {
		// Skip special sockets or unreadable links
		subPath := filepath.Join(resolved, entry.Name())

		wg.Add(1)
		go func(name, path string, isDir bool) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				return
			default:
			}

			var size uint64
			var count int64
			var modStr string

			fi, err := os.Lstat(path)
			if err == nil {
				modStr = fi.ModTime().Format("02.01.2006 15:04")
			}

			if isDir {
				// Don't traverse inside .photoslibrary as a dir
				if strings.HasSuffix(name, ".photoslibrary") || strings.HasSuffix(name, ".app") {
					s, c, _ := scanner.CalculateDirSize(path)
					size = s
					count = c
				} else {
					s, c, _ := scanner.CalculateDirSize(path)
					size = s
					count = c
				}
			} else if fi != nil {
				size = uint64(fi.Size())
				count = 1
			}

			item := DirItem{
				Name:         name,
				Path:         path,
				IsDir:        isDir,
				Size:         size,
				SizeStr:      FormatBytes(size),
				ItemCount:    count,
				LastModified: modStr,
			}

			mu.Lock()
			items = append(items, item)
			totalDirSize += size
			mu.Unlock()
		}(entry.Name(), subPath, entry.IsDir())
	}

	wg.Wait()

	// Calculate percentages and sort descending by size
	for i := range items {
		if totalDirSize > 0 {
			items[i].Percentage = (float64(items[i].Size) / float64(totalDirSize)) * 100.0
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	parent := filepath.Dir(resolved)
	if resolved == "/" || resolved == parent {
		parent = ""
	}

	return &DirBreakdown{
		CurrentPath: resolved,
		ParentPath:  parent,
		TotalSize:   totalDirSize,
		TotalStr:    FormatBytes(totalDirSize),
		Items:       items,
	}, nil
}

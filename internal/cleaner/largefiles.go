package cleaner

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// LargeFile represents a discovered large file
type LargeFile struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         uint64    `json:"size"`
	SizeStr      string    `json:"sizeStr"`
	Extension    string    `json:"extension"`
	Category     string    `json:"category"` // "video", "archive", "installer", "document", "audio", "other"
	LastModified time.Time `json:"lastModified"`
	LastModStr   string    `json:"lastModStr"`
	AgeDays      int       `json:"ageDays"`
	Selected     bool      `json:"selected"`
}

// LargeFileFilter options
type LargeFileFilter struct {
	Paths       []string `json:"paths"`
	MinSizeBytes uint64  `json:"minSizeBytes"` // e.g. 50*1024*1024 (50MB)
	MinAgeDays  int      `json:"minAgeDays"`   // e.g. 30 days
	Category    string   `json:"category"`     // "all", "video", "installer", etc.
}

// CategorizeExtension maps file extension to a user-friendly category
func CategorizeExtension(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".mp4", ".mov", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".m4v":
		return "video"
	case ".dmg", ".pkg", ".iso", ".app", ".ipa":
		return "installer"
	case ".zip", ".tar", ".gz", ".tgz", ".rar", ".7z", ".bz2", ".xz":
		return "archive"
	case ".pdf", ".psd", ".ai", ".sketch", ".fig", ".doc", ".docx", ".xls", ".xlsx":
		return "document"
	case ".mp3", ".wav", ".flac", ".aac", ".m4a", ".ogg":
		return "audio"
	default:
		return "other"
	}
}

// ScanLargeFiles searches for files exceeding the size threshold
func ScanLargeFiles(ctx context.Context, filter LargeFileFilter) ([]LargeFile, error) {
	home, _ := os.UserHomeDir()
	if len(filter.Paths) == 0 {
		filter.Paths = []string{
			filepath.Join(home, "Downloads"),
			filepath.Join(home, "Desktop"),
			filepath.Join(home, "Documents"),
			filepath.Join(home, "Movies"),
			filepath.Join(home, "Pictures"),
		}
	}

	if filter.MinSizeBytes == 0 {
		filter.MinSizeBytes = 50 * 1024 * 1024 // 50MB default
	}

	var results []LargeFile
	var mu sync.Mutex
	now := time.Now()

	// Skip system, cache or hidden directories during large file scan
	shouldSkipDir := func(name string) bool {
		if strings.HasPrefix(name, ".") {
			return true
		}
		switch name {
		case "Library", "node_modules", "Caches", "Trash", "System", "Applications", "Pods", "vendor", "build", "dist":
			return true
		}
		if strings.HasSuffix(name, ".photoslibrary") || strings.HasSuffix(name, ".app") {
			return true
		}
		return false
	}

	for _, root := range filter.Paths {
		resolvedRoot := ExpandPath(root)
		if _, err := os.Stat(resolvedRoot); err != nil {
			continue
		}

		_ = filepath.WalkDir(resolvedRoot, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if d.IsDir() {
				if shouldSkipDir(d.Name()) {
					return filepath.SkipDir
				}
				return nil
			}

			fi, err := d.Info()
			if err != nil {
				return nil
			}

			fileSize := uint64(fi.Size())
			if fileSize < filter.MinSizeBytes {
				return nil
			}

			ageDays := int(now.Sub(fi.ModTime()).Hours() / 24)
			if filter.MinAgeDays > 0 && ageDays < filter.MinAgeDays {
				return nil
			}

			ext := filepath.Ext(p)
			category := CategorizeExtension(ext)
			if filter.Category != "" && filter.Category != "all" && category != filter.Category {
				return nil
			}

			item := LargeFile{
				ID:           p,
				Name:         d.Name(),
				Path:         p,
				Size:         fileSize,
				SizeStr:      FormatBytes(fileSize),
				Extension:    ext,
				Category:     category,
				LastModified: fi.ModTime(),
				LastModStr:   fi.ModTime().Format("02.01.2006 15:04"),
				AgeDays:      ageDays,
				Selected:     false,
			}

			mu.Lock()
			results = append(results, item)
			mu.Unlock()

			return nil
		})
	}

	// Sort largest first
	sort.Slice(results, func(i, j int) bool {
		return results[i].Size > results[j].Size
	})

	return results, nil
}

// RevealInFinder opens macOS Finder highlighting the specified file
func RevealInFinder(filePath string) error {
	cmd := exec.Command("open", "-R", filePath)
	return cmd.Run()
}

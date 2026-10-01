package cleaner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DuplicateFile represents a file in a duplicate set
type DuplicateFile struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       uint64 `json:"size"`
	SizeStr    string `json:"sizeStr"`
	ModTimeStr string `json:"modTimeStr"`
	AgeDays    int    `json:"ageDays"`
	Selected   bool   `json:"selected"`
}

// DuplicateGroup represents a group of identical files
type DuplicateGroup struct {
	Hash          string          `json:"hash"`
	Size          uint64          `json:"size"`
	SizeStr       string          `json:"sizeStr"`
	WastedSize    uint64          `json:"wastedSize"`
	WastedSizeStr string          `json:"wastedSizeStr"`
	Files         []DuplicateFile `json:"files"`
}

// FindDuplicates searches for duplicate files using size + SHA-256
func FindDuplicates(ctx context.Context, roots []string, minSizeBytes uint64) ([]DuplicateGroup, error) {
	home, _ := os.UserHomeDir()
	if len(roots) == 0 {
		roots = []string{
			filepath.Join(home, "Desktop"),
			filepath.Join(home, "Downloads"),
			filepath.Join(home, "Documents"),
			filepath.Join(home, "Movies"),
		}
	}

	if minSizeBytes == 0 {
		minSizeBytes = 1024 * 1024 // 1 MB default
	}

	// 1. Group by exact file size
	sizeMap := make(map[uint64][]string)
	now := time.Now()

	shouldSkipDir := func(name string) bool {
		if strings.HasPrefix(name, ".") {
			return true
		}
		switch name {
		case "Library", "node_modules", "Caches", ".Trash", "System", "Applications", "Pods", "vendor", "build", "dist":
			return true
		}
		return false
	}

	for _, root := range roots {
		resolved := ExpandPath(root)
		_ = filepath.WalkDir(resolved, func(p string, d fs.DirEntry, err error) error {
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

			sz := uint64(fi.Size())
			if sz >= minSizeBytes {
				sizeMap[sz] = append(sizeMap[sz], p)
			}
			return nil
		})
	}

	// 2. For sizes with >1 files, calculate hash
	var groups []DuplicateGroup
	for sz, paths := range sizeMap {
		if len(paths) < 2 {
			continue
		}

		hashMap := make(map[string][]DuplicateFile)
		for _, path := range paths {
			h, err := computeFileHash(path)
			if err != nil {
				continue
			}

			fi, _ := os.Stat(path)
			modStr := ""
			ageDays := 0
			if fi != nil {
				modStr = fi.ModTime().Format("02.01.2006 15:04")
				ageDays = int(now.Sub(fi.ModTime()).Hours() / 24)
			}

			hashMap[h] = append(hashMap[h], DuplicateFile{
				Path:       path,
				Name:       filepath.Base(path),
				Size:       sz,
				SizeStr:    FormatBytes(sz),
				ModTimeStr: modStr,
				AgeDays:    ageDays,
				Selected:   false,
			})
		}

		for h, files := range hashMap {
			if len(files) > 1 {
				// Mark copies (except first) as pre-selected for convenience
				for i := 1; i < len(files); i++ {
					files[i].Selected = true
				}

				wasted := sz * uint64(len(files)-1)
				groups = append(groups, DuplicateGroup{
					Hash:          h,
					Size:          sz,
					SizeStr:       FormatBytes(sz),
					WastedSize:    wasted,
					WastedSizeStr: FormatBytes(wasted),
					Files:         files,
				})
			}
		}
	}

	// Sort by wasted size descending
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].WastedSize > groups[j].WastedSize
	})

	return groups, nil
}

// computeFileHash calculates SHA-256 (reading first 32KB + full if small)
func computeFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	// Read full file (streamed)
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

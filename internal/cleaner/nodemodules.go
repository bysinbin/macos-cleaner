package cleaner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// NodeModuleItem represents a discovered node_modules directory
type NodeModuleItem struct {
	ID           string    `json:"id"`
	ProjectName  string    `json:"projectName"`
	ProjectPath  string    `json:"projectPath"`
	ModulesPath  string    `json:"modulesPath"`
	Size         uint64    `json:"size"`
	SizeStr      string    `json:"sizeStr"`
	FileCount    int64     `json:"fileCount"`
	LastModified time.Time `json:"lastModified"`
	LastModStr   string    `json:"lastModStr"`
	HasPkgJSON   bool      `json:"hasPkgJson"`
	Selected     bool      `json:"selected"`
}

// NodeModuleProgress reports real-time scan progress
type NodeModuleProgress struct {
	CurrentPath string           `json:"currentPath"`
	FoundCount  int              `json:"foundCount"`
	TotalBytes  uint64           `json:"totalBytes"`
	Done        bool             `json:"done"`
	Item        *NodeModuleItem  `json:"item,omitempty"`
}

// ScanNodeModules scans the specified root paths for node_modules directories
func ScanNodeModules(ctx context.Context, roots []string, progressChan chan<- NodeModuleProgress) ([]NodeModuleItem, error) {
	if len(roots) == 0 {
		home, _ := os.UserHomeDir()
		// Search common project locations
		candidates := []string{
			filepath.Join(home, "Desktop"),
			filepath.Join(home, "Documents"),
			filepath.Join(home, "Projects"),
			filepath.Join(home, "Development"),
			filepath.Join(home, "workspace"),
			filepath.Join(home, "yazılım"),
			filepath.Join(home, "dev"),
			filepath.Join(home, "code"),
			filepath.Join(home, "src"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				roots = append(roots, c)
			}
		}
		if len(roots) == 0 {
			roots = []string{home}
		}
	}

	var items []NodeModuleItem
	var mu sync.Mutex
	scanner := NewScanner(ctx)

	// Skip system or internal directories
	shouldSkipDir := func(name string) bool {
		switch name {
		case ".git", ".svn", ".hg", ".gemini", ".vscode", "Library", "System", "Applications",
			"homebrew", "Cellar", ".Trash", ".cache", "Caches", "Pods", ".cargo", ".gradle", "vendor":
			return true
		default:
			return false
		}
	}

	var foundCount int
	var totalBytes uint64

	for _, root := range roots {
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
				dirName := d.Name()
				if shouldSkipDir(dirName) {
					return filepath.SkipDir
				}

				isBuildDir := false
				projectDir := filepath.Dir(p)

				if dirName == "node_modules" {
					isBuildDir = true
				} else if dirName == "target" {
					// Rust Cargo project build target
					if _, err := os.Stat(filepath.Join(projectDir, "Cargo.toml")); err == nil {
						isBuildDir = true
					}
				} else if dirName == "build" {
					// Gradle, CMake, Android build directory
					if _, err := os.Stat(filepath.Join(projectDir, "build.gradle")); err == nil {
						isBuildDir = true
					} else if _, err := os.Stat(filepath.Join(projectDir, "build.gradle.kts")); err == nil {
						isBuildDir = true
					} else if _, err := os.Stat(filepath.Join(projectDir, "CMakeLists.txt")); err == nil {
						isBuildDir = true
					}
				}

				if isBuildDir {
					projectName := filepath.Base(projectDir)
					hasPkg := false
					if _, err := os.Stat(filepath.Join(projectDir, "package.json")); err == nil {
						hasPkg = true
					} else if _, err := os.Stat(filepath.Join(projectDir, "Cargo.toml")); err == nil {
						hasPkg = true
					}

					info, _ := d.Info()
					var modTime time.Time
					if info != nil {
						modTime = info.ModTime()
					}

					size, count, _ := scanner.CalculateDirSize(p)

					item := NodeModuleItem{
						ID:           p,
						ProjectName:  projectName,
						ProjectPath:  projectDir,
						ModulesPath:  p,
						Size:         size,
						SizeStr:      FormatBytes(size),
						FileCount:    count,
						LastModified: modTime,
						LastModStr:   modTime.Format("02.01.2006 15:04"),
						HasPkgJSON:   hasPkg,
						Selected:     false,
					}

					mu.Lock()
					items = append(items, item)
					foundCount++
					totalBytes += size
					mu.Unlock()

					if progressChan != nil {
						progressChan <- NodeModuleProgress{
							CurrentPath: p,
							FoundCount:  foundCount,
							TotalBytes:  totalBytes,
							Item:        &item,
						}
					}

					// Crucial: Do not traverse inside node_modules!
					return filepath.SkipDir
				}
			}
			return nil
		})
	}

	// Sort largest first
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	if progressChan != nil {
		progressChan <- NodeModuleProgress{
			CurrentPath: "Tamamlandı",
			FoundCount:  len(items),
			TotalBytes:  totalBytes,
			Done:        true,
		}
	}

	return items, nil
}

package cleaner

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// mediaExtensions set of visual screenshot and recording extensions
var mediaExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
	".gif":  true,
	".mp4":  true,
	".webm": true,
	".mov":  true,
	".bmp":  true,
}

// IsAntigravityMediaFile checks if a filename is an image or video artifact
func IsAntigravityMediaFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return mediaExtensions[ext]
}

// isAntigravityBrainPath checks if a resolved path targets an Antigravity brain directory
func isAntigravityBrainPath(path string) bool {
	clean := filepath.Clean(path)
	return strings.HasSuffix(clean, filepath.Join(".gemini", "antigravity-ide", "brain")) ||
		strings.HasSuffix(clean, filepath.Join(".gemini", "antigravity", "brain"))
}

// CalculateBrainMediaSize scans only temporary screenshots, recordings, and media files in brain
func CalculateBrainMediaSize(brainPath string) (uint64, int64, error) {
	resolved := ExpandPath(brainPath)
	info, err := os.Lstat(resolved)
	if err != nil || !info.IsDir() {
		return 0, 0, err
	}

	var totalSize uint64
	var fileCount int64

	// Collect paths to check (both antigravity-ide and legacy antigravity if present)
	pathsToCheck := []string{resolved}
	altPath := ExpandPath("~/.gemini/antigravity/brain")
	if resolved != altPath {
		if altInfo, err := os.Lstat(altPath); err == nil && altInfo.IsDir() {
			pathsToCheck = append(pathsToCheck, altPath)
		}
	}

	for _, basePath := range pathsToCheck {
		_ = filepath.WalkDir(basePath, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				return nil
			}

			// Check if file is in .tempmediaStorage or .user_uploaded or has media extension
			parentDir := filepath.Base(filepath.Dir(p))
			isTempMedia := parentDir == ".tempmediaStorage" || parentDir == ".user_uploaded"
			isMedia := isTempMedia || IsAntigravityMediaFile(d.Name())

			if isMedia {
				fi, err := d.Info()
				if err == nil {
					atomic.AddUint64(&totalSize, uint64(fi.Size()))
					atomic.AddInt64(&fileCount, 1)
				}
			}
			return nil
		})
	}

	return totalSize, fileCount, nil
}

// CleanBrainMedia deletes screenshots and temporary media while preserving logs and markdown artifacts
func (c *Cleaner) CleanBrainMedia(brainPath string, moveToTrash bool) (uint64, int64, error) {
	resolved := ExpandPath(brainPath)
	if err := IsPathSafe(resolved); err != nil {
		return 0, 0, err
	}

	size, count, err := CalculateBrainMediaSize(resolved)
	if err != nil {
		return 0, 0, err
	}

	if c.dryRun {
		return size, count, nil
	}

	pathsToClean := []string{resolved}
	altPath := ExpandPath("~/.gemini/antigravity/brain")
	if resolved != altPath {
		if altInfo, err := os.Lstat(altPath); err == nil && altInfo.IsDir() {
			pathsToClean = append(pathsToClean, altPath)
		}
	}

	var freedBytes uint64
	var deletedFiles int64
	var tempDirsToRemove []string
	now := time.Now()

	for _, basePath := range pathsToClean {
		_ = filepath.WalkDir(basePath, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if d.IsDir() {
				if d.Name() == ".tempmediaStorage" || d.Name() == ".user_uploaded" {
					tempDirsToRemove = append(tempDirsToRemove, p)
				}
				return nil
			}

			parentDir := filepath.Base(filepath.Dir(p))
			isTempMedia := parentDir == ".tempmediaStorage" || parentDir == ".user_uploaded"
			isMedia := isTempMedia || IsAntigravityMediaFile(d.Name())

			if isMedia {
				fi, err := d.Info()
				if err != nil {
					return nil
				}

				// Safety guard: skip files modified in the last 2 minutes to protect in-flight active subagents
				if now.Sub(fi.ModTime()) < 2*time.Minute {
					return nil
				}

				fileSize := uint64(fi.Size())

				if moveToTrash {
					cmd := exec.Command("osascript", "-e", `tell application "Finder" to delete POSIX file "`+p+`"`)
					if cmd.Run() == nil {
						atomic.AddUint64(&freedBytes, fileSize)
						atomic.AddInt64(&deletedFiles, 1)
						return nil
					}
				}

				if err := os.Remove(p); err == nil {
					atomic.AddUint64(&freedBytes, fileSize)
					atomic.AddInt64(&deletedFiles, 1)
				}
			}
			return nil
		})

		// Try removing empty temp media directories
		for _, tempDir := range tempDirsToRemove {
			_ = os.Remove(tempDir) // Will succeed only if empty
		}
	}

	return freedBytes, deletedFiles, nil
}

// ScanAntigravityDebris scans all Antigravity temporary data for Smart Care
func ScanAntigravityDebris(ctx context.Context) (uint64, int64, error) {
	scanner := NewScanner(ctx)
	var totalSize uint64
	var totalCount int64

	var wg sync.WaitGroup

	// 1. Brain Media
	wg.Add(1)
	go func() {
		defer wg.Done()
		sz, cnt, err := CalculateBrainMediaSize("~/.gemini/antigravity-ide/brain")
		if err == nil {
			atomic.AddUint64(&totalSize, sz)
			atomic.AddInt64(&totalCount, cnt)
		}
	}()

	// 2. Browser Recordings
	recordingsDirs := []string{
		ExpandPath("~/.gemini/antigravity-ide/browser_recordings"),
		ExpandPath("~/.gemini/antigravity/browser_recordings"),
	}
	for _, recDir := range recordingsDirs {
		if _, err := os.Stat(recDir); err == nil {
			wg.Add(1)
			go func(d string) {
				defer wg.Done()
				sz, cnt, err := scanner.CalculateDirSize(d)
				if err == nil {
					atomic.AddUint64(&totalSize, sz)
					atomic.AddInt64(&totalCount, cnt)
				}
			}(recDir)
		}
	}

	// 3. Scratchpad
	scratchDir := ExpandPath("~/.gemini/antigravity-ide/scratch")
	if _, err := os.Stat(scratchDir); err == nil {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			sz, cnt, err := scanner.CalculateDirSize(d)
			if err == nil {
				atomic.AddUint64(&totalSize, sz)
				atomic.AddInt64(&totalCount, cnt)
			}
		}(scratchDir)
	}

	// 4. Browser Profile Cache
	browserCacheDir := ExpandPath("~/.gemini/antigravity-browser-profile/Default/Cache")
	if _, err := os.Stat(browserCacheDir); err == nil {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			sz, cnt, err := scanner.CalculateDirSize(d)
			if err == nil {
				atomic.AddUint64(&totalSize, sz)
				atomic.AddInt64(&totalCount, cnt)
			}
		}(browserCacheDir)
	}

	wg.Wait()
	return totalSize, totalCount, nil
}

// CleanAntigravityDebris executes safe one-click cleaning for Smart Care
func CleanAntigravityDebris() (uint64, int64, error) {
	cleaner := NewCleaner(false)
	var totalFreed uint64
	var totalCount int64

	// 1. Brain screenshots & temp media
	sz, cnt, err := cleaner.CleanBrainMedia("~/.gemini/antigravity-ide/brain", false)
	if err == nil {
		totalFreed += sz
		totalCount += cnt
	}

	// 2. Browser recordings
	recDirs := []string{
		"~/.gemini/antigravity-ide/browser_recordings",
		"~/.gemini/antigravity/browser_recordings",
	}
	for _, rd := range recDirs {
		resolved := ExpandPath(rd)
		if _, err := os.Stat(resolved); err == nil {
			sz, cnt, err := cleaner.CleanTarget(resolved, false)
			if err == nil {
				totalFreed += sz
				totalCount += cnt
			}
		}
	}

	// 3. Scratch folder contents
	scratchDir := ExpandPath("~/.gemini/antigravity-ide/scratch")
	if _, err := os.Stat(scratchDir); err == nil {
		sz, cnt, err := cleaner.CleanTarget(scratchDir, false)
		if err == nil {
			totalFreed += sz
			totalCount += cnt
		}
	}

	// 4. Browser Profile Cache
	browserCacheDir := ExpandPath("~/.gemini/antigravity-browser-profile/Default/Cache")
	if _, err := os.Stat(browserCacheDir); err == nil {
		sz, cnt, err := cleaner.CleanTarget(browserCacheDir, false)
		if err == nil {
			totalFreed += sz
			totalCount += cnt
		}
	}

	return totalFreed, totalCount, nil
}

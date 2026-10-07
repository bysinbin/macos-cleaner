package cleaner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
)

// CleanRequest specifies what to clean
type CleanRequest struct {
	TargetIDs      []string `json:"targetIds"`      // IDs of predefined targets
	CustomPaths    []string `json:"customPaths"`    // Extra paths (e.g. dynamic caches, node_modules)
	DryRun         bool     `json:"dryRun"`         // Simulate only
	MoveToTrash    bool     `json:"moveToTrash"`    // Use Finder Trash instead of rm
}

// CleanProgress reports real-time deletion progress
type CleanProgress struct {
	ItemName     string   `json:"itemName"`
	Path         string   `json:"path"`
	FreedBytes   uint64   `json:"freedBytes"`
	FreedStr     string   `json:"freedStr"`
	DeletedFiles int64    `json:"deletedFiles"`
	Errors       []string `json:"errors,omitempty"`
	Done         bool     `json:"done"`
}

// Cleaner handles safe file deletion
type Cleaner struct {
	dryRun bool
}

// NewCleaner creates a Cleaner instance
func NewCleaner(dryRun bool) *Cleaner {
	return &Cleaner{dryRun: dryRun}
}

// IsPathSafe performs safety validation on any path before deletion
func IsPathSafe(targetPath string) error {
	cleanPath := filepath.Clean(targetPath)
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("kullanıcı ana dizini alınamadı")
	}

	// Never allow root or system paths
	forbiddenPrefixes := []string{
		"/",
		"/System",
		"/usr",
		"/bin",
		"/sbin",
		"/Library",
		"/Applications",
		"/private",
		"/etc",
		"/var",
	}

	for _, p := range forbiddenPrefixes {
		if cleanPath == p {
			return fmt.Errorf("kritik sistem dizini silinemez: %s", cleanPath)
		}
	}

	// Never allow home directory root or critical home configs
	if cleanPath == home {
		return fmt.Errorf("ana dizin (home) doğrudan silinemez: %s", cleanPath)
	}

	forbiddenHomePaths := []string{
		filepath.Join(home, ".ssh"),
		filepath.Join(home, ".gnupg"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".gitconfig"),
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Library"),
	}

	for _, p := range forbiddenHomePaths {
		if cleanPath == p {
			return fmt.Errorf("korunan kullanıcı dizini doğrudan silinemez: %s", cleanPath)
		}
	}

	// Check if path exists
	if _, err := os.Lstat(cleanPath); err != nil {
		return fmt.Errorf("dosya veya dizin bulunamadı: %s", cleanPath)
	}

	return nil
}

// CleanTarget cleans a specific target directory or file safely
func (c *Cleaner) CleanTarget(targetPath string, moveToTrash bool) (uint64, int64, error) {
	resolved := ExpandPath(targetPath)
	if err := IsPathSafe(resolved); err != nil {
		return 0, 0, err
	}

	// Special handling for Antigravity brain media (screenshots, recordings, and temp media)
	if isAntigravityBrainPath(resolved) {
		return c.CleanBrainMedia(resolved, moveToTrash)
	}

	// Special handling for VS Code obsolete extensions
	home, _ := os.UserHomeDir()
	if filepath.Clean(resolved) == filepath.Join(home, ".vscode", "extensions") {
		return CleanVSCodeObsoleteExtensions()
	}

	scanner := NewScanner(nil)
	size, count, err := scanner.CalculateDirSize(resolved)
	if err != nil {
		return 0, 0, err
	}

	if c.dryRun {
		return size, count, nil
	}

	// Special handling for macOS Trash
	if filepath.Base(resolved) == ".Trash" {
		err := c.emptyTrash()
		return size, count, err
	}

	if moveToTrash {
		// Use osascript to move to Finder Trash safely
		cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, resolved))
		if err := cmd.Run(); err == nil {
			return size, count, nil
		}
		// Fallback to direct clean if AppleScript fails
	}

	// Clean directory contents or delete and recreate
	info, err := os.Lstat(resolved)
	if err != nil {
		return 0, 0, err
	}

	if info.IsDir() {
		// Remove contents instead of deleting the directory itself for top-level caches
		entries, err := os.ReadDir(resolved)
		if err != nil {
			return 0, 0, err
		}

		for _, entry := range entries {
			subPath := filepath.Join(resolved, entry.Name())
			_ = os.RemoveAll(subPath)
		}
	} else {
		_ = os.Remove(resolved)
	}

	return size, count, nil
}

// emptyTrash empties the macOS Trash
func (c *Cleaner) emptyTrash() error {
	cmd := exec.Command("osascript", "-e", `tell application "Finder" to empty trash`)
	return cmd.Run()
}

// ExecuteClean processes batch clean requests with progress streaming
func (c *Cleaner) ExecuteClean(req CleanRequest, progressChan chan<- CleanProgress) {
	defer func() {
		if progressChan != nil {
			close(progressChan)
		}
	}()

	var totalFreed uint64
	var totalFiles int64
	var errList []string

	targetsMap := make(map[string]CleanTarget)
	for _, t := range GetPredefinedTargets() {
		targetsMap[t.ID] = t
	}

	// 1. Process predefined target IDs
	for _, id := range req.TargetIDs {
		target, ok := targetsMap[id]
		if !ok {
			continue
		}

		resolved := ExpandPath(target.Path)
		freed, count, err := c.CleanTarget(resolved, req.MoveToTrash)
		if err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", target.Name, err))
		} else {
			atomic.AddUint64(&totalFreed, freed)
			atomic.AddInt64(&totalFiles, count)
		}

		if progressChan != nil {
			progressChan <- CleanProgress{
				ItemName:     target.Name,
				Path:         resolved,
				FreedBytes:   atomic.LoadUint64(&totalFreed),
				FreedStr:     FormatBytes(atomic.LoadUint64(&totalFreed)),
				DeletedFiles: atomic.LoadInt64(&totalFiles),
				Errors:       errList,
			}
		}
	}

	// 2. Process custom paths (dynamic caches, node_modules, large files)
	for _, customPath := range req.CustomPaths {
		resolved := ExpandPath(customPath)
		freed, count, err := c.CleanTarget(resolved, req.MoveToTrash)
		if err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", filepath.Base(resolved), err))
		} else {
			atomic.AddUint64(&totalFreed, freed)
			atomic.AddInt64(&totalFiles, count)
		}

		if progressChan != nil {
			progressChan <- CleanProgress{
				ItemName:     filepath.Base(resolved),
				Path:         resolved,
				FreedBytes:   atomic.LoadUint64(&totalFreed),
				FreedStr:     FormatBytes(atomic.LoadUint64(&totalFreed)),
				DeletedFiles: atomic.LoadInt64(&totalFiles),
				Errors:       errList,
			}
		}
	}

	// Final progress message
	if progressChan != nil {
		progressChan <- CleanProgress{
			ItemName:     "Tamamlandı",
			FreedBytes:   totalFreed,
			FreedStr:     FormatBytes(totalFreed),
			DeletedFiles: totalFiles,
			Errors:       errList,
			Done:         true,
		}
	}
}

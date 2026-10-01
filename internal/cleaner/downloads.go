package cleaner

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DownloadItem represents a file in ~/Downloads
type DownloadItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Category  string `json:"category"` // "installer", "archive", "document", "media", "other"
	Size      uint64 `json:"size"`
	SizeStr   string `json:"sizeStr"`
	ModDate   string `json:"modDate"`
	DaysOld   int    `json:"daysOld"`
	IsOld     bool   `json:"isOld"`     // > 30 days
	IsLarge   bool   `json:"isLarge"`   // > 50 MB
	Selected  bool   `json:"selected"`
}

// DownloadsResult represents the full categorized downloads scan
type DownloadsResult struct {
	Items            []DownloadItem `json:"items"`
	TotalSize        uint64         `json:"totalSize"`
	TotalSizeStr     string         `json:"totalSizeStr"`
	TotalCount       int64          `json:"totalCount"`
	InstallersSize   uint64         `json:"installersSize"`
	InstallersStr    string         `json:"installersStr"`
	InstallersCount  int64          `json:"installersCount"`
	ArchivesSize     uint64         `json:"archivesSize"`
	ArchivesStr      string         `json:"archivesStr"`
	ArchivesCount    int64          `json:"archivesCount"`
	OldFilesSize     uint64         `json:"oldFilesSize"`
	OldFilesStr      string         `json:"oldFilesStr"`
	OldFilesCount    int64          `json:"oldFilesCount"`
}

func categorizeDownloadFile(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".dmg", ".pkg", ".iso", ".appinstaller":
		return "installer"
	case ".zip", ".tar", ".gz", ".tgz", ".rar", ".7z", ".bz2", ".xz":
		return "archive"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".csv":
		return "document"
	case ".mov", ".mp4", ".mkv", ".avi", ".jpg", ".jpeg", ".png", ".webp", ".mp3", ".wav":
		return "media"
	default:
		return "other"
	}
}

// ScanDownloadsDirectory analyzes the user's Downloads folder
func ScanDownloadsDirectory(ctx context.Context) (*DownloadsResult, error) {
	home, _ := os.UserHomeDir()
	downloadsDir := filepath.Join(home, "Downloads")
	now := time.Now()

	result := &DownloadsResult{}

	if _, err := os.Stat(downloadsDir); err != nil {
		return result, nil
	}

	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		return result, err
	}

	scanner := NewScanner(ctx)

	for _, e := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		name := e.Name()
		if name == ".DS_Store" || name == ".localized" || strings.HasPrefix(name, ".") {
			continue
		}

		full := filepath.Join(downloadsDir, name)

		var sz uint64
		var modTime time.Time

		fi, err := e.Info()
		if err != nil {
			continue
		}
		modTime = fi.ModTime()

		if e.IsDir() {
			s, _, _ := scanner.CalculateDirSize(full)
			sz = s
		} else {
			sz = uint64(fi.Size())
		}

		if sz == 0 {
			continue
		}

		ext := filepath.Ext(name)
		cat := categorizeDownloadFile(ext)
		daysOld := int(now.Sub(modTime).Hours() / 24)
		isOld := daysOld >= 30
		isLarge := sz >= 50*1024*1024

		// By default select installers older than 7 days
		selected := cat == "installer" && daysOld >= 7

		item := DownloadItem{
			ID:       full,
			Name:     name,
			Path:     full,
			Category: cat,
			Size:     sz,
			SizeStr:  FormatBytes(sz),
			ModDate:  modTime.Format("02.01.2006 15:04"),
			DaysOld:  daysOld,
			IsOld:    isOld,
			IsLarge:  isLarge,
			Selected: selected,
		}

		result.Items = append(result.Items, item)
		result.TotalSize += sz
		result.TotalCount++

		switch cat {
		case "installer":
			result.InstallersSize += sz
			result.InstallersCount++
		case "archive":
			result.ArchivesSize += sz
			result.ArchivesCount++
		}

		if isOld {
			result.OldFilesSize += sz
			result.OldFilesCount++
		}
	}

	result.TotalSizeStr = FormatBytes(result.TotalSize)
	result.InstallersStr = FormatBytes(result.InstallersSize)
	result.ArchivesStr = FormatBytes(result.ArchivesSize)
	result.OldFilesStr = FormatBytes(result.OldFilesSize)

	// Sort largest first
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Size > result.Items[j].Size
	})

	return result, nil
}

// CleanDownloadsBatch cleans specified files from ~/Downloads
func CleanDownloadsBatch(paths []string, useTrash bool) (uint64, int64, error) {
	var freed uint64
	var count int64

	for _, p := range paths {
		cleanP := filepath.Clean(p)
		fi, err := os.Stat(cleanP)
		if err != nil {
			continue
		}
		var sz uint64
		if fi.IsDir() {
			scanner := NewScanner(context.Background())
			s, _, _ := scanner.CalculateDirSize(cleanP)
			sz = s
		} else {
			sz = uint64(fi.Size())
		}

		if useTrash {
			if err := MoveToTrash(cleanP); err == nil {
				freed += sz
				count++
			}
		} else {
			if err := os.RemoveAll(cleanP); err == nil {
				freed += sz
				count++
			}
		}
	}

	return freed, count, nil
}

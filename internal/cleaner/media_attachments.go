package cleaner

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// AttachmentItem represents a single media attachment
type AttachmentItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Category  string `json:"category"` // "video", "image", "audio", "document", "other"
	Size      uint64 `json:"size"`
	SizeStr   string `json:"sizeStr"`
	ModDate   string `json:"modDate"`
	DaysOld   int    `json:"daysOld"`
	Source    string `json:"source"` // "imessage" or "mail"
}

// MediaAttachmentsResult represents the full scan result
type MediaAttachmentsResult struct {
	Items            []AttachmentItem `json:"items"`
	TotalSize        uint64           `json:"totalSize"`
	TotalSizeStr     string           `json:"totalSizeStr"`
	TotalCount       int64            `json:"totalCount"`
	VideoSize        uint64           `json:"videoSize"`
	VideoSizeStr     string           `json:"videoSizeStr"`
	ImageSize        uint64           `json:"imageSize"`
	ImageSizeStr     string           `json:"imageSizeStr"`
	AudioSize        uint64           `json:"audioSize"`
	AudioSizeStr     string           `json:"audioSizeStr"`
	DocumentSize     uint64           `json:"documentSize"`
	DocumentSizeStr  string           `json:"documentSizeStr"`
	MailSize         uint64           `json:"mailSize"`
	MailSizeStr      string           `json:"mailSizeStr"`
	MessageSize      uint64           `json:"messageSize"`
	MessageSizeStr   string           `json:"messageSizeStr"`
}

func getAttachmentCategory(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".mov", ".mp4", ".m4v", ".avi", ".mkv", ".webm":
		return "video"
	case ".jpg", ".jpeg", ".png", ".heic", ".gif", ".webp", ".tiff", ".bmp":
		return "image"
	case ".m4a", ".mp3", ".wav", ".caf", ".aac", ".flac", ".ogg":
		return "audio"
	case ".pdf", ".zip", ".rar", ".7z", ".doc", ".docx", ".pages", ".xlsx", ".numbers", ".txt":
		return "document"
	default:
		return "other"
	}
}

// ScanMediaAttachments scans iMessage and Apple Mail downloads
func ScanMediaAttachments(ctx context.Context) (*MediaAttachmentsResult, error) {
	home, _ := os.UserHomeDir()
	now := time.Now()

	result := &MediaAttachmentsResult{}

	// 1. Scan iMessage Attachments
	msgDir := filepath.Join(home, "Library", "Messages", "Attachments")
	if _, err := os.Stat(msgDir); err == nil {
		_ = filepath.WalkDir(msgDir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if d.IsDir() {
				return nil
			}

			fi, err := d.Info()
			if err != nil || fi.Size() == 0 {
				return nil
			}

			sz := uint64(fi.Size())
			ext := filepath.Ext(p)
			cat := getAttachmentCategory(ext)
			modTime := fi.ModTime()
			daysOld := int(now.Sub(modTime).Hours() / 24)

			item := AttachmentItem{
				ID:       p,
				Name:     filepath.Base(p),
				Path:     p,
				Category: cat,
				Size:     sz,
				SizeStr:  FormatBytes(sz),
				ModDate:  modTime.Format("02.01.2006 15:04"),
				DaysOld:  daysOld,
				Source:   "imessage",
			}

			result.Items = append(result.Items, item)
			result.TotalSize += sz
			result.TotalCount++
			result.MessageSize += sz

			switch cat {
			case "video":
				result.VideoSize += sz
			case "image":
				result.ImageSize += sz
			case "audio":
				result.AudioSize += sz
			case "document":
				result.DocumentSize += sz
			}

			return nil
		})
	}

	// 2. Scan Mail Downloads
	mailDir := filepath.Join(home, "Library", "Containers", "com.apple.mail", "Data", "Library", "Mail Downloads")
	if _, err := os.Stat(mailDir); err == nil {
		_ = filepath.WalkDir(mailDir, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}

			fi, err := d.Info()
			if err != nil || fi.Size() == 0 {
				return nil
			}

			sz := uint64(fi.Size())
			ext := filepath.Ext(p)
			cat := getAttachmentCategory(ext)
			modTime := fi.ModTime()
			daysOld := int(now.Sub(modTime).Hours() / 24)

			item := AttachmentItem{
				ID:       p,
				Name:     filepath.Base(p),
				Path:     p,
				Category: cat,
				Size:     sz,
				SizeStr:  FormatBytes(sz),
				ModDate:  modTime.Format("02.01.2006 15:04"),
				DaysOld:  daysOld,
				Source:   "mail",
			}

			result.Items = append(result.Items, item)
			result.TotalSize += sz
			result.TotalCount++
			result.MailSize += sz

			switch cat {
			case "video":
				result.VideoSize += sz
			case "image":
				result.ImageSize += sz
			case "audio":
				result.AudioSize += sz
			case "document":
				result.DocumentSize += sz
			}

			return nil
		})
	}

	result.TotalSizeStr = FormatBytes(result.TotalSize)
	result.VideoSizeStr = FormatBytes(result.VideoSize)
	result.ImageSizeStr = FormatBytes(result.ImageSize)
	result.AudioSizeStr = FormatBytes(result.AudioSize)
	result.DocumentSizeStr = FormatBytes(result.DocumentSize)
	result.MessageSizeStr = FormatBytes(result.MessageSize)
	result.MailSizeStr = FormatBytes(result.MailSize)

	// Sort largest first
	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].Size > result.Items[j].Size
	})

	return result, nil
}

// CleanAttachmentsRequest defines the cleanup parameters
type CleanAttachmentsRequest struct {
	Sources       []string `json:"sources"`       // "imessage", "mail"
	Categories    []string `json:"categories"`    // "video", "image", "audio", "document", "other"
	OlderThanDays int      `json:"olderThanDays"` // e.g. 30, or 0 for all
	UseTrash      bool     `json:"useTrash"`
	SelectedPaths []string `json:"selectedPaths"` // Specific item IDs if selected
}

// CleanAttachments executes deletion of selected attachments
func CleanAttachments(req CleanAttachmentsRequest) (uint64, int64, error) {
	var pathsToDelete []string

	if len(req.SelectedPaths) > 0 {
		pathsToDelete = req.SelectedPaths
	} else {
		// Scan and filter by criteria
		res, err := ScanMediaAttachments(context.Background())
		if err != nil {
			return 0, 0, err
		}

		sourceMap := make(map[string]bool)
		for _, s := range req.Sources {
			sourceMap[s] = true
		}

		catMap := make(map[string]bool)
		for _, c := range req.Categories {
			catMap[c] = true
		}

		for _, item := range res.Items {
			if len(sourceMap) > 0 && !sourceMap[item.Source] {
				continue
			}
			if len(catMap) > 0 && !catMap[item.Category] {
				continue
			}
			if req.OlderThanDays > 0 && item.DaysOld < req.OlderThanDays {
				continue
			}
			pathsToDelete = append(pathsToDelete, item.Path)
		}
	}

	var freed uint64
	var deletedCount int64

	for _, p := range pathsToDelete {
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		sz := uint64(fi.Size())

		if req.UseTrash {
			if err := MoveToTrash(p); err == nil {
				freed += sz
				deletedCount++
			}
		} else {
			if err := os.Remove(p); err == nil {
				freed += sz
				deletedCount++
			}
		}
	}

	return freed, deletedCount, nil
}

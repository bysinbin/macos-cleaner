package cleaner

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func isProtectedSystemPath(p string) bool {
	p = filepath.Clean(p)
	home, _ := os.UserHomeDir()

	critical := []string{
		"/", "/System", "/usr", "/bin", "/sbin", "/etc", "/var", "/Library",
		"/Applications", "/private", "/dev",
		filepath.Join(home, ".ssh"),
		filepath.Join(home, ".gnupg"),
		filepath.Join(home, "Library/Keychains"),
	}

	for _, c := range critical {
		if p == c || p == filepath.Clean(c) {
			return true
		}
	}
	return false
}

// ShredItemResult represents shred outcome for a path
type ShredItemResult struct {
	Path     string `json:"path"`
	Success  bool   `json:"success"`
	Size     uint64 `json:"size"`
	SizeStr  string `json:"sizeStr"`
	Error    string `json:"error,omitempty"`
}

// ShredBatchResult holds the result of a shredding operation
type ShredBatchResult struct {
	TotalFiles   int               `json:"totalFiles"`
	TotalBytes   uint64            `json:"totalBytes"`
	TotalStr     string            `json:"totalStr"`
	Duration     string            `json:"duration"`
	Items        []ShredItemResult `json:"items"`
}

// ShredFile securely erases a single file by overwriting it multiple times before deletion
func ShredFile(filePath string, passes int) (uint64, error) {
	filePath = ExpandPath(filePath)

	// Safety protection
	if isProtectedSystemPath(filePath) {
		return 0, fmt.Errorf("güvenlik koruması: kritik sistem dosyası öğütülemez: %s", filePath)
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	if fi.IsDir() {
		return 0, fmt.Errorf("hedef bir klasör: %s", filePath)
	}

	fileSize := fi.Size()
	if fileSize == 0 {
		return 0, os.Remove(filePath)
	}

	if passes < 1 {
		passes = 1
	}
	if passes > 3 {
		passes = 3
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}

	bufSize := int64(64 * 1024) // 64KB buffer
	buf := make([]byte, bufSize)

	for pass := 1; pass <= passes; pass++ {
		_, _ = f.Seek(0, 0)
		var written int64 = 0

		for written < fileSize {
			toWrite := bufSize
			if fileSize-written < bufSize {
				toWrite = fileSize - written
			}

			if pass%2 == 1 {
				// Random bytes
				_, _ = rand.Read(buf[:toWrite])
			} else {
				// Zero-fill
				for i := range buf[:toWrite] {
					buf[i] = 0
				}
			}

			n, writeErr := f.Write(buf[:toWrite])
			if writeErr != nil {
				_ = f.Close()
				return 0, writeErr
			}
			written += int64(n)
		}
		_ = f.Sync()
	}

	// Truncate to 0
	_ = f.Truncate(0)
	_ = f.Close()

	// Rename to randomized name before deleting
	randomName := filepath.Join(filepath.Dir(filePath), fmt.Sprintf(".shredded_%d", time.Now().UnixNano()))
	_ = os.Rename(filePath, randomName)

	return uint64(fileSize), os.Remove(randomName)
}

// ShredPath securely shreds a file or directory tree
func ShredPath(targetPath string, passes int) (*ShredBatchResult, error) {
	start := time.Now()
	targetPath = ExpandPath(targetPath)

	if isProtectedSystemPath(targetPath) {
		return nil, fmt.Errorf("güvenlik koruması: kritik sistem yolu öğütülemez: %s", targetPath)
	}

	fi, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}

	result := &ShredBatchResult{}

	if !fi.IsDir() {
		sz, err := ShredFile(targetPath, passes)
		itemRes := ShredItemResult{
			Path:    targetPath,
			Success: err == nil,
			Size:    sz,
			SizeStr: FormatBytes(sz),
		}
		if err != nil {
			itemRes.Error = err.Error()
		} else {
			result.TotalFiles = 1
			result.TotalBytes = sz
		}
		result.Items = append(result.Items, itemRes)
	} else {
		// Walk and shred all files
		var files []string
		_ = filepath.Walk(targetPath, func(p string, info os.FileInfo, err error) error {
			if err == nil && info != nil && !info.IsDir() {
				files = append(files, p)
			}
			return nil
		})

		for _, f := range files {
			sz, err := ShredFile(f, passes)
			itemRes := ShredItemResult{
				Path:    f,
				Success: err == nil,
				Size:    sz,
				SizeStr: FormatBytes(sz),
			}
			if err != nil {
				itemRes.Error = err.Error()
			} else {
				result.TotalFiles++
				result.TotalBytes += sz
			}
			result.Items = append(result.Items, itemRes)
		}

		// Remove empty folders
		_ = os.RemoveAll(targetPath)
	}

	result.TotalStr = FormatBytes(result.TotalBytes)
	result.Duration = time.Since(start).Round(time.Millisecond).String()
	return result, nil
}

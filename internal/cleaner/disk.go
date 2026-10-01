package cleaner

import (
	"fmt"
	"golang.org/x/sys/unix"
)

// DiskStats represents filesystem space statistics
type DiskStats struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	TotalStr    string  `json:"totalStr"`
	FreeStr     string  `json:"freeStr"`
	UsedStr     string  `json:"usedStr"`
	MountPoint  string  `json:"mountPoint"`
}

// GetDiskStats returns storage statistics for a path (default "/")
func GetDiskStats(path string) (*DiskStats, error) {
	if path == "" {
		path = "/"
	}

	var stat unix.Statfs_t
	err := unix.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	bsize := uint64(stat.Bsize)
	total := stat.Blocks * bsize
	free := stat.Bavail * bsize
	used := total - free

	var usedPct float64
	if total > 0 {
		usedPct = (float64(used) / float64(total)) * 100.0
	}

	return &DiskStats{
		Total:       total,
		Free:        free,
		Used:        used,
		UsedPercent: usedPct,
		TotalStr:    FormatBytes(total),
		FreeStr:     FormatBytes(free),
		UsedStr:     FormatBytes(used),
		MountPoint:  path,
	}, nil
}

// FormatBytes converts bytes into human readable format (KB, MB, GB, TB)
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

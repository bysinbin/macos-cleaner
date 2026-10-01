package cleaner

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// APFSSnapshot represents a local Time Machine snapshot
type APFSSnapshot struct {
	ID        string `json:"id"`
	DateStr   string `json:"dateStr"`
	SnapshotID string `json:"snapshotId"`
}

// iOSBackupItem represents an iPhone or iPad local backup
type iOSBackupItem struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	DeviceName  string `json:"deviceName"`
	DeviceModel string `json:"deviceModel"`
	LastDate    string `json:"lastDate"`
	Size        uint64 `json:"size"`
	SizeStr     string `json:"sizeStr"`
	FileCount   int64  `json:"fileCount"`
}

// SimulatorItem represents an iOS Simulator device or cache
type SimulatorItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Runtime   string `json:"runtime"`
	Size      uint64 `json:"size"`
	SizeStr   string `json:"sizeStr"`
	FileCount int64  `json:"fileCount"`
}

// AppleSystemResult encapsulates Apple & System Data breakdown
type AppleSystemResult struct {
	Snapshots       []APFSSnapshot   `json:"snapshots"`
	SnapshotsCount  int              `json:"snapshotsCount"`
	Backups         []iOSBackupItem  `json:"backups"`
	TotalBackupSize uint64           `json:"totalBackupSize"`
	TotalBackupStr  string           `json:"totalBackupStr"`
	Simulators      []SimulatorItem  `json:"simulators"`
	TotalSimSize    uint64           `json:"totalSimSize"`
	TotalSimStr     string           `json:"totalSimStr"`
	TotalSize       uint64           `json:"totalSize"`
	TotalSizeStr    string           `json:"totalSizeStr"`
}

// ScanAppleSystem scans for APFS snapshots, iOS device backups, and simulators
func ScanAppleSystem(ctx context.Context) (*AppleSystemResult, error) {
	scanner := NewScanner(ctx)
	home, _ := os.UserHomeDir()

	result := &AppleSystemResult{}

	// 1. APFS Local Snapshots via tmutil
	cmd := exec.CommandContext(ctx, "tmutil", "listlocalsnapshots", "/")
	if out, err := cmd.Output(); err == nil {
		scannerBuf := bufio.NewScanner(bytes.NewReader(out))
		dateRegex := regexp.MustCompile(`com\.apple\.TimeMachine\.(\d{4}-\d{2}-\d{2}-\d{6})\.local`)
		for scannerBuf.Scan() {
			line := strings.TrimSpace(scannerBuf.Text())
			if m := dateRegex.FindStringSubmatch(line); len(m) > 1 {
				result.Snapshots = append(result.Snapshots, APFSSnapshot{
					ID:         line,
					DateStr:    m[1],
					SnapshotID: m[1],
				})
			}
		}
	}
	result.SnapshotsCount = len(result.Snapshots)

	// 2. iOS Backups in ~/Library/Application Support/MobileSync/Backup
	backupDir := filepath.Join(home, "Library", "Application Support", "MobileSync", "Backup")
	if entries, err := os.ReadDir(backupDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			folderPath := filepath.Join(backupDir, e.Name())
			sz, count, _ := scanner.CalculateDirSize(folderPath)
			if sz == 0 {
				continue
			}

			deviceName := "Apple Aygıt Yedeklemesi"
			deviceModel := "iOS Cihazı"
			lastDate := ""

			// Read Info.plist inside backup
			infoPlist := filepath.Join(folderPath, "Info.plist")
			if data, err := os.ReadFile(infoPlist); err == nil {
				reName := regexp.MustCompile(`(?s)<key>Device Name</key>\s*<string>([^<]+)</string>`)
				reModel := regexp.MustCompile(`(?s)<key>Product Type</key>\s*<string>([^<]+)</string>`)
				reDate := regexp.MustCompile(`(?s)<key>Last Backup Date</key>\s*<date>([^<]+)</date>`)

				if m := reName.FindSubmatch(data); len(m) > 1 {
					deviceName = string(m[1])
				}
				if m := reModel.FindSubmatch(data); len(m) > 1 {
					deviceModel = string(m[1])
				}
				if m := reDate.FindSubmatch(data); len(m) > 1 {
					if t, err := time.Parse(time.RFC3339, string(m[1])); err == nil {
						lastDate = t.Format("02.01.2006 15:04")
					}
				}
			}

			if lastDate == "" {
				if fi, err := e.Info(); err == nil {
					lastDate = fi.ModTime().Format("02.01.2006 15:04")
				}
			}

			item := iOSBackupItem{
				ID:          folderPath,
				Path:        folderPath,
				DeviceName:  deviceName,
				DeviceModel: deviceModel,
				LastDate:    lastDate,
				Size:        sz,
				SizeStr:     FormatBytes(sz),
				FileCount:   count,
			}

			result.Backups = append(result.Backups, item)
			result.TotalBackupSize += sz
		}
	}
	result.TotalBackupStr = FormatBytes(result.TotalBackupSize)

	// 3. iOS Simulator Devices in ~/Library/Developer/CoreSimulator/Devices
	simDir := filepath.Join(home, "Library", "Developer", "CoreSimulator", "Devices")
	if entries, err := os.ReadDir(simDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			full := filepath.Join(simDir, e.Name())
			sz, count, _ := scanner.CalculateDirSize(full)
			if sz < 10*1024 { // Ignore tiny empty folders < 10KB
				continue
			}

			name := "Simülatör Cihazı"
			runtime := "iOS"

			devPlist := filepath.Join(full, "device.plist")
			if data, err := os.ReadFile(devPlist); err == nil {
				reDevName := regexp.MustCompile(`(?s)<key>name</key>\s*<string>([^<]+)</string>`)
				reDevType := regexp.MustCompile(`(?s)<key>deviceType</key>\s*<string>([^<]+)</string>`)
				if m := reDevName.FindSubmatch(data); len(m) > 1 {
					name = string(m[1])
				}
				if m := reDevType.FindSubmatch(data); len(m) > 1 {
					runtime = filepath.Base(string(m[1]))
				}
			}

			item := SimulatorItem{
				ID:        full,
				Name:      name,
				Path:      full,
				Runtime:   runtime,
				Size:      sz,
				SizeStr:   FormatBytes(sz),
				FileCount: count,
			}
			result.Simulators = append(result.Simulators, item)
			result.TotalSimSize += sz
		}
	}
	result.TotalSimStr = FormatBytes(result.TotalSimSize)

	result.TotalSize = result.TotalBackupSize + result.TotalSimSize
	result.TotalSizeStr = FormatBytes(result.TotalSize)

	return result, nil
}

// DeleteAPFSSnapshot deletes a specific local Time Machine snapshot
func DeleteAPFSSnapshot(snapshotDate string) error {
	cmd := exec.Command("tmutil", "deletelocalsnapshots", snapshotDate)
	return cmd.Run()
}

// DeleteAllAPFSSnapshots deletes all local APFS snapshots
func DeleteAllAPFSSnapshots() (int, error) {
	result, err := ScanAppleSystem(context.Background())
	if err != nil {
		return 0, err
	}
	deleted := 0
	for _, snap := range result.Snapshots {
		if err := DeleteAPFSSnapshot(snap.SnapshotID); err == nil {
			deleted++
		}
	}
	return deleted, nil
}

// CleanSimulators prunes unavailable simulators and simulator caches
func CleanSimulators() error {
	// 1. Delete unavailable simulator runtimes via xcrun
	_ = exec.Command("xcrun", "simctl", "delete", "unavailable").Run()

	// 2. Clear Simulator Caches
	home, _ := os.UserHomeDir()
	simCaches := []string{
		filepath.Join(home, "Library", "Caches", "com.apple.CoreSimulator.CoreSimulatorService"),
		filepath.Join(home, "Library", "Developer", "CoreSimulator", "Caches"),
	}
	for _, p := range simCaches {
		_ = os.RemoveAll(p)
	}

	return nil
}

package cleaner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// APFSSnapshot represents a local Time Machine snapshot
type APFSSnapshot struct {
	ID         string `json:"id"`
	DateStr    string `json:"dateStr"`
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

// VolumeDetail represents an APFS volume part of macOS Signed System Volume
type VolumeDetail struct {
	Name       string `json:"name"`
	Role       string `json:"role"`
	MountPoint string `json:"mountPoint"`
	Size       uint64 `json:"size"`
	SizeStr    string `json:"sizeStr"`
	IsSealed   bool   `json:"isSealed"`
}

// MacOSVolumeInfo details the read-only Signed System Volume
type MacOSVolumeInfo struct {
	TotalSize    uint64         `json:"totalSize"`
	TotalSizeStr string         `json:"totalSizeStr"`
	Status       string         `json:"status"` // "Mühürlü ve Salt Okunur (Apple SSV)"
	Description  string         `json:"description"`
	Version      string         `json:"version"` // e.g. "macOS 26.7.1 (Derleme 25G241)"
	IsReadOnly   bool           `json:"isReadOnly"`
	Volumes      []VolumeDetail `json:"volumes"`
}

// SystemDataItem represents an individual cleanable or manageable system data entry
type SystemDataItem struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"categoryId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Resolved    string    `json:"resolved"`
	Size        uint64    `json:"size"`
	SizeStr     string    `json:"sizeStr"`
	FileCount   int64     `json:"fileCount"`
	Risk        RiskLevel `json:"risk"` // "safe", "recommended", "caution"
	Cleanable   bool      `json:"cleanable"`
	ActionType  string    `json:"actionType"` // "remove", "empty_contents", "brew_cleanup", "clean_containers_cache", "reveal_only"
}

// SystemDataCategory groups related system data components
type SystemDataCategory struct {
	ID           string           `json:"id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Icon         string           `json:"icon"` // "package_managers", "app_support", "developer", "system_frameworks", "containers", "logs"
	TotalSize    uint64           `json:"totalSize"`
	TotalSizeStr string           `json:"totalSizeStr"`
	SafeSize     uint64           `json:"safeSize"`
	SafeSizeStr  string           `json:"safeSizeStr"`
	Items        []SystemDataItem `json:"items"`
}

// AppleSystemResult encapsulates Apple & System Data breakdown
type AppleSystemResult struct {
	Snapshots            []APFSSnapshot       `json:"snapshots"`
	SnapshotsCount       int                  `json:"snapshotsCount"`
	Backups              []iOSBackupItem      `json:"backups"`
	TotalBackupSize      uint64               `json:"totalBackupSize"`
	TotalBackupStr       string               `json:"totalBackupStr"`
	Simulators           []SimulatorItem      `json:"simulators"`
	TotalSimSize         uint64               `json:"totalSimSize"`
	TotalSimStr          string               `json:"totalSimStr"`
	TotalSize            uint64               `json:"totalSize"`
	TotalSizeStr         string               `json:"totalSizeStr"`
	MacOSInfo            *MacOSVolumeInfo     `json:"macOSInfo"`
	SystemDataCategories []SystemDataCategory `json:"systemDataCategories"`
	TotalSystemDataSize  uint64               `json:"totalSystemDataSize"`
	TotalSystemDataStr   string               `json:"totalSystemDataStr"`
	SafeCleanableSize    uint64               `json:"safeCleanableSize"`
	SafeCleanableStr     string               `json:"safeCleanableStr"`
	VMSleepimageSize     uint64               `json:"vmSleepimageSize"`
	VMSleepimageStr      string               `json:"vmSleepimageStr"`
	HibernateMode        string               `json:"hibernateMode"`
	BrewPath             string               `json:"brewPath,omitempty"`
	BrewCleanupReclaim   string               `json:"brewCleanupReclaim,omitempty"`
}

// ScanAppleSystem scans for APFS snapshots, iOS backups, simulators, macOS SSV, and full System Data breakdown
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

	// 4. macOS Signed System Volume (SSV) Details
	result.MacOSInfo = getMacOSVolumeInfo(ctx)

	// 5. Sleepimage & Hibernate Mode Details
	sleepimagePath := "/private/var/vm/sleepimage"
	if fi, err := os.Stat(sleepimagePath); err == nil {
		result.VMSleepimageSize = uint64(fi.Size())
		result.VMSleepimageStr = FormatBytes(result.VMSleepimageSize)
	}
	if out, err := exec.CommandContext(ctx, "pmset", "-g").Output(); err == nil {
		re := regexp.MustCompile(`hibernatemode\s+(\d+)`)
		if m := re.FindStringSubmatch(string(out)); len(m) > 1 {
			result.HibernateMode = m[1]
		}
	}

	// 6. Homebrew detection
	brewPath := GetBrewPath()
	result.BrewPath = brewPath
	if brewPath != "" {
		_, reclaimStr := GetBrewCleanupSize()
		result.BrewCleanupReclaim = reclaimStr
	}

	// 7. Comprehensive System Data (Sistem Verileri) Scanning
	categories, totalSysData, safeCleanable := scanSystemDataBreakdown(ctx, scanner, home)
	result.SystemDataCategories = categories
	result.TotalSystemDataSize = totalSysData
	result.TotalSystemDataStr = FormatBytes(totalSysData)
	result.SafeCleanableSize = safeCleanable
	result.SafeCleanableStr = FormatBytes(safeCleanable)

	return result, nil
}

// getMacOSVolumeInfo parses diskutil apfs list and sw_vers to discover the Signed System Volume structure
func getMacOSVolumeInfo(ctx context.Context) *MacOSVolumeInfo {
	info := &MacOSVolumeInfo{
		Status:      "🔒 Mühürlü ve Salt Okunur (Apple SSV)",
		Description: "macOS Big Sur ve sonrasında Apple, çekirdek işletim sistemini şifreli ve mühürlü 'Signed System Volume (SSV)' içine yerleştirmiştir. Sistem Integrity Protection (SIP) ve donanımsal kök güvenliğiyle kilitlidir. Sistem çökmesini önlemek ve güvenliği korumak amacıyla salt okunurdur; üçüncü parti yazılımlar veya kullanıcılar tarafından doğrudan silinemez / temizlenemez.",
		IsReadOnly:  true,
	}

	// Fetch macOS version
	if out, err := exec.CommandContext(ctx, "sw_vers").Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		var prodName, prodVer, buildVer string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "ProductName:") {
				prodName = strings.TrimSpace(strings.TrimPrefix(line, "ProductName:"))
			} else if strings.HasPrefix(line, "ProductVersion:") {
				prodVer = strings.TrimSpace(strings.TrimPrefix(line, "ProductVersion:"))
			} else if strings.HasPrefix(line, "BuildVersion:") {
				buildVer = strings.TrimSpace(strings.TrimPrefix(line, "BuildVersion:"))
			}
		}
		if prodVer != "" {
			info.Version = fmt.Sprintf("%s %s (Derleme %s)", prodName, prodVer, buildVer)
		}
	}
	if info.Version == "" {
		info.Version = "macOS"
	}

	// Fetch APFS volume sizes via diskutil apfs list
	cmd := exec.CommandContext(ctx, "diskutil", "apfs", "list")
	out, err := cmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		var currentRole, currentName, currentMount string
		var currentBytes uint64
		var currentSealed bool

		reRole := regexp.MustCompile(`APFS Volume Disk \(Role\):\s+\S+\s+\(([^)]+)\)`)
		reName := regexp.MustCompile(`Name:\s+([^(]+)`)
		reCap := regexp.MustCompile(`Capacity Consumed:\s+(\d+)\s+B`)
		reSealed := regexp.MustCompile(`Sealed:\s+(\w+)`)
		reMount := regexp.MustCompile(`Mount Point:\s+(.+)`)

		flushVolume := func() {
			if currentRole != "" && currentBytes > 0 {
				roleLower := strings.ToLower(currentRole)
				if roleLower == "system" || roleLower == "preboot" || roleLower == "recovery" || roleLower == "vm" {
					vd := VolumeDetail{
						Name:       currentName,
						Role:       currentRole,
						MountPoint: currentMount,
						Size:       currentBytes,
						SizeStr:    FormatBytes(currentBytes),
						IsSealed:   currentSealed,
					}
					info.Volumes = append(info.Volumes, vd)
					info.TotalSize += currentBytes
				}
			}
			currentRole = ""
			currentName = ""
			currentMount = ""
			currentBytes = 0
			currentSealed = false
		}

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if m := reRole.FindStringSubmatch(line); len(m) > 1 {
				flushVolume()
				currentRole = m[1]
			}
			if m := reName.FindStringSubmatch(line); len(m) > 1 && currentRole != "" && currentName == "" {
				currentName = strings.TrimSpace(m[1])
			}
			if m := reMount.FindStringSubmatch(line); len(m) > 1 && currentRole != "" {
				currentMount = strings.TrimSpace(m[1])
			}
			if m := reCap.FindStringSubmatch(line); len(m) > 1 && currentRole != "" {
				currentBytes, _ = strconv.ParseUint(m[1], 10, 64)
			}
			if m := reSealed.FindStringSubmatch(line); len(m) > 1 && currentRole != "" {
				currentSealed = (m[1] == "Yes")
			}
		}
		flushVolume()
	}

	if info.TotalSize == 0 {
		info.TotalSize = 24 * 1024 * 1024 * 1024 // ~24 GB default
	}
	info.TotalSizeStr = FormatBytes(info.TotalSize)

	return info
}

type targetSpec struct {
	id          string
	categoryID  string
	name        string
	description string
	path        string
	risk        RiskLevel
	cleanable   bool
	actionType  string
}

// scanSystemDataBreakdown discovers all real System Data (Sistem Verileri) locations
func scanSystemDataBreakdown(ctx context.Context, scanner *Scanner, home string) ([]SystemDataCategory, uint64, uint64) {
	var totalAll uint64
	var safeAll uint64

	topCellar := GetHomebrewCellarSummary(home)
	cellarDesc := "Kullanıcı dizinindeki aktif Homebrew paket yöneticisi ortamı, Cellar paket ikilileri ve formül depoları."
	if len(topCellar) > 0 {
		cellarDesc = fmt.Sprintf("Aktif Homebrew ortamı (Cellar). En büyük paketler: %s.", strings.Join(topCellar, ", "))
	}

	// Curated prominent system data specs
	specs := []targetSpec{
		// 1. Package Managers & CLI Tools (Homebrew, npm, etc.)
		{
			id:          "sysdata-homebrew-user",
			categoryID:  "package_managers",
			name:        "Homebrew Kullanıcı Kurulumu & Paketleri (~/homebrew)",
			description: cellarDesc,
			path:        filepath.Join(home, "homebrew"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-homebrew-opt",
			categoryID:  "package_managers",
			name:        "Apple Silicon Sistem Homebrew (/opt/homebrew)",
			description: "Apple Silicon mimarisine ait sistem Homebrew kurulum dizini ve kütüphaneleri.",
			path:        "/opt/homebrew",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-usr-local",
			categoryID:  "package_managers",
			name:        "Yerel CLI Araçları & Kütüphaneler (/usr/local)",
			description: "Intel Homebrew veya sistem düzeyinde derlenmiş ikili araçlar (/usr/local/bin, /usr/local/share, /usr/local/Cellar).",
			path:        "/usr/local",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-homebrew-cache",
			categoryID:  "package_managers",
			name:        "Homebrew İndirme Önbelleği (Bottles)",
			description: "İndirilmiş paket kurulum arşivleri ve şişeler (~/Library/Caches/Homebrew).",
			path:        filepath.Join(home, "Library", "Caches", "Homebrew"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-npm-cache",
			categoryID:  "package_managers",
			name:        "npm Paket İndirme Önbelleği",
			description: "Node.js npm paket indirme önbelleği (~/.npm).",
			path:        filepath.Join(home, ".npm"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-gradle-caches",
			categoryID:  "package_managers",
			name:        "Gradle Bağımlılık & Derleme Önbellekleri",
			description: "İndirilen Gradle kütüphaneleri ve derleme önbellekleri (~/.gradle/caches).",
			path:        filepath.Join(home, ".gradle", "caches"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-cargo-cache",
			categoryID:  "package_managers",
			name:        "Rust Cargo Paket İndirme Önbelleği",
			description: "Rust kütüphane indirme önbelleği (~/.cargo/registry/cache).",
			path:        filepath.Join(home, ".cargo", "registry", "cache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},

		// 2. Application Support: Games & Apps
		{
			id:          "sysdata-steam-games",
			categoryID:  "app_support",
			name:        "Steam Oyun Kütüphanesi & Verileri",
			description: "Steam tarafından indirilen oyun dosyaları ve yerel depolar (~/Library/Application Support/Steam/steamapps/common).",
			path:        filepath.Join(home, "Library", "Application Support", "Steam", "steamapps", "common"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-steam-bundle",
			categoryID:  "app_support",
			name:        "Steam İstemci Uygulama Dosyaları",
			description: "Steam istemcisinin çalışma ikilileri ve güncellemeleri (~/Library/Application Support/Steam/Steam.AppBundle).",
			path:        filepath.Join(home, "Library", "Application Support", "Steam", "Steam.AppBundle"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-steam-cache",
			categoryID:  "app_support",
			name:        "Steam İndirme Kalıntıları & Önbellekleri",
			description: "Steam'in geçici indirme ve istemci önbellekleri (appcache, downloading, temp).",
			path:        filepath.Join(home, "Library", "Application Support", "Steam", "appcache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-steam-shadercache",
			categoryID:  "app_support",
			name:        "Steam GPU Shader Önbelleği",
			description: "Oyunların grafik sürücüsü ve GPU shader derleme kalıntıları.",
			path:        filepath.Join(home, "Library", "Application Support", "Steam", "shadercache"),
			risk:        RiskRecommended,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-epicgames-shared",
			categoryID:  "app_support",
			name:        "Epic Games Paylaşılan Oyun Kütüphanesi (/Users/Shared)",
			description: "Epic Games Launcher tarafından indirilen bağımsız oyunlar (TRDeliveryService: 4.3 GB, Torchlight2: 1.5 GB). macOS bu alanı kullanıcı profili dışında olduğundan doğrudan Sistem Verisi olarak sayar.",
			path:        "/Users/Shared/Epic Games",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-adobe-temp",
			categoryID:  "app_support",
			name:        "Adobe Geçici Yükleme Arşivleri (.adobeTemp)",
			description: "Adobe Creative Cloud yükleyicisinin disk kökünde bıraktığı geçici yükleme arşivleri (/System/Volumes/Data/.adobeTemp).",
			path:        "/System/Volumes/Data/.adobeTemp",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-adobe-system-appsupport",
			categoryID:  "app_support",
			name:        "Adobe Kök Sistem Kütüphaneleri & Destek",
			description: "Adobe uygulamalarının paylaşılan sistem bileşenleri ve CameraRaw profilleri (/Library/Application Support/Adobe).",
			path:        "/Library/Application Support/Adobe",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-adobe-common",
			categoryID:  "app_support",
			name:        "Adobe Ortak Medya Önbellekleri",
			description: "Adobe Premiere, After Effects ve Photoshop'un ortak medya ve render önbellekleri.",
			path:        filepath.Join(home, "Library", "Application Support", "Adobe", "Common"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-adobe-cameraraw-cache",
			categoryID:  "app_support",
			name:        "Adobe CameraRaw Önbelleği",
			description: "Camera Raw önizleme ve fotoğraf önbellek verileri.",
			path:        filepath.Join(home, "Library", "Application Support", "Adobe", "CameraRaw", "Cache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-blackmagic",
			categoryID:  "app_support",
			name:        "Blackmagic Design & DaVinci Destek Verileri",
			description: "DaVinci Resolve ve Blackmagic araçları sistem destek dosyaları (/Library/Application Support/Blackmagic Design).",
			path:        "/Library/Application Support/Blackmagic Design",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-google-appsupport",
			categoryID:  "app_support",
			name:        "Google Chrome / Chromium Profil & Uygulama Verileri",
			description: "Google Chrome kullanıcı profilleri, uzantıları ve veritabanları (~/Library/Application Support/Google).",
			path:        filepath.Join(home, "Library", "Application Support", "Google"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-chrome-shader",
			categoryID:  "app_support",
			name:        "Google Chrome Shader & Webview Önbelleği",
			description: "Chrome GPU shader derlemeleri ve WebGL/GrShaderCache önbellekleri.",
			path:        filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "ShaderCache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-chrome-serviceworker",
			categoryID:  "app_support",
			name:        "Google Chrome Service Worker Önbelleği",
			description: "Web sitelerinin arka planda depoladığı çevrimdışı ve geçici veriler.",
			path:        filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default", "Service Worker", "CacheStorage"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-vscode-appsupport",
			categoryID:  "app_support",
			name:        "VS Code Uygulama Verileri & Profilleri",
			description: "VS Code kullanıcı çalışma alanı durumu, yapılandırmaları ve oturumları (~/Library/Application Support/Code).",
			path:        filepath.Join(home, "Library", "Application Support", "Code"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-vscode-cacheddata",
			categoryID:  "app_support",
			name:        "VS Code Derleme & Webview Önbelleği",
			description: "VS Code tarafından oluşturulan V8 bytecode ve geçici çalışma önbellekleri (CachedData, Cache).",
			path:        filepath.Join(home, "Library", "Application Support", "Code", "CachedData"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-visualstudio-appsupport",
			categoryID:  "app_support",
			name:        "Visual Studio for Mac Uygulama Verileri",
			description: "Visual Studio masaüstü ortamı ve proje çalışma verileri (~/Library/Application Support/VisualStudio).",
			path:        filepath.Join(home, "Library", "Application Support", "VisualStudio"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-battlenet-appsupport",
			categoryID:  "app_support",
			name:        "Battle.net Oyun İstemcisi Deposu",
			description: "Blizzard Battle.net istemci dosyaları ve indirme verileri (~/Library/Application Support/Battle.net).",
			path:        filepath.Join(home, "Library", "Application Support", "Battle.net"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-antigravity-appsupport",
			categoryID:  "app_support",
			name:        "Antigravity AI Ajanı Oturum Depoları",
			description: "Antigravity yapay zeka ajan oturumu kayıtları ve çalışma verileri (~/Library/Application Support/Antigravity).",
			path:        filepath.Join(home, "Library", "Application Support", "Antigravity"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-discord-cache",
			categoryID:  "app_support",
			name:        "Discord Medya ve Görsel Önbelleği",
			description: "Discord kanallarından indirilen profil resimleri, gifler ve medya önbellekleri.",
			path:        filepath.Join(home, "Library", "Application Support", "discord", "Cache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-techsmith",
			categoryID:  "app_support",
			name:        "TechSmith (Camtasia / Snagit) Medya Deposu",
			description: "Camtasia kayıtları, geçici render dosyaları ve medya havuzları (~/Library/Application Support/TechSmith).",
			path:        filepath.Join(home, "Library", "Application Support", "TechSmith"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-eagle",
			categoryID:  "app_support",
			name:        "Eagle Görsel Arşivleme Kütüphanesi",
			description: "Eagle tasarım ve görsel kütüphane verileri (~/Library/Application Support/Eagle).",
			path:        filepath.Join(home, "Library", "Application Support", "Eagle"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},

		// 3. Developer Runtimes, SDKs & Toolchains
		{
			id:          "sysdata-cmdlinetools",
			categoryID:  "developer",
			name:        "macOS Command Line Tools SDK",
			description: "Apple Xcode Komut Satırı Araçları SDK ve başlık dosyaları (/Library/Developer/CommandLineTools).",
			path:        "/Library/Developer/CommandLineTools",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-arduino-packages",
			categoryID:  "developer",
			name:        "Arduino Donanım Paketleri & Çekirdekleri",
			description: "Arduino IDE tarafından indirilen ESP32, AVR ve ARM kart paketleri (~/Library/Arduino15/packages).",
			path:        filepath.Join(home, "Library", "Arduino15", "packages"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-arduino-staging",
			categoryID:  "developer",
			name:        "Arduino Paket İndirme Arşivi (Staging)",
			description: "İndirilmiş Arduino kart ve çekirdek paketleri. Kurulum sonrasında güvenle temizlenebilir.",
			path:        filepath.Join(home, "Library", "Arduino15", "staging"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-arduino-cache",
			categoryID:  "developer",
			name:        "Arduino IDE Derleme ve İndirme Önbelleği",
			description: "Arduino derleyici ve kütüphane önbellek dosyaları (~/Library/Arduino15/cache).",
			path:        filepath.Join(home, "Library", "Arduino15", "cache"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-android-sdk",
			categoryID:  "developer",
			name:        "Android SDK & Platform Paketleri",
			description: "Android Studio SDK platformları, derleme araçları ve emülatör bileşenleri (~/Library/Android/sdk).",
			path:        filepath.Join(home, "Library", "Android", "sdk"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-android-user",
			categoryID:  "developer",
			name:        "Android Studio Kullanıcı Önbellekleri (~/.android)",
			description: "Android Studio derleme önbellekleri, avd kalıpları ve eklenti verileri.",
			path:        filepath.Join(home, ".android"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-vscode-extensions",
			categoryID:  "developer",
			name:        "VS Code Yüklü Eklentiler Deposu",
			description: "VS Code kullanıcı eklentileri (~/.vscode/extensions).",
			path:        filepath.Join(home, ".vscode", "extensions"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-rustup",
			categoryID:  "developer",
			name:        "Rustup Toolchain ve Bileşenleri",
			description: "Rust programlama dili derleyicileri ve platform toolchain'leri (~/.rustup).",
			path:        filepath.Join(home, ".rustup"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-postgresql",
			categoryID:  "developer",
			name:        "PostgreSQL Yerel Sunucu ve Veritabanları",
			description: "macOS yerel PostgreSQL servis dosyaları ve veritabanı kümeleri (/Library/PostgreSQL).",
			path:        "/Library/PostgreSQL",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-java",
			categoryID:  "developer",
			name:        "Java Çalışma Zamanları / JDKs",
			description: "macOS kullanıcı alanında kurulu Java SDK ve çalışma zamanı ikilileri (~/Library/Java).",
			path:        filepath.Join(home, "Library", "Java"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-go-workspace",
			categoryID:  "developer",
			name:        "Go Çalışma Alanı & İndirilen Modüller",
			description: "Go programlama dili modül önbellekleri ve ikili araçları (~/go).",
			path:        filepath.Join(home, "go"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-go-build-cache",
			categoryID:  "developer",
			name:        "Go Derleme Önbelleği (go-build)",
			description: "Go derleyicisinin ara derleme önbellekleri (~/Library/Caches/go-build).",
			path:        filepath.Join(home, "Library", "Caches", "go-build"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "remove",
		},
		{
			id:          "sysdata-phpwebstudy",
			categoryID:  "developer",
			name:        "PhpWebStudy Yerel Geliştirme Ortamı",
			description: "Yerel PHP/Web sunucusu ortamı ve çalışma verileri (~/Library/PhpWebStudy).",
			path:        filepath.Join(home, "Library", "PhpWebStudy"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-gemini-dot",
			categoryID:  "developer",
			name:        "Gemini CLI & Ajan Çalışma Verileri",
			description: "Gemini AI geliştirici araçları, konfigürasyonları ve eklentileri (~/.gemini).",
			path:        filepath.Join(home, ".gemini"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-antigravity-dot",
			categoryID:  "developer",
			name:        "Antigravity IDE Çalışma Alanı Verileri",
			description: "Antigravity IDE dahili çalışma dizini ve oturum logları (~/.antigravity-ide).",
			path:        filepath.Join(home, ".antigravity-ide"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-dotnet",
			categoryID:  "developer",
			name:        ".NET SDK & Çalışma Zamanları",
			description: ".NET Core ve C# yerel araçları ve çalışma paketleri (~/.net).",
			path:        filepath.Join(home, ".net"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-cmdline-tools",
			categoryID:  "developer",
			name:        "Komut Satırı SDK Araçları (~/cmdline-tools)",
			description: "Kullanıcı dizinindeki Android veya diğer CLI araçları (~/cmdline-tools).",
			path:        filepath.Join(home, "cmdline-tools"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-node-modules",
			categoryID:  "developer",
			name:        "Kullanıcı Ana Dizini node_modules",
			description: "Kullanıcı kök dizininde unutulmuş Node.js bağımlılıkları (~/node_modules).",
			path:        filepath.Join(home, "node_modules"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},

		// 4. System Frameworks, Containers & Sleepimage
		{
			id:          "sysdata-frameworks",
			categoryID:  "system_frameworks",
			name:        "macOS Paylaşılan Sistem Çerçeveleri (/Library/Frameworks)",
			description: "macOS üçüncü parti yazılımlar ve sistem servislerinin paylaştığı dinamik kütüphaneler.",
			path:        "/Library/Frameworks",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-group-containers",
			categoryID:  "system_frameworks",
			name:        "Paylaşılan Grup Kapsayıcıları (Group Containers)",
			description: "Birden fazla Mac uygulamasının ortak kullandığı sandbox veri depoları (~/Library/Group Containers).",
			path:        filepath.Join(home, "Library", "Group Containers"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-httpstorages",
			categoryID:  "system_frameworks",
			name:        "HTTPStorages Ağ & İndirme Önbellekleri",
			description: "Uygulamaların yerel ağ transfer ve önbellek depoları (~/Library/HTTPStorages).",
			path:        filepath.Join(home, "Library", "HTTPStorages"),
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-sleepimage",
			categoryID:  "system_frameworks",
			name:        "macOS Uyku & Bellek İmajı (Sleepimage)",
			description: "Mac uyku modundayken RAM belleğinin diske kopyalanan yedeği (/private/var/vm/sleepimage). Sistem veri boyutunun doğrudan parçasıdır.",
			path:        "/private/var/vm/sleepimage",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},
		{
			id:          "sysdata-var-folders",
			categoryID:  "system_frameworks",
			name:        "macOS Kullanıcı Geçici & Önbellek Ağacı (/private/var/folders)",
			description: "macOS GUI uygulamalarının, Safari webview bileşenlerinin ve sistem servislerinin izole geçici çalışma ve önbellek depoları.",
			path:        "/private/var/folders",
			risk:        RiskCaution,
			cleanable:   false,
			actionType:  "reveal_only",
		},

		// 5. System Logs, Crash Reports & Diagnostics
		{
			id:          "sysdata-user-logs",
			categoryID:  "logs",
			name:        "Kullanıcı Uygulama Günlükleri (User Logs)",
			description: "Kullanıcı kütüphanesinde biriken uygulama ve servis logları (~/Library/Logs).",
			path:        filepath.Join(home, "Library", "Logs"),
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "empty_contents",
		},
		{
			id:          "sysdata-system-logs",
			categoryID:  "logs",
			name:        "Sistem Genel Günlükleri & Çökme Raporları",
			description: "macOS hata, kilitlenme ve tanılama dökümleri (/Library/Logs & DiagnosticReports).",
			path:        "/Library/Logs",
			risk:        RiskSafe,
			cleanable:   true,
			actionType:  "empty_contents",
		},
	}

	categoryMap := map[string]*SystemDataCategory{
		"package_managers": {
			ID:          "package_managers",
			Title:       "Paket Yöneticileri & CLI Araçları (Homebrew, npm vb.)",
			Description: "Homebrew (3.9 GB+), yerel kütüphaneler (/usr/local), npm ve bağımlılık önbellekleri. Buradaki eski sürümler 'brew cleanup' ile temizlenebilir.",
			Icon:        "package_managers",
		},
		"app_support": {
			ID:          "app_support",
			Title:       "Uygulama İçi Veriler & Medya Depoları (Application Support)",
			Description: "Steam oyunları (4.3 GB), Adobe sistem bileşenleri (8.4 GB), Google Chrome, VS Code ve diğer masaüstü yazılımların büyük depoları.",
			Icon:        "app_support",
		},
		"developer": {
			ID:          "developer",
			Title:       "Geliştirici SDK, Araçlar & Çekirdek Paketleri",
			Description: "macOS Command Line Tools SDK (5.9 GB), Arduino paketleri (4.2 GB), Android Studio SDK (2.8 GB), PostgreSQL, Rustup ve Java.",
			Icon:        "developer",
		},
		"containers": {
			ID:          "containers",
			Title:       "Korumalı Uygulama Kapsayıcıları (Sandboxed Containers)",
			Description: "macOS sandbox mimarisinde çalışan izole uygulamaların (Oyunlar, Office, QuickLook, Sistem Ajanları) veri depoları ve geçici önbellekleri.",
			Icon:        "containers",
		},
		"system_frameworks": {
			ID:          "system_frameworks",
			Title:       "Sistem Çerçeveleri & Sanal Bellek (macOS Frameworks & VM)",
			Description: "macOS paylaşımlı kütüphaneleri (/Library/Frameworks - 3.4 GB) ve uyku dosyası (/private/var/vm/sleepimage - 2.15 GB).",
			Icon:        "system_frameworks",
		},
		"logs": {
			ID:          "logs",
			Title:       "Sistem Günlükleri, Çökme Raporları & Tanılama",
			Description: "Uygulama hata kayıtları, kilitlenme dökümleri ve macOS sistem logları.",
			Icon:        "logs",
		},
	}

	// Dynamic detection of Homebrew Cleanup reclaimable space
	if brewCleanupBytes, brewCleanupStr := GetBrewCleanupSize(); brewCleanupBytes > 0 {
		item := SystemDataItem{
			ID:          "sysdata-homebrew-cleanup",
			CategoryID:  "package_managers",
			Name:        "Homebrew Eski Sürüm & Artık Paketleri (brew cleanup)",
			Description: fmt.Sprintf("Homebrew tarafından otomatik kaldırılabilir eski formül sürümleri ve gereksiz bağımlılıklar (%s alan açılabilir).", brewCleanupStr),
			Path:        GetBrewPath(),
			Resolved:    GetBrewPath(),
			Size:        brewCleanupBytes,
			SizeStr:     brewCleanupStr,
			FileCount:   1,
			Risk:        RiskSafe,
			Cleanable:   true,
			ActionType:  "brew_cleanup",
		}
		cat := categoryMap["package_managers"]
		cat.Items = append(cat.Items, item)
		cat.TotalSize += brewCleanupBytes
		cat.SafeSize += brewCleanupBytes
		safeAll += brewCleanupBytes
		totalAll += brewCleanupBytes
	}

	type scanItemResult struct {
		spec  targetSpec
		size  uint64
		count int64
	}

	resChan := make(chan scanItemResult, len(specs))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 12) // Concurrency limit of 12 workers

	for _, s := range specs {
		wg.Add(1)
		go func(spec targetSpec) {
			defer wg.Done()
			if _, err := os.Stat(spec.path); err != nil {
				return
			}
			sem <- struct{}{}
			sz, count, _ := scanner.CalculateDirSize(spec.path)
			<-sem
			if sz > 0 {
				resChan <- scanItemResult{spec: spec, size: sz, count: count}
			}
		}(s)
	}

	wg.Wait()
	close(resChan)

	for r := range resChan {
		item := SystemDataItem{
			ID:          r.spec.id,
			CategoryID:  r.spec.categoryID,
			Name:        r.spec.name,
			Description: r.spec.description,
			Path:        r.spec.path,
			Resolved:    r.spec.path,
			Size:        r.size,
			SizeStr:     FormatBytes(r.size),
			FileCount:   r.count,
			Risk:        r.spec.risk,
			Cleanable:   r.spec.cleanable,
			ActionType:  r.spec.actionType,
		}

		cat := categoryMap[r.spec.categoryID]
		if cat != nil {
			cat.Items = append(cat.Items, item)
			cat.TotalSize += r.size
			if r.spec.risk == RiskSafe {
				cat.SafeSize += r.size
				safeAll += r.size
			}
			totalAll += r.size
		}
	}

	// Dynamic Sandboxed Container Caches scan (~/Library/Containers/*/Data/Library/Caches)
	containersPattern := filepath.Join(home, "Library", "Containers", "*", "Data", "Library", "Caches")
	if matches, err := filepath.Glob(containersPattern); err == nil && len(matches) > 0 {
		var containerCacheBytes uint64
		var containerCacheFiles int64
		for _, m := range matches {
			sz, count, _ := scanner.CalculateDirSize(m)
			containerCacheBytes += sz
			containerCacheFiles += count
		}

		if containerCacheBytes > 0 {
			item := SystemDataItem{
				ID:          "sysdata-containers-caches",
				CategoryID:  "containers",
				Name:        "Tüm Sandbox Kapsayıcı Önbellekleri (Containers Caches)",
				Description: fmt.Sprintf("%d farklı korumalı Mac uygulamasının (QuickLook, Wallpaper, Office vb.) geçici çalışma önbellekleri.", len(matches)),
				Path:        filepath.Join(home, "Library", "Containers"),
				Resolved:    containersPattern,
				Size:        containerCacheBytes,
				SizeStr:     FormatBytes(containerCacheBytes),
				FileCount:   containerCacheFiles,
				Risk:        RiskSafe,
				Cleanable:   true,
				ActionType:  "clean_containers_cache",
			}
			cat := categoryMap["containers"]
			cat.Items = append(cat.Items, item)
			cat.TotalSize += containerCacheBytes
			cat.SafeSize += containerCacheBytes
			totalAll += containerCacheBytes
			safeAll += containerCacheBytes
		}
	}

	// Prominent large containers to inspect
	prominentContainers := []string{
		"com.wemade.mir4global",
		"com.microsoft.rdc.macos",
		"com.utmapp.UTM",
		"com.apple.geod",
		"com.apple.mediaanalysisd",
		"com.microsoft.Powerpoint",
		"com.microsoft.Word",
		"com.microsoft.Excel",
		"com.adobe.accmac.ACCFinderSync",
	}

	for _, cName := range prominentContainers {
		p := filepath.Join(home, "Library", "Containers", cName)
		if _, err := os.Stat(p); err != nil {
			continue
		}
		sz, count, _ := scanner.CalculateDirSize(p)
		if sz > 50*1024*1024 { // Only list if > 50MB
			item := SystemDataItem{
				ID:          "sysdata-container-" + cName,
				CategoryID:  "containers",
				Name:        fmt.Sprintf("Kapsayıcı: %s", cName),
				Description: "İzole Mac uygulaması tam veri alanı.",
				Path:        p,
				Resolved:    p,
				Size:        sz,
				SizeStr:     FormatBytes(sz),
				FileCount:   count,
				Risk:        RiskCaution,
				Cleanable:   false,
				ActionType:  "reveal_only",
			}
			cat := categoryMap["containers"]
			cat.Items = append(cat.Items, item)
			cat.TotalSize += sz
			totalAll += sz
		}
	}

	// Sort items in each category from largest to smallest
	order := []string{"package_managers", "app_support", "developer", "containers", "system_frameworks", "logs"}
	var categories []SystemDataCategory
	for _, k := range order {
		cat := categoryMap[k]
		if cat != nil && len(cat.Items) > 0 {
			sort.Slice(cat.Items, func(i, j int) bool {
				return cat.Items[i].Size > cat.Items[j].Size
			})
			cat.TotalSizeStr = FormatBytes(cat.TotalSize)
			cat.SafeSizeStr = FormatBytes(cat.SafeSize)
			categories = append(categories, *cat)
		}
	}

	return categories, totalAll, safeAll
}

// GetHomebrewCellarSummary returns top packages in Cellar and their formatted strings
func GetHomebrewCellarSummary(home string) []string {
	cellarPaths := []string{
		filepath.Join(home, "homebrew", "Cellar"),
		"/opt/homebrew/Cellar",
		"/usr/local/Cellar",
	}
	for _, cp := range cellarPaths {
		entries, err := os.ReadDir(cp)
		if err != nil || len(entries) == 0 {
			continue
		}
		type pkgSize struct {
			name string
			size uint64
		}
		var pkgs []pkgSize
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			p := filepath.Join(cp, e.Name())
			var sz uint64
			_ = filepath.Walk(p, func(_ string, info os.FileInfo, err error) error {
				if err == nil && info != nil && !info.IsDir() {
					sz += uint64(info.Size())
				}
				return nil
			})
			if sz > 30*1024*1024 { // only > 30MB
				pkgs = append(pkgs, pkgSize{name: e.Name(), size: sz})
			}
		}
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].size > pkgs[j].size
		})
		var topSummary []string
		for i, p := range pkgs {
			if i >= 6 {
				break
			}
			topSummary = append(topSummary, fmt.Sprintf("%s (%s)", p.name, FormatBytes(p.size)))
		}
		if len(topSummary) > 0 {
			return topSummary
		}
	}
	return nil
}

// GetBrewPath returns the absolute path of the brew binary if found
func GetBrewPath() string {
	home, _ := os.UserHomeDir()
	candidates := []string{
		filepath.Join(home, "homebrew", "bin", "brew"),
		"/opt/homebrew/bin/brew",
		"/usr/local/bin/brew",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	if p, err := exec.LookPath("brew"); err == nil {
		return p
	}
	return ""
}

// GetBrewCleanupSize checks how much space brew cleanup would free
func GetBrewCleanupSize() (uint64, string) {
	brewPath := GetBrewPath()
	if brewPath == "" {
		return 0, ""
	}
	cmd := exec.Command(brewPath, "cleanup", "-s", "-n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, ""
	}
	re := regexp.MustCompile(`free approximately (\d+(?:\.\d+)?)\s*([KMGTP]?B)`)
	m := re.FindStringSubmatch(string(out))
	if len(m) > 2 {
		val, _ := strconv.ParseFloat(m[1], 64)
		unit := strings.ToUpper(m[2])
		var multiplier float64 = 1
		switch unit {
		case "KB":
			multiplier = 1024
		case "MB":
			multiplier = 1024 * 1024
		case "GB":
			multiplier = 1024 * 1024 * 1024
		}
		bytes := uint64(val * multiplier)
		return bytes, fmt.Sprintf("%.1f %s", val, unit)
	}
	return 0, ""
}

// ExecuteBrewCleanup runs brew cleanup and autoremove to reclaim disk space
func ExecuteBrewCleanup() (string, error) {
	brewPath := GetBrewPath()
	if brewPath == "" {
		return "", fmt.Errorf("Homebrew ikili dosyası bulunamadı")
	}
	out, err := exec.Command(brewPath, "cleanup", "-s").CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

// CleanSystemDataItem cleans a specific system data item
func CleanSystemDataItem(id string, path string, useTrash bool) (uint64, error) {
	scanner := NewScanner(nil)
	home, _ := os.UserHomeDir()

	// Special action: brew cleanup
	if id == "sysdata-homebrew-cleanup" {
		beforeSz, _ := GetBrewCleanupSize()
		_, err := ExecuteBrewCleanup()
		if err != nil {
			return 0, err
		}
		if beforeSz == 0 {
			beforeSz = 350 * 1024 * 1024 // ~350MB
		}
		return beforeSz, nil
	}

	// Handle special containers cache batch clean
	if id == "sysdata-containers-caches" {
		pattern := filepath.Join(home, "Library", "Containers", "*", "Data", "Library", "Caches")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return 0, err
		}
		var freed uint64
		for _, m := range matches {
			sz, _, _ := scanner.CalculateDirSize(m)
			entries, err := os.ReadDir(m)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				itemPath := filepath.Join(m, entry.Name())
				_ = os.RemoveAll(itemPath)
			}
			freed += sz
		}
		return freed, nil
	}

	resolved := ExpandPath(path)
	if err := IsPathSafe(resolved); err != nil {
		return 0, err
	}

	sz, _, _ := scanner.CalculateDirSize(resolved)

	// If empty contents requested (e.g. Logs directory)
	if strings.HasSuffix(id, "-logs") {
		entries, err := os.ReadDir(resolved)
		if err != nil {
			return 0, err
		}
		for _, entry := range entries {
			sub := filepath.Join(resolved, entry.Name())
			if useTrash {
				_ = MoveToTrash(sub)
			} else {
				_ = os.RemoveAll(sub)
			}
		}
		return sz, nil
	}

	// Default: remove target
	var err error
	if useTrash {
		err = MoveToTrash(resolved)
	} else {
		err = os.RemoveAll(resolved)
	}
	if err != nil {
		return 0, err
	}
	return sz, nil
}

// CleanAllSafeSystemData cleans all items marked RiskSafe across all categories
func CleanAllSafeSystemData() (uint64, int, error) {
	result, err := ScanAppleSystem(context.Background())
	if err != nil {
		return 0, 0, err
	}

	var totalFreed uint64
	var count int

	for _, cat := range result.SystemDataCategories {
		for _, item := range cat.Items {
			if item.Risk == RiskSafe && item.Cleanable {
				freed, err := CleanSystemDataItem(item.ID, item.Path, false)
				if err == nil && freed > 0 {
					totalFreed += freed
					count++
				}
			}
		}
	}

	return totalFreed, count, nil
}

// ReclaimPurgeableSpace triggers macOS purge routines and cache flushes
func ReclaimPurgeableSpace() (string, error) {
	var results []string

	// 1. Purge command (frees OS file-backed inactive cache)
	out, err := exec.Command("purge").CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	if err == nil && !strings.Contains(outStr, "not permitted") && !strings.Contains(outStr, "Unable") {
		results = append(results, "macOS purge çalıştırıldı")
	} else {
		triggerMemoryPressure()
		results = append(results, "Bellek önbellek basıncı döngüsü tetiklendi")
	}

	// 2. QuickLook thumbnail cache reset
	_ = exec.Command("qlmanage", "-r", "cache").Run()
	results = append(results, "QuickLook önbelleği sıfırlandı")

	// 3. User Font cache reset
	_ = exec.Command("atsutil", "databases", "-removeUser").Run()
	results = append(results, "Yazı tipi önbellekleri yenilendi")

	// 4. DNS cache flush
	_ = exec.Command("dscacheutil", "-flushcache").Run()
	_ = exec.Command("killall", "-HUP", "mDNSResponder").Run()
	results = append(results, "DNS önbelleği boşaltıldı")

	return strings.Join(results, ", "), nil
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

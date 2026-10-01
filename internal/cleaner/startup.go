package cleaner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// StartupItemType represents the classification of startup/background items
type StartupItemType string

const (
	ItemTypeLoginItem    StartupItemType = "login_item"
	ItemTypeUserAgent    StartupItemType = "user_agent"
	ItemTypeSystemAgent  StartupItemType = "system_agent"
	ItemTypeSystemDaemon StartupItemType = "system_daemon"
)

// StartupItem represents an individual startup program, login item, or background daemon
type StartupItem struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Label       string          `json:"label"`
	Type        StartupItemType `json:"type"`
	TypeName    string          `json:"typeName"`
	Path        string          `json:"path"`
	Program     string          `json:"program"`
	Enabled     bool            `json:"enabled"`
	Running     bool            `json:"running"`
	PID         int             `json:"pid"`
	RunAtLoad   bool            `json:"runAtLoad"`
	Vendor      string          `json:"vendor"`
	Hidden      bool            `json:"hidden"`
	Writable    bool            `json:"writable"`
	Description string          `json:"description"`
}

// StartupSummary provides aggregate statistics and the complete list of startup items
type StartupSummary struct {
	TotalCount        int           `json:"totalCount"`
	ActiveCount       int           `json:"activeCount"`
	RunningCount      int           `json:"runningCount"`
	LoginItemsCount   int           `json:"loginItemsCount"`
	UserAgentsCount   int           `json:"userAgentsCount"`
	SystemAgentsCount int           `json:"systemAgentsCount"`
	DaemonsCount      int           `json:"daemonsCount"`
	Items             []StartupItem `json:"items"`
}

// PlistData represents the JSON converted structure of a launchd plist
type PlistData struct {
	Label            string   `json:"Label"`
	Program          string   `json:"Program"`
	ProgramArguments []string `json:"ProgramArguments"`
	RunAtLoad        *bool    `json:"RunAtLoad"`
	Disabled         *bool    `json:"Disabled"`
}

// DisabledLoginRecord stores disabled user login items persistently
type DisabledLoginRecord struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Hidden bool   `json:"hidden"`
}

// ScanStartupItems collects all login items, user launch agents, system launch agents, and daemons
func ScanStartupItems(ctx context.Context) (*StartupSummary, error) {
	summary := &StartupSummary{
		Items: make([]StartupItem, 0),
	}

	// 1. Collect running launchctl jobs (label -> pid)
	runningJobs := getRunningLaunchdJobs()

	// 2. Collect launchctl disabled states
	disabledJobs := getLaunchctlDisabledServices()

	// 3. Scan User Login Items (via AppleScript System Events + persistent disabled store)
	loginItems := scanUserLoginItems(ctx)
	summary.Items = append(summary.Items, loginItems...)

	// 4. Scan User LaunchAgents (~/Library/LaunchAgents)
	home, _ := os.UserHomeDir()
	userAgentsDir := filepath.Join(home, "Library", "LaunchAgents")
	userAgents := scanLaunchDirectory(userAgentsDir, ItemTypeUserAgent, "Kullanıcı Ajanı", runningJobs, disabledJobs)
	summary.Items = append(summary.Items, userAgents...)

	// 5. Scan System LaunchAgents (/Library/LaunchAgents)
	sysAgents := scanLaunchDirectory("/Library/LaunchAgents", ItemTypeSystemAgent, "Sistem Ajanı", runningJobs, disabledJobs)
	summary.Items = append(summary.Items, sysAgents...)

	// 6. Scan System LaunchDaemons (/Library/LaunchDaemons)
	sysDaemons := scanLaunchDirectory("/Library/LaunchDaemons", ItemTypeSystemDaemon, "Arka Plan Hizmeti", runningJobs, disabledJobs)
	summary.Items = append(summary.Items, sysDaemons...)

	// Calculate counts
	summary.TotalCount = len(summary.Items)
	for _, item := range summary.Items {
		if item.Enabled {
			summary.ActiveCount++
		}
		if item.Running {
			summary.RunningCount++
		}
		switch item.Type {
		case ItemTypeLoginItem:
			summary.LoginItemsCount++
		case ItemTypeUserAgent:
			summary.UserAgentsCount++
		case ItemTypeSystemAgent:
			summary.SystemAgentsCount++
		case ItemTypeSystemDaemon:
			summary.DaemonsCount++
		}
	}

	return summary, nil
}

// getDisabledLoginFilePath returns the path to ~/.disk-cleaner-disabled-login-items.json
func getDisabledLoginFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".disk-cleaner-disabled-login-items.json"
	}
	return filepath.Join(home, ".disk-cleaner-disabled-login-items.json")
}

// loadDisabledLoginRecords loads persistent disabled login items
func loadDisabledLoginRecords() map[string]DisabledLoginRecord {
	records := make(map[string]DisabledLoginRecord)
	path := getDisabledLoginFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return records
	}
	var list []DisabledLoginRecord
	if err := json.Unmarshal(data, &list); err == nil {
		for _, rec := range list {
			records[rec.Name] = rec
		}
	}
	return records
}

// saveDisabledLoginRecords writes persistent disabled login items to disk
func saveDisabledLoginRecords(records map[string]DisabledLoginRecord) {
	list := make([]DisabledLoginRecord, 0, len(records))
	for _, rec := range records {
		list = append(list, rec)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err == nil {
		_ = os.WriteFile(getDisabledLoginFilePath(), data, 0644)
	}
}

// scanUserLoginItems queries macOS System Events for registered login items and merges disabled list
func scanUserLoginItems(ctx context.Context) []StartupItem {
	items := make([]StartupItem, 0)
	activeNames := make(map[string]bool)

	script := `
tell application "System Events"
    set outList to ""
    set allItems to every login item
    repeat with anItem in allItems
        set itemName to name of anItem
        set itemPath to ""
        try
            set itemPath to path of anItem
        end try
        set itemHidden to hidden of anItem
        set outList to outList & itemName & "||" & itemPath & "||" & itemHidden & "\n"
    end repeat
    return outList
end tell
`
	cmd := exec.CommandContext(ctx, "osascript", "-e", script)
	out, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Split(line, "||")
			if len(parts) < 3 {
				continue
			}
			name := strings.TrimSpace(parts[0])
			path := strings.TrimSpace(parts[1])
			hiddenStr := strings.TrimSpace(parts[2])

			if path == "missing value" {
				path = ""
			}

			hidden := hiddenStr == "true"
			vendor := inferVendorFromPathOrName(path, name)
			running := isProcessRunning(name, path)

			activeNames[name] = true

			items = append(items, StartupItem{
				ID:          "login:" + name,
				Name:        name,
				Label:       name,
				Type:        ItemTypeLoginItem,
				TypeName:    "Oturum Açma Öğesi",
				Path:        path,
				Program:     path,
				Enabled:     true,
				Running:     running,
				RunAtLoad:   true,
				Hidden:      hidden,
				Vendor:      vendor,
				Writable:    true,
				Description: "Kullanıcı oturum açtığında otomatik başlatılan masaüstü uygulaması",
			})
		}
	}

	// Now merge with persistent disabled login items
	disabledRecords := loadDisabledLoginRecords()
	cleanedRecords := make(map[string]DisabledLoginRecord)

	for name, rec := range disabledRecords {
		if activeNames[name] {
			// If it's already active in System Events, don't list as disabled
			continue
		}
		cleanedRecords[name] = rec
		vendor := inferVendorFromPathOrName(rec.Path, rec.Name)
		running := isProcessRunning(rec.Name, rec.Path)

		items = append(items, StartupItem{
			ID:          "login:" + rec.Name,
			Name:        rec.Name,
			Label:       rec.Name,
			Type:        ItemTypeLoginItem,
			TypeName:    "Oturum Açma Öğesi (Devre Dışı)",
			Path:        rec.Path,
			Program:     rec.Path,
			Enabled:     false,
			Running:     running,
			RunAtLoad:   false,
			Hidden:      rec.Hidden,
			Vendor:      vendor,
			Writable:    true,
			Description: "Geçici olarak devre dışı bırakılmış oturum açma öğesi",
		})
	}

	if len(cleanedRecords) != len(disabledRecords) {
		saveDisabledLoginRecords(cleanedRecords)
	}

	return items
}

// getLaunchctlDisabledServices parses `launchctl print-disabled` outputs
// Returns a map where label -> true means disabled, false means explicitly enabled
func getLaunchctlDisabledServices() map[string]bool {
	disabledMap := make(map[string]bool)

	// 1. User GUI domain
	uid := os.Getuid()
	if out, err := exec.Command("launchctl", "print-disabled", fmt.Sprintf("gui/%d", uid)).Output(); err == nil {
		parsePrintDisabledOutput(string(out), disabledMap)
	}

	// 2. System domain
	if out, err := exec.Command("launchctl", "print-disabled", "system").Output(); err == nil {
		parsePrintDisabledOutput(string(out), disabledMap)
	}

	return disabledMap
}

func parsePrintDisabledOutput(out string, dest map[string]bool) {
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		// e.g. "com.adobe.AdobeCreativeCloud" => enabled
		// or "com.adobe.GC.Scheduler-1.0" => disabled
		if strings.Contains(l, "=>") {
			parts := strings.Split(l, "=>")
			if len(parts) == 2 {
				key := strings.Trim(strings.TrimSpace(parts[0]), `"`)
				val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				if val == "disabled" {
					dest[key] = true
				} else if val == "enabled" {
					dest[key] = false
				}
			}
		}
	}
}

// scanLaunchDirectory reads a directory of launchd plists (.plist or .plist.disabled)
func scanLaunchDirectory(dir string, itemType StartupItemType, typeName string, runningJobs map[string]int, disabledJobs map[string]bool) []StartupItem {
	items := make([]StartupItem, 0)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || strings.HasPrefix(name, ".") {
			continue
		}

		isDisabledFile := strings.HasSuffix(name, ".plist.disabled") || strings.HasSuffix(name, ".disabled")
		isPlist := strings.HasSuffix(name, ".plist")
		if !isPlist && !isDisabledFile {
			continue
		}

		fullPath := filepath.Join(dir, name)
		data, err := parsePlistViaPlutil(fullPath)
		if err != nil {
			continue
		}

		label := data.Label
		if label == "" {
			label = strings.TrimSuffix(strings.TrimSuffix(name, ".disabled"), ".plist")
		}

		program := data.Program
		if program == "" && len(data.ProgramArguments) > 0 {
			program = strings.Join(data.ProgramArguments, " ")
		}

		runAtLoad := true
		if data.RunAtLoad != nil {
			runAtLoad = *data.RunAtLoad
		}

		// Calculate enabled state:
		// 1. If file has .disabled extension -> false
		// 2. If launchctl print-disabled explicitly says disabled -> false
		// 3. If plist internal Disabled is true -> false
		// 4. Otherwise -> true
		enabled := !isDisabledFile
		if data.Disabled != nil && *data.Disabled {
			enabled = false
		}
		if isDisabledInLaunchd, exists := disabledJobs[label]; exists {
			if isDisabledInLaunchd {
				enabled = false
			} else if !isDisabledFile {
				enabled = true
			}
		}

		pid, running := runningJobs[label]
		vendor := inferVendorFromLabel(label, program)
		displayName := cleanDisplayName(label, name)
		writable := isFileWritable(fullPath)

		var desc string
		switch itemType {
		case ItemTypeUserAgent:
			desc = "Kullanıcı oturumu başladığında arka planda çalışan başlatıcı ajan"
		case ItemTypeSystemAgent:
			desc = "Sistem çapında oturum açıldığında çalışan servis ajanı"
		case ItemTypeSystemDaemon:
			desc = "macOS açılışında sistem düzeyinde çalışan arka plan servisi (Daemon)"
		}

		items = append(items, StartupItem{
			ID:          string(itemType) + ":" + name,
			Name:        displayName,
			Label:       label,
			Type:        itemType,
			TypeName:    typeName,
			Path:        fullPath,
			Program:     program,
			Enabled:     enabled,
			Running:     running,
			PID:         pid,
			RunAtLoad:   runAtLoad,
			Vendor:      vendor,
			Writable:    writable,
			Description: desc,
		})
	}

	return items
}

// parsePlistViaPlutil converts a plist into JSON via plutil -convert json
func parsePlistViaPlutil(path string) (*PlistData, error) {
	cmd := exec.Command("plutil", "-convert", "json", "-o", "-", path)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data PlistData
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// getRunningLaunchdJobs queries `launchctl list`
func getRunningLaunchdJobs() map[string]int {
	running := make(map[string]int)
	cmd := exec.Command("launchctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return running
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			pidStr := fields[0]
			label := fields[2]
			if pidStr != "-" {
				if pid, err := strconv.Atoi(pidStr); err == nil {
					running[label] = pid
				}
			}
		}
	}
	return running
}

// ToggleStartupItem enables or disables a startup item
func ToggleStartupItem(id string, enable bool) error {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("geçersiz öğe kimliği: %s", id)
	}

	itemType := parts[0]
	target := parts[1]

	switch itemType {
	case "login":
		disabledRecords := loadDisabledLoginRecords()

		if !enable {
			// Disable: find active item in System Events to get its path, then delete it and record to persistent store
			findScript := fmt.Sprintf(`
tell application "System Events"
    set foundPath to ""
    set foundHidden to false
    try
        set theItem to login item "%s"
        try
            set foundPath to path of theItem
        end try
        set foundHidden to hidden of theItem
        delete theItem
    end try
    return foundPath & "||" & foundHidden
end tell
`, target)
			out, err := exec.Command("osascript", "-e", findScript).Output()
			if err != nil {
				return fmt.Errorf("oturum açma öğesi silinemedi: %w", err)
			}
			parts := strings.Split(strings.TrimSpace(string(out)), "||")
			path := ""
			hidden := false
			if len(parts) >= 1 && parts[0] != "missing value" {
				path = parts[0]
			}
			if len(parts) >= 2 && parts[1] == "true" {
				hidden = true
			}

			// If path is missing, check if we already had a record
			if path == "" {
				if oldRec, ok := disabledRecords[target]; ok {
					path = oldRec.Path
					hidden = oldRec.Hidden
				}
			}

			disabledRecords[target] = DisabledLoginRecord{
				Name:   target,
				Path:   path,
				Hidden: hidden,
			}
			saveDisabledLoginRecords(disabledRecords)
			return nil

		} else {
			// Enable: retrieve from persistent store and add back to System Events
			rec, ok := disabledRecords[target]
			if !ok || rec.Path == "" {
				// Try using target as app name if in /Applications
				potentialPath := filepath.Join("/Applications", target+".app")
				if _, err := os.Stat(potentialPath); err == nil {
					rec.Path = potentialPath
				} else {
					return fmt.Errorf("etkinleştirmek için uygulama yolu bulunamadı: %s", target)
				}
			}

			addScript := fmt.Sprintf(`tell application "System Events" to make login item at end with properties {path:"%s", hidden:%t}`, rec.Path, rec.Hidden)
			cmd := exec.Command("osascript", "-e", addScript)
			if out, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("oturum açma öğesi eklenemedi: %s (%w)", string(out), err)
			}

			delete(disabledRecords, target)
			saveDisabledLoginRecords(disabledRecords)
			return nil
		}

	case "user_agent", "system_agent", "system_daemon":
		home, _ := os.UserHomeDir()
		var dir string
		switch itemType {
		case "user_agent":
			dir = filepath.Join(home, "Library", "LaunchAgents")
		case "system_agent":
			dir = "/Library/LaunchAgents"
		case "system_daemon":
			dir = "/Library/LaunchDaemons"
		}

		currentPath := filepath.Join(dir, target)
		if _, err := os.Stat(currentPath); err != nil {
			// Check if exists with other suffix (.plist or .disabled)
			if strings.HasSuffix(target, ".disabled") {
				alt := strings.TrimSuffix(target, ".disabled")
				if _, errAlt := os.Stat(filepath.Join(dir, alt)); errAlt == nil {
					currentPath = filepath.Join(dir, alt)
					target = alt
				} else {
					return fmt.Errorf("öğe dosyası bulunamadı: %s", currentPath)
				}
			} else if strings.HasSuffix(target, ".plist") {
				alt := currentPath + ".disabled"
				if _, errAlt := os.Stat(alt); errAlt == nil {
					currentPath = alt
					target = target + ".disabled"
				} else {
					return fmt.Errorf("öğe dosyası bulunamadı: %s", currentPath)
				}
			} else {
				return fmt.Errorf("öğe dosyası bulunamadı: %s", currentPath)
			}
		}

		// Parse plist to get true label
		label := strings.TrimSuffix(strings.TrimSuffix(target, ".disabled"), ".plist")
		if data, err := parsePlistViaPlutil(currentPath); err == nil && data.Label != "" {
			label = data.Label
		}

		uid := os.Getuid()

		if itemType == "system_daemon" {
			// System Daemons require system domain
			action := "disable"
			if enable {
				action = "enable"
			}
			cmd := exec.Command("launchctl", action, "system/"+label)
			if err := cmd.Run(); err != nil {
				// Fallback to osascript admin privileges
				adminScript := fmt.Sprintf(`do shell script "launchctl %s system/%s" with administrator privileges`, action, label)
				if errAdm := exec.Command("osascript", "-e", adminScript).Run(); errAdm != nil {
					return fmt.Errorf("sistem servisi değiştirilemedi: %w", errAdm)
				}
			}
			return nil
		}

		// User Agent & System Agent run in gui domain
		if enable {
			// launchctl enable gui/<uid>/<label>
			_ = exec.Command("launchctl", "enable", fmt.Sprintf("gui/%d/%s", uid, label)).Run()
			_ = exec.Command("launchctl", "bootstrap", fmt.Sprintf("gui/%d", uid), currentPath).Run()

			// If file was renamed to .disabled and we have write permission, restore .plist
			if strings.HasSuffix(target, ".disabled") && isFileWritable(currentPath) {
				newPath := strings.TrimSuffix(currentPath, ".disabled")
				_ = os.Rename(currentPath, newPath)
			}
			return nil
		} else {
			// Disable: bootout & disable via launchctl
			_ = exec.Command("launchctl", "bootout", fmt.Sprintf("gui/%d/%s", uid, label)).Run()
			cmd := exec.Command("launchctl", "disable", fmt.Sprintf("gui/%d/%s", uid, label))
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("launchctl servisi devre dışı bırakılamadı: %w", err)
			}

			// If file is writable and ends in .plist, optionally rename to .disabled
			if strings.HasSuffix(target, ".plist") && isFileWritable(currentPath) {
				newPath := currentPath + ".disabled"
				_ = os.Rename(currentPath, newPath)
			}
			return nil
		}

	default:
		return fmt.Errorf("bilinmeyen öğe türü: %s", itemType)
	}
}

// RemoveStartupItem removes a login item or moves a launch agent plist to Trash
func RemoveStartupItem(id string, useTrash bool) error {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("geçersiz öğe kimliği: %s", id)
	}

	itemType := parts[0]
	target := parts[1]

	switch itemType {
	case "login":
		// Remove from System Events and from disabled store
		script := fmt.Sprintf(`tell application "System Events" to delete login item "%s"`, target)
		_ = exec.Command("osascript", "-e", script).Run()

		disabledRecords := loadDisabledLoginRecords()
		delete(disabledRecords, target)
		saveDisabledLoginRecords(disabledRecords)
		return nil

	case "user_agent", "system_agent", "system_daemon":
		home, _ := os.UserHomeDir()
		var dir string
		switch itemType {
		case "user_agent":
			dir = filepath.Join(home, "Library", "LaunchAgents")
		case "system_agent":
			dir = "/Library/LaunchAgents"
		case "system_daemon":
			dir = "/Library/LaunchDaemons"
		}

		plistPath := filepath.Join(dir, target)
		if _, err := os.Stat(plistPath); err != nil {
			return fmt.Errorf("öğe dosyası bulunamadı: %s", plistPath)
		}

		// First unload / bootout
		label := strings.TrimSuffix(strings.TrimSuffix(target, ".disabled"), ".plist")
		uid := os.Getuid()
		if itemType == "system_daemon" {
			_ = exec.Command("launchctl", "bootout", "system/"+label).Run()
		} else {
			_ = exec.Command("launchctl", "bootout", fmt.Sprintf("gui/%d/%s", uid, label)).Run()
		}
		_ = exec.Command("launchctl", "unload", "-w", plistPath).Run()

		// If user has write permission, try Finder trash or direct remove
		if isFileWritable(plistPath) {
			if useTrash {
				script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, plistPath)
				if err := exec.Command("osascript", "-e", script).Run(); err == nil {
					return nil
				}
			}
			if err := os.Remove(plistPath); err == nil {
				return nil
			}
		}

		// Fallback to administrator privileges for root-owned files in /Library
		adminScript := fmt.Sprintf(`do shell script "rm -f " & quoted form of "%s" with administrator privileges`, plistPath)
		cmd := exec.Command("osascript", "-e", adminScript)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("dosya silinemedi (yönetici izni hatası): %s (%w)", string(out), err)
		}
		return nil

	default:
		return fmt.Errorf("bilinmeyen öğe türü: %s", itemType)
	}
}

// AddLoginItem adds a new application to user login items
func AddLoginItem(appPath string, hidden bool) error {
	if _, err := os.Stat(appPath); err != nil {
		return fmt.Errorf("uygulama yolu bulunamadı: %s", appPath)
	}

	script := fmt.Sprintf(`tell application "System Events" to make login item at end with properties {path:"%s", hidden:%t}`, appPath, hidden)
	cmd := exec.Command("osascript", "-e", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("başlangıç öğesi eklenemedi: %s", string(out))
	}
	return nil
}

// Helpers

func isProcessRunning(name, path string) bool {
	if name == "" {
		return false
	}
	cmd := exec.Command("pgrep", "-f", name)
	if err := cmd.Run(); err == nil {
		return true
	}
	if path != "" {
		base := filepath.Base(path)
		cmd = exec.Command("pgrep", "-f", base)
		if err := cmd.Run(); err == nil {
			return true
		}
	}
	return false
}

func isFileWritable(path string) bool {
	return syscall.Access(path, 2) == nil // 2 is W_OK
}

func inferVendorFromLabel(label, program string) string {
	l := strings.ToLower(label)
	p := strings.ToLower(program)

	switch {
	case strings.Contains(l, "adobe") || strings.Contains(p, "adobe"):
		return "Adobe"
	case strings.Contains(l, "microsoft") || strings.Contains(p, "microsoft") || strings.Contains(l, "onedrive"):
		return "Microsoft"
	case strings.Contains(l, "google") || strings.Contains(p, "google"):
		return "Google"
	case strings.Contains(l, "apple") || strings.Contains(p, "/system/"):
		return "Apple"
	case strings.Contains(l, "valvesoftware") || strings.Contains(l, "steam"):
		return "Valve (Steam)"
	case strings.Contains(l, "macpaw") || strings.Contains(l, "cleanmymac"):
		return "MacPaw"
	case strings.Contains(l, "homebrew"):
		return "Homebrew"
	case strings.Contains(l, "cloudflare") || strings.Contains(l, "warp"):
		return "Cloudflare"
	case strings.Contains(l, "freedownloadmanager") || strings.Contains(l, "fdm"):
		return "Free Download Manager"
	case strings.Contains(l, "tunnelblick"):
		return "Tunnelblick"
	case strings.Contains(l, "cindori") || strings.Contains(l, "sensei"):
		return "Cindori (Sensei)"
	case strings.Contains(l, "oracle") || strings.Contains(l, "java"):
		return "Oracle"
	case strings.Contains(l, "bitgapp") || strings.Contains(l, "eqmac"):
		return "eqMac"
	case strings.Contains(l, "flyenv"):
		return "FlyEnv"
	case strings.Contains(l, "wireshark"):
		return "Wireshark"
	case strings.Contains(l, "mongodb"):
		return "MongoDB"
	case strings.Contains(l, "postgres"):
		return "PostgreSQL"
	default:
		parts := strings.Split(label, ".")
		if len(parts) >= 2 && parts[0] == "com" {
			return strings.Title(parts[1])
		}
		return "Üçüncü Parti"
	}
}

func inferVendorFromPathOrName(path, name string) string {
	p := strings.ToLower(path)
	n := strings.ToLower(name)

	switch {
	case strings.Contains(p, "keyboard maestro") || strings.Contains(n, "keyboard maestro"):
		return "Stairways Software"
	case strings.Contains(p, "lightshot") || strings.Contains(n, "lightshot"):
		return "Skillbrains"
	case strings.Contains(p, "touchswitcher") || strings.Contains(n, "touchswitcher"):
		return "TouchSwitcher"
	case strings.Contains(p, "spotify") || strings.Contains(n, "spotify"):
		return "Spotify"
	case strings.Contains(p, "discord") || strings.Contains(n, "discord"):
		return "Discord"
	case strings.Contains(p, "slack") || strings.Contains(n, "slack"):
		return "Slack"
	case strings.Contains(p, "docker") || strings.Contains(n, "docker"):
		return "Docker"
	case strings.Contains(p, "telegram") || strings.Contains(n, "telegram"):
		return "Telegram"
	case strings.Contains(p, "whatsapp") || strings.Contains(n, "whatsapp"):
		return "WhatsApp"
	case strings.Contains(p, "dropbox") || strings.Contains(n, "dropbox"):
		return "Dropbox"
	case strings.Contains(p, "adobe"):
		return "Adobe"
	case strings.Contains(p, "google"):
		return "Google"
	case strings.Contains(p, "/system/"):
		return "Apple"
	default:
		if path != "" {
			base := filepath.Base(path)
			return strings.TrimSuffix(base, ".app")
		}
		return "Uygulama"
	}
}

func cleanDisplayName(label, filename string) string {
	name := strings.TrimSuffix(strings.TrimSuffix(filename, ".disabled"), ".plist")
	parts := strings.Split(name, ".")
	if len(parts) > 2 {
		last := parts[len(parts)-1]
		if len(last) > 2 {
			var b bytes.Buffer
			for i, r := range last {
				if i == 0 {
					b.WriteString(strings.ToUpper(string(r)))
				} else {
					b.WriteRune(r)
				}
			}
			return b.String() + " (" + name + ")"
		}
	}
	return name
}

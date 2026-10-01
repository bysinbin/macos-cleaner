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

// ScanStartupItems collects all login items, user launch agents, system launch agents, and daemons
func ScanStartupItems(ctx context.Context) (*StartupSummary, error) {
	summary := &StartupSummary{
		Items: make([]StartupItem, 0),
	}

	// 1. Collect running launchctl jobs (label -> pid)
	runningJobs := getRunningLaunchdJobs()

	// 2. Scan User Login Items (via AppleScript System Events)
	loginItems := scanUserLoginItems(ctx)
	summary.Items = append(summary.Items, loginItems...)

	// 3. Scan User LaunchAgents (~/Library/LaunchAgents)
	home, _ := os.UserHomeDir()
	userAgentsDir := filepath.Join(home, "Library", "LaunchAgents")
	userAgents := scanLaunchDirectory(userAgentsDir, ItemTypeUserAgent, "Kullanıcı Ajanı", runningJobs)
	summary.Items = append(summary.Items, userAgents...)

	// 4. Scan System LaunchAgents (/Library/LaunchAgents)
	sysAgents := scanLaunchDirectory("/Library/LaunchAgents", ItemTypeSystemAgent, "Sistem Ajanı", runningJobs)
	summary.Items = append(summary.Items, sysAgents...)

	// 5. Scan System LaunchDaemons (/Library/LaunchDaemons)
	sysDaemons := scanLaunchDirectory("/Library/LaunchDaemons", ItemTypeSystemDaemon, "Arka Plan Hizmeti", runningJobs)
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

// scanUserLoginItems queries macOS System Events for registered login items
func scanUserLoginItems(ctx context.Context) []StartupItem {
	items := make([]StartupItem, 0)

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
	if err != nil {
		return items
	}

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

	return items
}

// scanLaunchDirectory reads a directory of launchd plists (.plist or .plist.disabled)
func scanLaunchDirectory(dir string, itemType StartupItemType, typeName string, runningJobs map[string]int) []StartupItem {
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

		isDisabledFile := strings.HasSuffix(name, ".plist.disabled")
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

		enabled := !isDisabledFile
		if data.Disabled != nil && *data.Disabled {
			enabled = false
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
		// For login items, toggle hidden state
		script := fmt.Sprintf(`tell application "System Events" to set hidden of login item "%s" to %t`, target, !enable)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()

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
			return fmt.Errorf("öğe dosyası bulunamadı: %s", currentPath)
		}

		if enable {
			// If it's .disabled, rename to .plist and launchctl load
			if strings.HasSuffix(target, ".disabled") {
				newPath := strings.TrimSuffix(currentPath, ".disabled")
				if err := os.Rename(currentPath, newPath); err != nil {
					return fmt.Errorf("dosya yeniden adlandırılamadı: %w", err)
				}
				_ = exec.Command("launchctl", "load", "-w", newPath).Run()
			} else {
				_ = exec.Command("launchctl", "load", "-w", currentPath).Run()
			}
		} else {
			// Disable: launchctl unload and rename to .disabled
			_ = exec.Command("launchctl", "unload", "-w", currentPath).Run()
			if strings.HasSuffix(target, ".plist") {
				newPath := currentPath + ".disabled"
				if err := os.Rename(currentPath, newPath); err != nil {
					return fmt.Errorf("dosya devre dışı bırakılamadı: %w", err)
				}
			}
		}
		return nil

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
		script := fmt.Sprintf(`tell application "System Events" to delete login item "%s"`, target)
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Run()

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

		// First unload
		_ = exec.Command("launchctl", "unload", "-w", plistPath).Run()

		if useTrash {
			script := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, plistPath)
			if err := exec.Command("osascript", "-e", script).Run(); err == nil {
				return nil
			}
		}

		// Fallback to direct remove
		return os.Remove(plistPath)

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
		// e.g. com.valvesoftware.steamclean -> Steamclean
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

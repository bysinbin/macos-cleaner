package cleaner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PermissionStatus represents macOS privacy and security permission states
type PermissionStatus struct {
	HasFullDiskAccess bool              `json:"hasFullDiskAccess"`
	FDAStatusMessage  string            `json:"fdaStatusMessage"`
	SettingsURL       string            `json:"settingsUrl"`
	Platform          string            `json:"platform"`
	UserHome          string            `json:"userHome"`
	Diagnostics       map[string]string `json:"diagnostics,omitempty"`
}

// CheckFullDiskAccess verifies if the current process has macOS Full Disk Access (FDA)
// by attempting to read multiple TCC-protected system directories and files.
func CheckFullDiskAccess() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	dirTargets := []string{
		filepath.Join(home, "Library", "Safari"),
		filepath.Join(home, "Library", "Messages"),
		filepath.Join(home, "Library", "Mail"),
		filepath.Join(home, "Library", "Suggestions"),
		filepath.Join(home, "Library", "HomeKit"),
		filepath.Join(home, "Library", "Cookies"),
		filepath.Join(home, "Library", "PersonalizationPortrait"),
		filepath.Join(home, "Library", "Containers", "com.apple.mail", "Data", "Library", "Mail"),
	}

	for _, dir := range dirTargets {
		if _, err := os.ReadDir(dir); err == nil {
			return true
		}
	}

	fileTargets := []string{
		filepath.Join(home, "Library", "Safari", "Bookmarks.plist"),
		filepath.Join(home, "Library", "Safari", "CloudTabs.db"),
		filepath.Join(home, "Library", "Messages", "chat.db"),
		filepath.Join(home, "Library", "Suggestions", "suggestions.db"),
		filepath.Join(home, "Library", "Cookies", "Cookies.binarycookies"),
		"/Library/Preferences/com.apple.TimeMachine.plist",
	}

	for _, file := range fileTargets {
		if f, err := os.Open(file); err == nil {
			_ = f.Close()
			return true
		}
	}

	return false
}

// ProbeFDADiagnostics returns detailed probe results for each TCC target
func ProbeFDADiagnostics() map[string]string {
	home, _ := os.UserHomeDir()
	results := make(map[string]string)

	targets := []string{
		filepath.Join(home, "Library", "Safari"),
		filepath.Join(home, "Library", "Safari", "Bookmarks.plist"),
		filepath.Join(home, "Library", "Messages"),
		filepath.Join(home, "Library", "Messages", "chat.db"),
		filepath.Join(home, "Library", "Mail"),
		filepath.Join(home, "Library", "Cookies"),
		filepath.Join(home, "Library", "Suggestions"),
		"/Library/Preferences/com.apple.TimeMachine.plist",
	}

	for _, t := range targets {
		fi, err := os.Stat(t)
		if err != nil {
			results[t] = err.Error()
			continue
		}
		if fi.IsDir() {
			if entries, err := os.ReadDir(t); err != nil {
				results[t] = "readdir: " + err.Error()
			} else {
				results[t] = fmt.Sprintf("OK (dir, %d items)", len(entries))
			}
		} else {
			if f, err := os.Open(t); err != nil {
				results[t] = "open: " + err.Error()
			} else {
				_ = f.Close()
				results[t] = "OK (file readable)"
			}
		}
	}
	return results
}

// GetSystemPermissions gathers current permission status
func GetSystemPermissions() PermissionStatus {
	fda := CheckFullDiskAccess()
	home, _ := os.UserHomeDir()
	diag := ProbeFDADiagnostics()

	msg := "Tam Disk Erişimi (FDA) etkin. Sistem ve korumalı dizinler eksiksiz taranabilir."
	if !fda {
		msg = "Tam Disk Erişimi (FDA) eksik. Safari, Mail, Mesajlar ve Time Machine dizinlerini tam taramak için Sistem Ayarları'ndan izin verilmesi önerilir."
	}

	return PermissionStatus{
		HasFullDiskAccess: fda,
		FDAStatusMessage:  msg,
		SettingsURL:       "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles",
		Platform:          "darwin",
		UserHome:          home,
		Diagnostics:       diag,
	}
}

// OpenFDASettings opens macOS System Settings directly at the Full Disk Access pane
func OpenFDASettings() error {
	_ = exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles").Run()
	_ = exec.Command("osascript", "-e", `tell application "System Settings" to activate`).Run()
	return nil
}

// PromptForAdmin executes a shell command with native macOS Touch ID / password prompt
func PromptForAdmin(promptMsg, bashCmd string) (string, error) {
	script := strings.ReplaceAll(bashCmd, `"`, `\"`)
	prompt := strings.ReplaceAll(promptMsg, `"`, `\"`)
	appleScript := `do shell script "` + script + `" with prompt "` + prompt + `" with administrator privileges`
	
	cmd := exec.Command("osascript", "-e", appleScript)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

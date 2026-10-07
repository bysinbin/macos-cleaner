package cleaner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PermissionStatus represents macOS privacy and security permission states
type PermissionStatus struct {
	HasFullDiskAccess bool   `json:"hasFullDiskAccess"`
	FDAStatusMessage  string `json:"fdaStatusMessage"`
	SettingsURL       string `json:"settingsUrl"`
	Platform          string `json:"platform"`
	UserHome          string `json:"userHome"`
}

// CheckFullDiskAccess verifies if the current process has macOS Full Disk Access (FDA)
// by attempting to read protected system directories (Safari or Messages).
func CheckFullDiskAccess() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	// Target 1: Safari directory / Bookmarks (TCC protected)
	safariPath := filepath.Join(home, "Library", "Safari")
	if entries, err := os.ReadDir(safariPath); err == nil && len(entries) > 0 {
		return true
	}

	// Target 2: Messages chat database directory (TCC protected)
	messagesPath := filepath.Join(home, "Library", "Messages")
	if entries, err := os.ReadDir(messagesPath); err == nil && len(entries) > 0 {
		return true
	}

	// Target 3: Mail container
	mailPath := filepath.Join(home, "Library", "Mail")
	if entries, err := os.ReadDir(mailPath); err == nil && len(entries) > 0 {
		return true
	}

	return false
}

// GetSystemPermissions gathers current permission status
func GetSystemPermissions() PermissionStatus {
	fda := CheckFullDiskAccess()
	home, _ := os.UserHomeDir()

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
	}
}

// OpenFDASettings opens macOS System Settings directly at the Full Disk Access pane
func OpenFDASettings() error {
	cmd := exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles")
	return cmd.Run()
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

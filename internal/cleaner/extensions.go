package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// ExtensionType represents the category of the extension
type ExtensionType string

const (
	ExtTypePluginKit ExtensionType = "pluginkit"
	ExtTypeSystem    ExtensionType = "system_ext"
	ExtTypeQuickLook ExtensionType = "quicklook"
	ExtTypeSpotlight ExtensionType = "spotlight"
	ExtTypePrefPane  ExtensionType = "prefpane"
)

// ExtensionItem represents an installed system or app extension
type ExtensionItem struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	BundleID    string        `json:"bundleId"`
	Type        ExtensionType `json:"type"`
	TypeName    string        `json:"typeName"`
	Path        string        `json:"path"`
	Version     string        `json:"version"`
	TeamID      string        `json:"teamId"`
	Enabled     bool          `json:"enabled"`
	Active      bool          `json:"active"`
	Vendor      string        `json:"vendor"`
	Description string        `json:"description"`
	Writable    bool          `json:"writable"`
	IsApple     bool          `json:"isApple"`
}

// ExtensionsSummary provides counts and lists of all extensions
type ExtensionsSummary struct {
	TotalCount     int             `json:"totalCount"`
	ActiveCount    int             `json:"activeCount"`
	PluginKitCount int             `json:"pluginKitCount"`
	SystemExtCount int             `json:"systemExtCount"`
	QuickLookCount int             `json:"quickLookCount"`
	SpotlightCount int             `json:"spotlightCount"`
	PrefPaneCount  int             `json:"prefPaneCount"`
	Items          []ExtensionItem `json:"items"`
}

// ScanExtensions gathers all installed extensions across macOS
func ScanExtensions(ctx context.Context) (*ExtensionsSummary, error) {
	summary := &ExtensionsSummary{
		Items: make([]ExtensionItem, 0),
	}

	// 1. Scan System Extensions (driver & network extensions via systemextensionsctl)
	sysExts := scanSystemExtensionsCtl(ctx)
	summary.Items = append(summary.Items, sysExts...)

	// 2. Scan PluginKit App Extensions (Finder Sync, Share, Action, Widgets)
	pluginKitExts := scanPluginKitExtensions(ctx)
	summary.Items = append(summary.Items, pluginKitExts...)

	// 3. Scan QuickLook Generators (/Library/QuickLook & ~/Library/QuickLook)
	home, _ := os.UserHomeDir()
	qlDirs := []string{"/Library/QuickLook", filepath.Join(home, "Library", "QuickLook")}
	for _, dir := range qlDirs {
		qlItems := scanBundleDirectory(dir, ExtTypeQuickLook, "QuickLook Önizleme", ".qlgenerator")
		summary.Items = append(summary.Items, qlItems...)
	}

	// 4. Scan Spotlight Importers (/Library/Spotlight & ~/Library/Spotlight)
	spotlightDirs := []string{"/Library/Spotlight", filepath.Join(home, "Library", "Spotlight")}
	for _, dir := range spotlightDirs {
		spItems := scanBundleDirectory(dir, ExtTypeSpotlight, "Spotlight Arama Aktarıcı", ".mdimporter")
		summary.Items = append(summary.Items, spItems...)
	}

	// 5. Scan Preference Panes (/Library/PreferencePanes & ~/Library/PreferencePanes)
	prefDirs := []string{"/Library/PreferencePanes", filepath.Join(home, "Library", "PreferencePanes")}
	for _, dir := range prefDirs {
		prefItems := scanBundleDirectory(dir, ExtTypePrefPane, "Sistem Tercih Paneli", ".prefPane")
		summary.Items = append(summary.Items, prefItems...)
	}

	// Calculate counts
	summary.TotalCount = len(summary.Items)
	for _, item := range summary.Items {
		if item.Enabled {
			summary.ActiveCount++
		}
		switch item.Type {
		case ExtTypeSystem:
			summary.SystemExtCount++
		case ExtTypePluginKit:
			summary.PluginKitCount++
		case ExtTypeQuickLook:
			summary.QuickLookCount++
		case ExtTypeSpotlight:
			summary.SpotlightCount++
		case ExtTypePrefPane:
			summary.PrefPaneCount++
		}
	}

	return summary, nil
}

// scanSystemExtensionsCtl parses output of `systemextensionsctl list`
func scanSystemExtensionsCtl(ctx context.Context) []ExtensionItem {
	items := make([]ExtensionItem, 0)
	cmd := exec.CommandContext(ctx, "systemextensionsctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return items
	}

	lines := strings.Split(string(out), "\n")
	category := "Sürücü / Ağ Eklentisi"

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "driver_extension") {
			category = "Sistem Sürücü Eklentisi (DriverKit)"
			continue
		} else if strings.Contains(trimmed, "network_extension") {
			category = "Ağ Filtresi & VPN Eklentisi (NetworkExtension)"
			continue
		}

		if strings.HasPrefix(trimmed, "---") || strings.HasPrefix(trimmed, "enabled\t") || strings.HasPrefix(trimmed, "enabled ") || strings.Contains(trimmed, "extension(s)") {
			continue
		}

		// Row format: enabled \t active \t teamID \t bundleID (version) \t name \t [state]
		parts := strings.Split(line, "\t")
		if len(parts) < 5 {
			continue
		}

		enabledMark := strings.TrimSpace(parts[0])
		activeMark := strings.TrimSpace(parts[1])
		teamID := strings.TrimSpace(parts[2])
		bundleVer := strings.TrimSpace(parts[3])
		name := strings.TrimSpace(parts[4])
		state := ""
		if len(parts) >= 6 {
			state = strings.TrimSpace(parts[5])
		}

		enabled := enabledMark == "*" || strings.Contains(state, "enabled")
		active := activeMark == "*" || strings.Contains(state, "activated")

		bundleID := bundleVer
		version := ""
		if idx := strings.Index(bundleVer, " ("); idx != -1 {
			bundleID = strings.TrimSpace(bundleVer[:idx])
			version = strings.TrimSuffix(strings.TrimSpace(bundleVer[idx+2:]), ")")
		}

		vendor := inferVendorFromBundleID(bundleID, name)
		desc := fmt.Sprintf("%s — Durum: %s", category, state)

		items = append(items, ExtensionItem{
			ID:          "sysext:" + bundleID,
			Name:        name,
			BundleID:    bundleID,
			Type:        ExtTypeSystem,
			TypeName:    category,
			Path:        "",
			Version:     version,
			TeamID:      teamID,
			Enabled:     enabled,
			Active:      active,
			Vendor:      vendor,
			Description: desc,
			Writable:    false,
			IsApple:     false,
		})
	}

	return items
}

// scanPluginKitExtensions queries `pluginkit -m -A -v` and filters 3rd party plugins
func scanPluginKitExtensions(ctx context.Context) []ExtensionItem {
	items := make([]ExtensionItem, 0)
	cmd := exec.CommandContext(ctx, "pluginkit", "-m", "-A", "-v")
	out, err := cmd.Output()
	if err != nil {
		return items
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "(") {
			continue
		}

		// Example line:
		// +    com.trendmicro.DrUnzip.FastShareExtension(4.0.9)\tB0EC... \t 2026-07-06... \t /Applications/UnzipOne.app/...
		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}

		header := parts[0]
		path := strings.TrimSpace(parts[len(parts)-1])

		// Determine enabled status: line starts with '+'
		enabled := strings.HasPrefix(strings.TrimSpace(header), "+")

		// Extract bundle ID and version
		cleanHeader := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(header), "+-"))
		bundleID := cleanHeader
		version := ""
		if idx := strings.Index(cleanHeader, "("); idx != -1 {
			bundleID = strings.TrimSpace(cleanHeader[:idx])
			version = strings.TrimSuffix(strings.TrimSpace(cleanHeader[idx+1:]), ")")
		}

		// Filter out internal apple plugins to highlight 3rd-party user extensions
		isApple := strings.HasPrefix(bundleID, "com.apple.")
		if isApple {
			continue
		}

		name := filepath.Base(path)
		name = strings.TrimSuffix(name, ".appex")
		if name == "" {
			name = bundleID
		}

		vendor := inferVendorFromPathOrName(path, name)
		typeName := "Uygulama Eklentisi (Appex)"
		if strings.Contains(strings.ToLower(bundleID), "findersync") || strings.Contains(strings.ToLower(name), "finder") {
			typeName = "Finder Senkronizasyon Eklentisi"
		} else if strings.Contains(strings.ToLower(bundleID), "share") || strings.Contains(strings.ToLower(name), "share") {
			typeName = "Paylaşım Menüsü Eklentisi"
		} else if strings.Contains(strings.ToLower(bundleID), "widget") || strings.Contains(strings.ToLower(name), "widget") {
			typeName = "Masaüstü / Bildirim Araç Takımı (Widget)"
		} else if strings.Contains(strings.ToLower(bundleID), "quicklook") || strings.Contains(strings.ToLower(name), "ql") {
			typeName = "Hızlı Bakış Önizleme Eklentisi"
		}

		items = append(items, ExtensionItem{
			ID:          "pluginkit:" + bundleID,
			Name:        name,
			BundleID:    bundleID,
			Type:        ExtTypePluginKit,
			TypeName:    typeName,
			Path:        path,
			Version:     version,
			Enabled:     enabled,
			Active:      enabled,
			Vendor:      vendor,
			Description: fmt.Sprintf("%s tarafından sağlanan %s", vendor, typeName),
			Writable:    true,
			IsApple:     false,
		})
	}

	return items
}

// scanBundleDirectory scans a folder for .qlgenerator, .mdimporter, .prefPane
func scanBundleDirectory(dir string, extType ExtensionType, typeName, expectedSuffix string) []ExtensionItem {
	items := make([]ExtensionItem, 0)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		isDisabled := strings.HasSuffix(name, ".disabled")
		cleanName := strings.TrimSuffix(name, ".disabled")
		if !strings.HasSuffix(cleanName, expectedSuffix) {
			continue
		}

		fullPath := filepath.Join(dir, name)
		writable := syscall.Access(fullPath, 2) == nil

		baseName := strings.TrimSuffix(cleanName, expectedSuffix)
		vendor := inferVendorFromPathOrName(fullPath, baseName)

		items = append(items, ExtensionItem{
			ID:          string(extType) + ":" + fullPath,
			Name:        baseName,
			BundleID:    baseName,
			Type:        extType,
			TypeName:    typeName,
			Path:        fullPath,
			Enabled:     !isDisabled,
			Active:      !isDisabled,
			Vendor:      vendor,
			Description: fmt.Sprintf("%s dizininde bulunan %s", dir, typeName),
			Writable:    writable,
			IsApple:     strings.Contains(dir, "/System/"),
		})
	}

	return items
}

// ToggleExtension enables or disables an extension
func ToggleExtension(id string, enable bool) error {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("geçersiz eklenti kimliği: %s", id)
	}

	extType := parts[0]
	target := parts[1]

	switch extType {
	case "pluginkit":
		// target is bundleID
		action := "ignore"
		if enable {
			action = "use"
		}
		cmd := exec.Command("pluginkit", "-e", action, "-i", target)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("pluginkit eklentisi değiştirilemedi: %s (%w)", string(out), err)
		}
		return nil

	case "quicklook", "spotlight", "prefpane":
		// target is fullPath
		path := target
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("eklenti dosyası bulunamadı: %s", path)
		}

		if enable {
			if strings.HasSuffix(path, ".disabled") {
				newPath := strings.TrimSuffix(path, ".disabled")
				if err := os.Rename(path, newPath); err != nil {
					// Fallback to admin privileges if root-owned
					adminScript := fmt.Sprintf(`do shell script "mv -f " & quoted form of "%s" & " " & quoted form of "%s" with administrator privileges`, path, newPath)
					if errAdm := exec.Command("osascript", "-e", adminScript).Run(); errAdm != nil {
						return fmt.Errorf("eklenti etkinleştirilemedi: %w", errAdm)
					}
				}
			}
		} else {
			if !strings.HasSuffix(path, ".disabled") {
				newPath := path + ".disabled"
				if err := os.Rename(path, newPath); err != nil {
					// Fallback to admin privileges if root-owned
					adminScript := fmt.Sprintf(`do shell script "mv -f " & quoted form of "%s" & " " & quoted form of "%s" with administrator privileges`, path, newPath)
					if errAdm := exec.Command("osascript", "-e", adminScript).Run(); errAdm != nil {
						return fmt.Errorf("eklenti devre dışı bırakılamadı: %w", errAdm)
					}
				}
			}
		}
		return nil

	case "sysext":
		// System extensions are managed by macOS Settings
		return fmt.Errorf("sistem sürücü eklentileri Apple güvenlik ilkeleri gereği macOS Sistem Ayarları > Gizlilik ve Güvenlik > Eklentiler bölümünden yönetilmelidir")

	default:
		return fmt.Errorf("bilinmeyen eklenti türü: %s", extType)
	}
}

// DeleteExtension removes an extension or moves it to trash
func DeleteExtension(id string, path string) error {
	if path == "" {
		parts := strings.SplitN(id, ":", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[1], "/") {
			path = parts[1]
		}
	}

	if path == "" {
		return fmt.Errorf("silinecek eklenti dosya yolu bulunamadı")
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("dosya bulunamadı: %s", path)
	}

	// Try Trash via AppleScript
	trashScript := fmt.Sprintf(`tell application "Finder" to delete POSIX file "%s"`, path)
	if err := exec.Command("osascript", "-e", trashScript).Run(); err == nil {
		return nil
	}

	// Try direct remove
	if err := os.RemoveAll(path); err == nil {
		return nil
	}

	// Fallback to admin privileges
	adminScript := fmt.Sprintf(`do shell script "rm -rf " & quoted form of "%s" with administrator privileges`, path)
	cmd := exec.Command("osascript", "-e", adminScript)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("eklenti silinemedi: %s (%w)", string(out), err)
	}

	return nil
}

func inferVendorFromBundleID(bundleID, name string) string {
	b := strings.ToLower(bundleID)
	n := strings.ToLower(name)

	switch {
	case strings.Contains(b, "silabs"):
		return "Silicon Laboratories"
	case strings.Contains(b, "protonvpn") || strings.Contains(n, "proton"):
		return "Proton AG"
	case strings.Contains(b, "wireguard"):
		return "WireGuard LLC"
	case strings.Contains(b, "trendmicro") || strings.Contains(n, "unzip"):
		return "Trend Micro"
	case strings.Contains(b, "adobe"):
		return "Adobe"
	case strings.Contains(b, "microsoft"):
		return "Microsoft"
	case strings.Contains(b, "blender"):
		return "Blender Foundation"
	case strings.Contains(b, "google"):
		return "Google"
	case strings.Contains(b, "cleanmymac") || strings.Contains(b, "macpaw"):
		return "MacPaw"
	case strings.Contains(b, "apple"):
		return "Apple"
	default:
		parts := strings.Split(bundleID, ".")
		if len(parts) >= 2 {
			return strings.Title(parts[1])
		}
		return "Geliştirici"
	}
}

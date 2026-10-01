package cleaner

import (
	"os"
	"path/filepath"
	"strings"
)

// RiskLevel indicates how safe an item is to delete
type RiskLevel string

const (
	RiskSafe        RiskLevel = "safe"        // Completely safe, automatically regenerated
	RiskRecommended RiskLevel = "recommended" // Very safe, minor rebuild time penalty
	RiskCaution     RiskLevel = "caution"     // May require re-login or full re-download
)

// CleanTarget defines a known cleanable location or rule
type CleanTarget struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`    // "developer", "system", "app-cache"
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`        // Template path with ~
	Risk        RiskLevel `json:"risk"`
	Size        uint64    `json:"size"`
	SizeStr     string    `json:"sizeStr"`
	ItemCount   int64     `json:"itemCount"`
	Exists      bool      `json:"exists"`
	Selected    bool      `json:"selected"`    // Default selected in UI
	Resolved    string    `json:"resolved"`    // Expanded absolute path
	CleanCmd    string    `json:"cleanCmd"`    // Optional command instead of rm
}

// ExpandPath expands ~ to the user's home directory
func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~/") || p == "~" {
		home, err := os.UserHomeDir()
		if err == nil {
			if p == "~" {
				return home
			}
			return filepath.Join(home, p[2:])
		}
	}
	return p
}

// GetPredefinedTargets returns the curated list of macOS & developer clean targets
func GetPredefinedTargets() []CleanTarget {
	return []CleanTarget{
		// Developer: Xcode
		{
			ID:          "xcode-previews",
			Category:    "developer",
			Name:        "Xcode SwiftUI Previews",
			Description: "SwiftUI Canvas önizleme önbellekleri. Xcode bunları gerektiğinde yeniden üretir.",
			Path:        "~/Library/Developer/Xcode/UserData/Previews",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "xcode-deriveddata",
			Category:    "developer",
			Name:        "Xcode DerivedData",
			Description: "Xcode derleme çıktıları ve indeksler. Projelerinizi bir sonraki derlemede yeniden derler.",
			Path:        "~/Library/Developer/Xcode/DerivedData",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "xcode-archives",
			Category:    "developer",
			Name:        "Xcode Archives",
			Description: "Eski derlenmiş ve paketlenmiş uygulama arşivleri.",
			Path:        "~/Library/Developer/Xcode/Archives",
			Risk:        RiskCaution,
			Selected:    false,
		},
		{
			ID:          "xcode-devicesupport",
			Category:    "developer",
			Name:        "iOS Device Support",
			Description: "Eski iOS sürümlerine ait hata ayıklama sembolleri.",
			Path:        "~/Library/Developer/Xcode/iOS DeviceSupport",
			Risk:        RiskRecommended,
			Selected:    true,
		},
		{
			ID:          "xcode-doc-cache",
			Category:    "developer",
			Name:        "Xcode Documentation Cache",
			Description: "İndirilen ve önbelleğe alınan Xcode dokümantasyon verileri.",
			Path:        "~/Library/Developer/Xcode/DocumentationCache",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "ios-simulators-cache",
			Category:    "developer",
			Name:        "iOS Simülatör Önbellekleri",
			Description: "Kullanılmayan simülatör geçici dosyaları.",
			Path:        "~/Library/Developer/CoreSimulator/Caches",
			Risk:        RiskSafe,
			Selected:    true,
		},

		// Developer: Package Managers & Languages
		{
			ID:          "gradle-caches",
			Category:    "developer",
			Name:        "Gradle Önbellekleri",
			Description: "İndirilen Gradle bağımlılıkları ve derleme önbellekleri (~/.gradle/caches).",
			Path:        "~/.gradle/caches",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "homebrew-cache",
			Category:    "developer",
			Name:        "Homebrew İndirme Önbelleği",
			Description: "İndirilmiş paket arşivleri ve şişeler (bottles). 'brew cleanup' eşdeğeridir.",
			Path:        "~/Library/Caches/Homebrew",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "go-build-cache",
			Category:    "developer",
			Name:        "Go Derleme Önbelleği (go-build)",
			Description: "Go derleyicisinin derleme önbelleği. 'go clean -cache' eşdeğeridir.",
			Path:        "~/Library/Caches/go-build",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "go-mod-cache",
			Category:    "developer",
			Name:        "Go Modül İndirme Önbelleği",
			Description: "İndirilen Go modüllerinin önbelleği (~/go/pkg/mod/cache).",
			Path:        "~/go/pkg/mod/cache",
			Risk:        RiskRecommended,
			Selected:    false,
		},
		{
			ID:          "npm-cache",
			Category:    "developer",
			Name:        "npm Paket Önbelleği",
			Description: "npm paket indirme önbelleği (~/.npm/_cacache).",
			Path:        "~/.npm/_cacache",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "yarn-cache",
			Category:    "developer",
			Name:        "Yarn Paket Önbelleği",
			Description: "Yarn paket indirme önbelleği.",
			Path:        "~/Library/Caches/Yarn",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "pnpm-cache",
			Category:    "developer",
			Name:        "pnpm Önbelleği",
			Description: "pnpm geçici paket önbellekleri.",
			Path:        "~/Library/Caches/pnpm",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "cocoapods-cache",
			Category:    "developer",
			Name:        "CocoaPods Önbelleği",
			Description: "iOS bağımlılıkları indirme önbelleği.",
			Path:        "~/Library/Caches/CocoaPods",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "pip-cache",
			Category:    "developer",
			Name:        "Python Pip Önbelleği",
			Description: "Python pip paket indirme önbelleği.",
			Path:        "~/Library/Caches/pip",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "cargo-cache",
			Category:    "developer",
			Name:        "Rust Cargo Önbelleği",
			Description: "Rust kütüphane indirme önbelleği (~/.cargo/registry/cache).",
			Path:        "~/.cargo/registry/cache",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "arduino-staging",
			Category:    "developer",
			Name:        "Arduino Paket İndirme Arşivi (Staging)",
			Description: "İndirilmiş Arduino kart ve çekirdek paketleri. Kurulum sonrasında güvenle temizlenebilir (~2.9 GB).",
			Path:        "~/Library/Arduino15/staging",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "android-system-images",
			Category:    "developer",
			Name:        "Android Emülatör Sistem İmajları",
			Description: "Android Studio sanal cihazları için indirilmiş Google APIs / Play Store imajları (~4.1 GB).",
			Path:        "~/Library/Android/sdk/system-images",
			Risk:        RiskRecommended,
			Selected:    false,
		},
		{
			ID:          "android-avd",
			Category:    "developer",
			Name:        "Android Sanal Cihazları (AVD Kalıpları)",
			Description: "Kullanılmayan Android Studio emülatör sanal diskleri.",
			Path:        "~/.android/avd",
			Risk:        RiskCaution,
			Selected:    false,
		},
		{
			ID:          "vscode-cached-vsix",
			Category:    "developer",
			Name:        "VS Code Eklenti Kurulum Arşivi (VSIX)",
			Description: "VS Code tarafından indirilen eklenti kurulum paketleri (CachedExtensionVSIXs).",
			Path:        "~/Library/Application Support/Code/CachedExtensionVSIXs",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "vscode-app-cache",
			Category:    "developer",
			Name:        "VS Code Uygulama Önbelleği",
			Description: "VS Code webview, GPU ve geçici çalışma önbellekleri.",
			Path:        "~/Library/Application Support/Code/Cache",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "chrome-serviceworker-cache",
			Category:    "system",
			Name:        "Google Chrome Service Worker Önbelleği",
			Description: "Web sitelerinin arka planda depoladığı çevrimdışı ve geçici veriler.",
			Path:        "~/Library/Application Support/Google/Chrome/Default/Service Worker/CacheStorage",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "discord-cache",
			Category:    "system",
			Name:        "Discord Medya ve Görsel Önbelleği",
			Description: "Discord kanallarından indirilen profil resimleri, gifler ve medya önbellekleri.",
			Path:        "~/Library/Application Support/discord/Cache",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "steam-appcache",
			Category:    "system",
			Name:        "Steam İstemci & İndirme Önbelleği",
			Description: "Steam geçici indirme ve web arayüzü önbellekleri.",
			Path:        "~/Library/Application Support/Steam/appcache",
			Risk:        RiskSafe,
			Selected:    true,
		},

		// System Junk & App Updates
		{
			ID:          "vscode-shipit",
			Category:    "system",
			Name:        "VS Code Güncelleme Artıkları (ShipIt)",
			Description: "VS Code otomatik güncelleme indiricisine ait eski kurulum paketleri.",
			Path:        "~/Library/Caches/com.microsoft.VSCode.ShipIt",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "google-cache",
			Category:    "system",
			Name:        "Google / Chrome Önbellek Artıkları",
			Description: "Chrome ve Google servislerine ait geçici indirme ve önbellek dosyaları.",
			Path:        "~/Library/Caches/Google",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "adobe-cache",
			Category:    "system",
			Name:        "Adobe Geçici Önbellekleri",
			Description: "Adobe Creative Cloud ve uygulamalarına ait geçici dosyalar.",
			Path:        "~/Library/Caches/Adobe",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "antigravity-shipit",
			Category:    "system",
			Name:        "Antigravity / IDE Güncelleme Artıkları",
			Description: "Eski IDE güncelleme indiricisi paketleri ve geçici dosyalar.",
			Path:        "~/Library/Caches/com.google.antigravity.ShipIt",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "antigravity-updater",
			Category:    "system",
			Name:        "Antigravity Güncelleyici Önbelleği",
			Description: "Arka plan güncelleme önbelleği.",
			Path:        "~/Library/Caches/antigravity-updater",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "crossover-cache",
			Category:    "system",
			Name:        "CrossOver Geçici Dosyaları",
			Description: "CrossOver şişeleme ve indirme önbelleği.",
			Path:        "~/Library/Caches/com.codeweavers.CrossOver",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "phpwebstudy-updater",
			Category:    "system",
			Name:        "PhpWebStudy Güncelleyici Önbelleği",
			Description: "Eski güncelleme indiricisi dosyaları.",
			Path:        "~/Library/Caches/phpwebstudy-updater",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "freecad-cache",
			Category:    "system",
			Name:        "FreeCAD Önbelleği",
			Description: "FreeCAD geçici hesaplama ve render önbelleği.",
			Path:        "~/Library/Caches/FreeCAD",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "epicgames-cache",
			Category:    "system",
			Name:        "Epic Games Launcher Önbelleği",
			Description: "Epic Games Launcher web görünümü ve geçici indirme önbelleği.",
			Path:        "~/Library/Caches/com.epicgames.EpicGamesLauncher",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "steam-cache",
			Category:    "system",
			Name:        "Steam İstemci Önbelleği",
			Description: "Steam web tarayıcısı ve önbellek dosyaları.",
			Path:        "~/Library/Caches/Steam",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "system-logs",
			Category:    "system",
			Name:        "Kullanıcı ve Uygulama Günlükleri (Logs)",
			Description: "Zamanla biriken hata, çökme raporları ve teşhis logları.",
			Path:        "~/Library/Logs",
			Risk:        RiskSafe,
			Selected:    true,
		},
		{
			ID:          "trash-bin",
			Category:    "system",
			Name:        "macOS Çöp Kutusu",
			Description: "Çöp kutusunda bekleyen ve diskte yer kaplayan dosyalar.",
			Path:        "~/.Trash",
			Risk:        RiskRecommended,
			Selected:    true,
		},
		{
			ID:          "quicklook-cache",
			Category:    "system",
			Name:        "QuickLook Küçük Resim Önbelleği",
			Description: "macOS Finder önizleme ve küçük resim önbelleği.",
			Path:        "~/Library/Caches/com.apple.QuickLook.thumbnailcache",
			Risk:        RiskSafe,
			Selected:    true,
		},
	}
}

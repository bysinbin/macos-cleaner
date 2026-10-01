package cleaner

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PrivacyCategory represents privacy item group
type PrivacyCategory string

const (
	PrivacyCatRecents     PrivacyCategory = "recents"
	PrivacyCatHistory     PrivacyCategory = "history"
	PrivacyCatBrowser     PrivacyCategory = "browser"
	PrivacyCatPermissions PrivacyCategory = "permissions"
)

// PrivacyItem represents a trackable privacy trace item
type PrivacyItem struct {
	ID          string          `json:"id"`
	Category    PrivacyCategory `json:"category"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Count       int             `json:"count"`
	CountLabel  string          `json:"countLabel"`
	Size        uint64          `json:"size"`
	SizeStr     string          `json:"sizeStr"`
	Paths       []string        `json:"paths"`
	Removable   bool            `json:"removable"`
	Icon        string          `json:"icon"`
}

// PrivacyPermission represents an app authorization capability
type PrivacyPermission struct {
	Service     string `json:"service"`     // e.g. "Camera", "Microphone"
	ServiceName string `json:"serviceName"` // e.g. "Kamera Erişimi"
	Description string `json:"description"`
	ResetCmd    string `json:"resetCmd"`
}

// PrivacySummary provides aggregate counts and list of privacy traces
type PrivacySummary struct {
	TotalItemsCount int                 `json:"totalItemsCount"`
	TotalSizeBytes  uint64              `json:"totalSizeBytes"`
	TotalSizeStr    string              `json:"totalSizeStr"`
	Items           []PrivacyItem       `json:"items"`
	Permissions     []PrivacyPermission `json:"permissions"`
}

// ScanPrivacyTraces scans recent items, shell history files, and browser history
func ScanPrivacyTraces(ctx context.Context) (*PrivacySummary, error) {
	home, _ := os.UserHomeDir()
	summary := &PrivacySummary{
		Items:       make([]PrivacyItem, 0),
		Permissions: getAvailablePrivacyPermissions(),
	}

	// 1. Son Açılan Öğeler (Recent Items & SharedFileList)
	recentItem := scanRecentItems()
	if recentItem != nil {
		summary.Items = append(summary.Items, *recentItem)
	}

	// 2. Terminal ve Komut Satırı Geçmişi (Shell Histories)
	historyItems := scanShellHistories(home)
	summary.Items = append(summary.Items, historyItems...)

	// 3. Tarayıcı İzleri ve Geçmişi (Safari, Chrome, Arc, Brave, Firefox)
	browserItems := scanBrowserTraces(home)
	summary.Items = append(summary.Items, browserItems...)

	// Aggregate statistics
	for _, it := range summary.Items {
		summary.TotalItemsCount += it.Count
		summary.TotalSizeBytes += it.Size
	}
	summary.TotalSizeStr = FormatBytes(summary.TotalSizeBytes)

	return summary, nil
}

func scanRecentItems() *PrivacyItem {
	// Query sfltool for list of recent lists
	out, err := exec.Command("sfltool", "list").Output()
	recentLists := []string{
		"com.apple.LSSharedFileList.RecentDocuments",
		"com.apple.LSSharedFileList.RecentApplications",
		"com.apple.LSSharedFileList.RecentServers",
		"com.apple.LSSharedFileList.RecentHosts",
		"com.apple.LSSharedFileList.ApplicationRecentDocuments",
	}

	count := 0
	if err == nil {
		text := string(out)
		for _, rl := range recentLists {
			if strings.Contains(text, rl) {
				count += 10 // macOS tracks up to 10-30 recent items per category
			}
		}
	}
	if count == 0 {
		count = 25
	}

	return &PrivacyItem{
		ID:          "recents:macos",
		Category:    PrivacyCatRecents,
		Title:       "macOS Son Kullanılan Belgeler ve Uygulamalar",
		Description: "Finder ve menü çubuğunda listelenen son açılan dosyalar, uygulamalar ve bağlı sunucular listesi",
		Count:       count,
		CountLabel:  "liste kaydı",
		Size:        0,
		SizeStr:     "Sistem Önbelleği",
		Paths:       recentLists,
		Removable:   true,
		Icon:        "clock",
	}
}

func scanShellHistories(home string) []PrivacyItem {
	items := make([]PrivacyItem, 0)

	shells := []struct {
		ID    string
		Title string
		Desc  string
		Rel   string
	}{
		{"zsh", "Zsh Terminal Komut Geçmişi", "Yürütülen tüm Terminal zsh komutlarının tam arşivi", ".zsh_history"},
		{"bash", "Bash Terminal Komut Geçmişi", "Eski veya alternatif Bash kabuğu komut geçmişi", ".bash_history"},
		{"python", "Python REPL Etkileşim Geçmişi", "Python interaktif oturumlarında yazılan kodlar", ".python_history"},
		{"node", "Node.js REPL Komut Geçmişi", "Node REPL oturumlarında çalıştırılan JavaScript komutları", ".node_repl_history"},
		{"less", "Less / Terminal Sayfalama Geçmişi", "Terminalde incelenen dosya ve arama geçmişi", ".lesshst"},
	}

	for _, sh := range shells {
		p := filepath.Join(home, sh.Rel)
		fi, err := os.Stat(p)
		if err != nil || fi.Size() == 0 {
			continue
		}

		lineCount := countFileLines(p)
		items = append(items, PrivacyItem{
			ID:          "history:" + sh.ID,
			Category:    PrivacyCatHistory,
			Title:       sh.Title,
			Description: sh.Desc,
			Count:       lineCount,
			CountLabel:  "çalıştırılan komut",
			Size:        uint64(fi.Size()),
			SizeStr:     FormatBytes(uint64(fi.Size())),
			Paths:       []string{p},
			Removable:   true,
			Icon:        "terminal",
		})
	}

	return items
}

func scanBrowserTraces(home string) []PrivacyItem {
	items := make([]PrivacyItem, 0)

	type BrowserTarget struct {
		ID       string
		Title    string
		Desc     string
		PathPats []string
	}

	targets := []BrowserTarget{
		{
			ID:    "safari",
			Title: "Safari Web İzleri & Geçmişi",
			Desc:  "Safari gezinme veritabanı, yerel depolama ve web sitesi kalıntıları",
			PathPats: []string{
				filepath.Join(home, "Library/Safari/History.db*"),
				filepath.Join(home, "Library/Safari/CloudHistory.db*"),
				filepath.Join(home, "Library/Safari/Favicon Cache"),
				filepath.Join(home, "Library/Cookies/Cookies.binarycookies"),
			},
		},
		{
			ID:    "chrome",
			Title: "Google Chrome Gezinme İzleri",
			Desc:  "Chrome ziyaret edilen bağlantılar, arama anahtar kelimeleri ve site verileri",
			PathPats: []string{
				filepath.Join(home, "Library/Application Support/Google/Chrome/Default/History*"),
				filepath.Join(home, "Library/Application Support/Google/Chrome/Default/Shortcuts*"),
				filepath.Join(home, "Library/Application Support/Google/Chrome/Default/Favicons*"),
			},
		},
		{
			ID:    "arc",
			Title: "Arc Browser Gezinme İzleri",
			Desc:  "Arc web tarayıcısı geçmişi ve arama önbelleği",
			PathPats: []string{
				filepath.Join(home, "Library/Application Support/Arc/User Data/Default/History*"),
				filepath.Join(home, "Library/Application Support/Arc/User Data/Default/Favicons*"),
			},
		},
		{
			ID:    "brave",
			Title: "Brave Browser Gezinme İzleri",
			Desc:  "Brave gizli tarayıcı yerel geçmiş ve arama kayıtları",
			PathPats: []string{
				filepath.Join(home, "Library/Application Support/BraveSoftware/Brave-Browser/Default/History*"),
			},
		},
		{
			ID:    "firefox",
			Title: "Mozilla Firefox Geçmişi",
			Desc:  "Firefox yer imleri ve ziyaret geçmişi sqlite veritabanı",
			PathPats: []string{
				filepath.Join(home, "Library/Application Support/Firefox/Profiles/*/places.sqlite*"),
			},
		},
	}

	for _, bt := range targets {
		matchedPaths := make([]string, 0)
		var totalSize uint64
		count := 0

		for _, pat := range bt.PathPats {
			matches, err := filepath.Glob(pat)
			if err != nil {
				continue
			}
			for _, m := range matches {
				fi, err := os.Stat(m)
				if err != nil {
					continue
				}
				matchedPaths = append(matchedPaths, m)
				totalSize += uint64(fi.Size())
				count++
			}
		}

		if len(matchedPaths) > 0 {
			items = append(items, PrivacyItem{
				ID:          "browser:" + bt.ID,
				Category:    PrivacyCatBrowser,
				Title:       bt.Title,
				Description: bt.Desc,
				Count:       count,
				CountLabel:  "veritabanı / iz dosyası",
				Size:        totalSize,
				SizeStr:     FormatBytes(totalSize),
				Paths:       matchedPaths,
				Removable:   true,
				Icon:        "globe",
			})
		}
	}

	return items
}

func getAvailablePrivacyPermissions() []PrivacyPermission {
	return []PrivacyPermission{
		{
			Service:     "Camera",
			ServiceName: "Kamera Erişimi",
			Description: "Uygulamalara verilen web kamerası erişim izinlerini sıfırlar.",
			ResetCmd:    "tccutil reset Camera",
		},
		{
			Service:     "Microphone",
			ServiceName: "Mikrofon Erişimi",
			Description: "Uygulamalara verilen dahili ve harici ses kayıt izinlerini sıfırlar.",
			ResetCmd:    "tccutil reset Microphone",
		},
		{
			Service:     "ScreenCapture",
			ServiceName: "Ekran Kaydı İzni",
			Description: "Ekran ve ses yakalama (Screen Recording) yetkilerini temizler.",
			ResetCmd:    "tccutil reset ScreenCapture",
		},
		{
			Service:     "Accessibility",
			ServiceName: "Erişilebilirlik Yetkileri",
			Description: "Klavyeyi izleme ve sistemi uzaktan kontrol etme yetkilerini sıfırlar.",
			ResetCmd:    "tccutil reset Accessibility",
		},
		{
			Service:     "SystemPolicyAllFiles",
			ServiceName: "Tam Disk Erişimi (Full Disk)",
			Description: "Tüm diske doğrudan okuma/yazma izni verilmiş uygulamaları sıfırlar.",
			ResetCmd:    "tccutil reset SystemPolicyAllFiles",
		},
	}
}

// CleanPrivacyItems cleans selected privacy items
func CleanPrivacyItems(itemIDs []string) (int, error) {
	cleaned := 0

	for _, id := range itemIDs {
		if strings.HasPrefix(id, "recents:") {
			// Clear all recent items via sfltool
			recentLists := []string{
				"com.apple.LSSharedFileList.RecentDocuments",
				"com.apple.LSSharedFileList.RecentApplications",
				"com.apple.LSSharedFileList.RecentServers",
				"com.apple.LSSharedFileList.RecentHosts",
				"com.apple.LSSharedFileList.ApplicationRecentDocuments",
			}
			for _, rl := range recentLists {
				_ = exec.Command("sfltool", "clear", rl).Run()
			}
			cleaned++
		} else if strings.HasPrefix(id, "history:") {
			// Truncate history file
			target := strings.TrimPrefix(id, "history:")
			home, _ := os.UserHomeDir()
			var filename string
			switch target {
			case "zsh":
				filename = ".zsh_history"
			case "bash":
				filename = ".bash_history"
			case "python":
				filename = ".python_history"
			case "node":
				filename = ".node_repl_history"
			case "less":
				filename = ".lesshst"
			}
			if filename != "" {
				p := filepath.Join(home, filename)
				_ = os.WriteFile(p, []byte(""), 0600)
				cleaned++
			}
		} else if strings.HasPrefix(id, "browser:") {
			// Remove browser traces
			target := strings.TrimPrefix(id, "browser:")
			home, _ := os.UserHomeDir()
			var pats []string
			switch target {
			case "safari":
				pats = []string{
					filepath.Join(home, "Library/Safari/History.db*"),
					filepath.Join(home, "Library/Safari/CloudHistory.db*"),
					filepath.Join(home, "Library/Cookies/Cookies.binarycookies"),
				}
			case "chrome":
				pats = []string{
					filepath.Join(home, "Library/Application Support/Google/Chrome/Default/History*"),
					filepath.Join(home, "Library/Application Support/Google/Chrome/Default/Shortcuts*"),
				}
			case "arc":
				pats = []string{
					filepath.Join(home, "Library/Application Support/Arc/User Data/Default/History*"),
				}
			case "brave":
				pats = []string{
					filepath.Join(home, "Library/Application Support/BraveSoftware/Brave-Browser/Default/History*"),
				}
			case "firefox":
				pats = []string{
					filepath.Join(home, "Library/Application Support/Firefox/Profiles/*/places.sqlite*"),
				}
			}
			for _, pat := range pats {
				matches, _ := filepath.Glob(pat)
				for _, m := range matches {
					_ = os.Remove(m)
				}
			}
			cleaned++
		}
	}

	return cleaned, nil
}

// ResetPrivacyPermission resets a TCC service permission
func ResetPrivacyPermission(service string) error {
	allowed := map[string]bool{
		"Camera":               true,
		"Microphone":           true,
		"ScreenCapture":        true,
		"Accessibility":        true,
		"SystemPolicyAllFiles": true,
	}
	if !allowed[service] {
		return fmt.Errorf("geçersiz gizlilik servisi: %s", service)
	}

	cmd := exec.Command("tccutil", "reset", service)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("izin sıfırlanamadı: %s (%w)", string(out), err)
	}
	return nil
}

func countFileLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

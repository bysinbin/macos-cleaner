package cleaner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// SmartCareItem represents an individual safe-to-clean category
type SmartCareItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Size        uint64 `json:"size"`
	SizeStr     string `json:"sizeStr"`
	FileCount   int64  `json:"fileCount"`
	Selected    bool   `json:"selected"`
}

// SmartCareResult encapsulates the one-click system health check
type SmartCareResult struct {
	Items             []SmartCareItem `json:"items"`
	TotalCleanable    uint64          `json:"totalCleanable"`
	TotalCleanableStr string          `json:"totalCleanableStr"`
	TotalFiles        int64           `json:"totalFiles"`
	DiskUsedPct       float64         `json:"diskUsedPct"`
	DiskFreeStr       string          `json:"diskFreeStr"`
	DiskTotalStr      string          `json:"diskTotalStr"`
	IsCritical        bool            `json:"isCritical"` // Used > 85%
	HealthStatus      string          `json:"healthStatus"` // "Harika", "İyi", "Dikkat", "Kritik"
	Recommendations   []string        `json:"recommendations"`
}

// ScanSmartCare performs a rapid, safe health analysis across the system
func ScanSmartCare(ctx context.Context) (*SmartCareResult, error) {
	scanner := NewScanner(ctx)
	home, _ := os.UserHomeDir()

	stats, _ := GetDiskStats("/")
	usedPct := 0.0
	freeStr := "--"
	totalStr := "--"
	if stats != nil {
		usedPct = stats.UsedPercent
		freeStr = stats.FreeStr
		totalStr = stats.TotalStr
	}

	result := &SmartCareResult{
		DiskUsedPct:  usedPct,
		DiskFreeStr:  freeStr,
		DiskTotalStr: totalStr,
		IsCritical:   usedPct >= 85.0,
	}

	if usedPct >= 90.0 {
		result.HealthStatus = "Kritik"
		result.Recommendations = append(result.Recommendations, "Disk doluluğu %90'ın üzerinde. Büyük dosyaları ve uygulama artıklarını temizlemeniz önerilir.")
	} else if usedPct >= 80.0 {
		result.HealthStatus = "Dikkat"
		result.Recommendations = append(result.Recommendations, "Disk doluluk oranı yükseliyor. Akıllı temizlik ile yer açabilirsiniz.")
	} else {
		result.HealthStatus = "İyi"
		result.Recommendations = append(result.Recommendations, "Sistem depolama durumu sağlıklı. Rutin bakım gerçekleştirebilirsiniz.")
	}

	var items []SmartCareItem
	var mu sync.Mutex
	var wg sync.WaitGroup

	var totalCleanable uint64
	var totalFiles int64

	EmitProgress("smartcare", "Akıllı Bakım: Sistem, önbellek ve çöp taraması başlatılıyor...", 1, 6)

	// 1. User Caches
	wg.Add(1)
	go func() {
		defer wg.Done()
		cacheDir := filepath.Join(home, "Library", "Caches")
		sz, cnt, _ := scanner.CalculateDirSize(cacheDir)
		if sz > 0 {
			it := SmartCareItem{
				ID:          "user_caches",
				Title:       "Kullanıcı ve Sistem Önbellekleri",
				Description: "Uygulamaların geçici görsel, web ve operasyonel önbellek dosyaları.",
				Size:        sz,
				SizeStr:     FormatBytes(sz),
				FileCount:   cnt,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, sz)
			atomic.AddInt64(&totalFiles, cnt)
			mu.Unlock()
		}
	}()

	// 2. System Logs
	wg.Add(1)
	go func() {
		defer wg.Done()
		logsDir := filepath.Join(home, "Library", "Logs")
		sz, cnt, _ := scanner.CalculateDirSize(logsDir)
		if sz > 0 {
			it := SmartCareItem{
				ID:          "system_logs",
				Title:       "Uygulama ve Sistem Günlükleri (Logs)",
				Description: "Arka planda kaydedilen eski çalışma ve hata günlükleri.",
				Size:        sz,
				SizeStr:     FormatBytes(sz),
				FileCount:   cnt,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, sz)
			atomic.AddInt64(&totalFiles, cnt)
			mu.Unlock()
		}
	}()

	// 3. Browser Caches
	wg.Add(1)
	go func() {
		defer wg.Done()
		brRes, err := ScanBrowserCaches(ctx)
		if err == nil && brRes.TotalSize > 0 {
			it := SmartCareItem{
				ID:          "browser_caches",
				Title:       "Tarayıcı Web Önbellekleri",
				Description: "Chrome, Safari, Firefox vb. tarayıcıların web sayfalarını hızlandırmak için tuttuğu önbellekler.",
				Size:        brRes.TotalSize,
				SizeStr:     brRes.TotalSizeStr,
				FileCount:   brRes.TotalFiles,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, brRes.TotalSize)
			atomic.AddInt64(&totalFiles, brRes.TotalFiles)
			mu.Unlock()
		}
	}()

	// 4. Trash
	wg.Add(1)
	go func() {
		defer wg.Done()
		trashDir := filepath.Join(home, ".Trash")
		sz, cnt, _ := scanner.CalculateDirSize(trashDir)
		if sz > 0 {
			it := SmartCareItem{
				ID:          "trash",
				Title:       "macOS Çöp Sepeti (Trash)",
				Description: "Çöp kutusunda bekleyen ve henüz kalıcı olarak silinmemiş dosyalar.",
				Size:        sz,
				SizeStr:     FormatBytes(sz),
				FileCount:   cnt,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, sz)
			atomic.AddInt64(&totalFiles, cnt)
			mu.Unlock()
		}
	}()

	// 5. Leftover Debris (.DS_Store & Broken symlinks)
	wg.Add(1)
	go func() {
		defer wg.Done()
		leftovers, err := ScanLeftovers(ctx)
		if err == nil && leftovers.TotalSize > 0 {
			it := SmartCareItem{
				ID:          "leftovers_debris",
				Title:       "Finder Artıkları ve Eski Durumlar",
				Description: ".DS_Store pencere ayarları, kırık sembolik linkler ve oturum kalıntıları.",
				Size:        leftovers.TotalSize,
				SizeStr:     leftovers.TotalSizeStr,
				FileCount:   leftovers.TotalFiles,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, leftovers.TotalSize)
			atomic.AddInt64(&totalFiles, leftovers.TotalFiles)
			mu.Unlock()
		}
	}()

	// 6. Antigravity Ajan & Tarayıcı Kayıtları
	wg.Add(1)
	go func() {
		defer wg.Done()
		agSz, agCnt, err := ScanAntigravityDebris(ctx)
		if err == nil && agSz > 0 {
			it := SmartCareItem{
				ID:          "antigravity_debris",
				Title:       "Antigravity Ajan & Tarayıcı Kayıtları",
				Description: "Ajan oturumlarında kaydedilen ekran görüntüleri, tarayıcı videoları ve geçici dosyalar.",
				Size:        agSz,
				SizeStr:     FormatBytes(agSz),
				FileCount:   agCnt,
				Selected:    true,
			}
			mu.Lock()
			items = append(items, it)
			atomic.AddUint64(&totalCleanable, agSz)
			atomic.AddInt64(&totalFiles, agCnt)
			mu.Unlock()
		}
	}()

	wg.Wait()

	result.Items = items
	result.TotalCleanable = totalCleanable
	result.TotalCleanableStr = FormatBytes(totalCleanable)
	result.TotalFiles = totalFiles

	EmitProgress("smartcare", fmt.Sprintf("Akıllı Bakım taraması tamamlandı: %s alan kazanılabilir.", result.TotalCleanableStr), 6, 6)

	return result, nil
}

// ExecuteSmartCareClean performs one-click cleaning of all selected safe components
func ExecuteSmartCareClean(selectedIDs []string) (uint64, int64, error) {
	home, _ := os.UserHomeDir()

	idMap := make(map[string]bool)
	for _, id := range selectedIDs {
		idMap[id] = true
	}

	var freedBytes uint64
	var deletedFiles int64

	// 1. User Caches
	if len(idMap) == 0 || idMap["user_caches"] {
		cacheDir := filepath.Join(home, "Library", "Caches")
		entries, _ := os.ReadDir(cacheDir)
		for _, e := range entries {
			// Don't wipe the root cache dir itself, delete contents
			p := filepath.Join(cacheDir, e.Name())
			fi, err := os.Stat(p)
			if err != nil {
				continue
			}
			sz := uint64(fi.Size())
			if fi.IsDir() {
				scanner := NewScanner(context.Background())
				s, c, _ := scanner.CalculateDirSize(p)
				sz = s
				deletedFiles += c
			} else {
				deletedFiles++
			}
			if err := os.RemoveAll(p); err == nil {
				freedBytes += sz
			}
		}
	}

	// 2. System Logs
	if len(idMap) == 0 || idMap["system_logs"] {
		logsDir := filepath.Join(home, "Library", "Logs")
		entries, _ := os.ReadDir(logsDir)
		for _, e := range entries {
			p := filepath.Join(logsDir, e.Name())
			fi, err := os.Stat(p)
			if err != nil {
				continue
			}
			sz := uint64(fi.Size())
			if fi.IsDir() {
				scanner := NewScanner(context.Background())
				s, c, _ := scanner.CalculateDirSize(p)
				sz = s
				deletedFiles += c
			} else {
				deletedFiles++
			}
			if err := os.RemoveAll(p); err == nil {
				freedBytes += sz
			}
		}
	}

	// 3. Browser Caches
	if len(idMap) == 0 || idMap["browser_caches"] {
		sz, cnt, _ := CleanBrowserCaches(nil)
		freedBytes += sz
		deletedFiles += cnt
	}

	// 4. Trash
	if len(idMap) == 0 || idMap["trash"] {
		trashDir := filepath.Join(home, ".Trash")
		entries, _ := os.ReadDir(trashDir)
		for _, e := range entries {
			p := filepath.Join(trashDir, e.Name())
			fi, err := os.Stat(p)
			if err != nil {
				continue
			}
			sz := uint64(fi.Size())
			if fi.IsDir() {
				scanner := NewScanner(context.Background())
				s, c, _ := scanner.CalculateDirSize(p)
				sz = s
				deletedFiles += c
			} else {
				deletedFiles++
			}
			if err := os.RemoveAll(p); err == nil {
				freedBytes += sz
			}
		}
	}

	// 5. Leftover Debris (.DS_Store)
	if len(idMap) == 0 || idMap["leftovers_debris"] {
		sz, cnt, _ := CleanDSStoreBatch()
		freedBytes += sz
		deletedFiles += cnt
	}

	// 6. Antigravity Debris (Screenshots, Browser recordings, Temp media)
	if len(idMap) == 0 || idMap["antigravity_debris"] {
		sz, cnt, _ := CleanAntigravityDebris()
		freedBytes += sz
		deletedFiles += cnt
	}

	return freedBytes, deletedFiles, nil
}

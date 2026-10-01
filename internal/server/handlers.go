package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"disk-cleaner/internal/cleaner"
	"disk-cleaner/internal/config"
)

func jsonEncode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// SystemResponse carries system and storage status
type SystemResponse struct {
	Stats   *cleaner.DiskStats `json:"stats"`
	HomeDir string             `json:"homeDir"`
	OS      string             `json:"os"`
}

func (s *Server) handleSystemStats(w http.ResponseWriter, r *http.Request) {
	stats, err := cleaner.GetDiskStats("/")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := SystemResponse{
		Stats:   stats,
		HomeDir: getCurrentUserHome(),
		OS:      "macOS",
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	stream := r.URL.Query().Get("stream") == "true" || strings.Contains(r.Header.Get("Accept"), "text/event-stream")

	scanner := cleaner.NewScanner(r.Context())

	if stream {
		flusher, ok := setupSSE(w)
		if !ok {
			return
		}

		progressChan := make(chan cleaner.ScanProgress, 20)
		doneChan := make(chan struct{})

		var result *cleaner.ScanResult
		var scanErr error

		go func() {
			result, scanErr = scanner.ScanTargets(progressChan)
			close(doneChan)
		}()

		for {
			select {
			case prog, ok := <-progressChan:
				if !ok {
					continue
				}
				data, _ := json.Marshal(prog)
				fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
				flusher.Flush()
			case <-doneChan:
				if scanErr != nil {
					fmt.Fprintf(w, "event: error\ndata: {\"error\": %q}\n\n", scanErr.Error())
				} else {
					data, _ := json.Marshal(result)
					fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
				}
				flusher.Flush()
				return
			case <-r.Context().Done():
				return
			}
		}
	} else {
		result, err := scanner.ScanTargets(nil)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func (s *Server) handleClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req cleaner.CleanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek: " + err.Error()})
		return
	}

	stream := r.URL.Query().Get("stream") == "true" || strings.Contains(r.Header.Get("Accept"), "text/event-stream")
	c := cleaner.NewCleaner(req.DryRun)

	if stream {
		flusher, ok := setupSSE(w)
		if !ok {
			return
		}

		progressChan := make(chan cleaner.CleanProgress, 20)
		go c.ExecuteClean(req, progressChan)

		for prog := range progressChan {
			data, _ := json.Marshal(prog)
			fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
			flusher.Flush()
		}
	} else {
		progressChan := make(chan cleaner.CleanProgress, 20)
		go c.ExecuteClean(req, progressChan)

		var lastProg cleaner.CleanProgress
		for p := range progressChan {
			lastProg = p
		}
		writeJSON(w, http.StatusOK, lastProg)
	}
}

func (s *Server) handleNodeModulesScan(w http.ResponseWriter, r *http.Request) {
	customPath := r.URL.Query().Get("root")
	var roots []string
	if customPath != "" {
		roots = []string{customPath}
	}

	stream := r.URL.Query().Get("stream") == "true" || strings.Contains(r.Header.Get("Accept"), "text/event-stream")

	if stream {
		flusher, ok := setupSSE(w)
		if !ok {
			return
		}

		progressChan := make(chan cleaner.NodeModuleProgress, 20)
		doneChan := make(chan struct{})
		var results []cleaner.NodeModuleItem
		var scanErr error

		go func() {
			results, scanErr = cleaner.ScanNodeModules(r.Context(), roots, progressChan)
			close(doneChan)
		}()

		for {
			select {
			case prog, ok := <-progressChan:
				if !ok {
					continue
				}
				data, _ := json.Marshal(prog)
				fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
				flusher.Flush()
			case <-doneChan:
				if scanErr != nil {
					fmt.Fprintf(w, "event: error\ndata: {\"error\": %q}\n\n", scanErr.Error())
				} else {
					data, _ := json.Marshal(results)
					fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
				}
				flusher.Flush()
				return
			case <-r.Context().Done():
				return
			}
		}
	} else {
		results, err := cleaner.ScanNodeModules(r.Context(), roots, nil)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, results)
	}
}

func (s *Server) handleNodeModulesClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Paths       []string `json:"paths"`
		MoveToTrash bool     `json:"moveToTrash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c := cleaner.NewCleaner(false)
	var totalFreed uint64
	var deletedCount int
	var errors []string

	for _, p := range req.Paths {
		// Verify it's actually a node_modules folder for safety
		if !strings.HasSuffix(p, "node_modules") && !strings.Contains(p, "node_modules") {
			errors = append(errors, fmt.Sprintf("Güvenlik hatası: Sadece node_modules klasörleri silinebilir: %s", p))
			continue
		}

		freed, _, err := c.CleanTarget(p, req.MoveToTrash)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s silinemedi: %v", p, err))
		} else {
			totalFreed += freed
			deletedCount++
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deletedCount": deletedCount,
		"freedBytes":   totalFreed,
		"freedStr":     cleaner.FormatBytes(totalFreed),
		"errors":       errors,
	})
}

func (s *Server) handleLargeFilesScan(w http.ResponseWriter, r *http.Request) {
	minSizeMB, _ := strconv.Atoi(r.URL.Query().Get("minSizeMb"))
	if minSizeMB <= 0 {
		minSizeMB = 50
	}

	minAgeDays, _ := strconv.Atoi(r.URL.Query().Get("minAgeDays"))
	category := r.URL.Query().Get("category")
	customPath := r.URL.Query().Get("path")

	var paths []string
	if customPath != "" {
		paths = []string{customPath}
	}

	filter := cleaner.LargeFileFilter{
		Paths:        paths,
		MinSizeBytes: uint64(minSizeMB) * 1024 * 1024,
		MinAgeDays:   minAgeDays,
		Category:     category,
	}

	results, err := cleaner.ScanLargeFiles(r.Context(), filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, results)
}

func (s *Server) handleLargeFilesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Paths       []string `json:"paths"`
		MoveToTrash bool     `json:"moveToTrash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c := cleaner.NewCleaner(false)
	var totalFreed uint64
	var deletedCount int
	var errors []string

	for _, p := range req.Paths {
		// Stat before deleting to get size
		fi, err := os.Stat(p)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s bulunamadı: %v", p, err))
			continue
		}
		size := uint64(fi.Size())

		_, _, err = c.CleanTarget(p, req.MoveToTrash)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s silinemedi: %v", p, err))
		} else {
			totalFreed += size
			deletedCount++
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deletedCount": deletedCount,
		"freedBytes":   totalFreed,
		"freedStr":     cleaner.FormatBytes(totalFreed),
		"errors":       errors,
	})
}

func (s *Server) handleReveal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := cleaner.RevealInFinder(cleaner.ExpandPath(req.Path))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (s *Server) handleDirTree(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Query().Get("path")
	breakdown, err := cleaner.AnalyzeDirectory(r.Context(), targetPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, breakdown)
}

func (s *Server) handleDuplicates(w http.ResponseWriter, r *http.Request) {
	minSizeMB, _ := strconv.Atoi(r.URL.Query().Get("minSizeMb"))
	if minSizeMB <= 0 {
		minSizeMB = 1
	}

	customPath := r.URL.Query().Get("path")
	var roots []string
	if customPath != "" {
		roots = []string{customPath}
	}

	groups, err := cleaner.FindDuplicates(r.Context(), roots, uint64(minSizeMB)*1024*1024)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, groups)
}

func (s *Server) handleDuplicatesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Paths       []string `json:"paths"`
		MoveToTrash bool     `json:"moveToTrash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c := cleaner.NewCleaner(false)
	var totalFreed uint64
	var deletedCount int
	var errors []string

	for _, p := range req.Paths {
		fi, err := os.Stat(p)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s bulunamadı: %v", p, err))
			continue
		}
		size := uint64(fi.Size())

		_, _, err = c.CleanTarget(p, req.MoveToTrash)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s silinemedi: %v", p, err))
		} else {
			totalFreed += size
			deletedCount++
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deletedCount": deletedCount,
		"freedBytes":   totalFreed,
		"freedStr":     cleaner.FormatBytes(totalFreed),
		"errors":       errors,
	})
}

func (s *Server) handleLeftoversScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanLeftovers(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleLeftoversClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IDs         []string `json:"ids"`
		MoveToTrash bool     `json:"moveToTrash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c := cleaner.NewCleaner(false)
	var totalFreed uint64
	var deletedCount int
	var errors []string

	for _, id := range req.IDs {
		if id == "all-ds-store" {
			freed, count, err := cleaner.CleanDSStoreBatch()
			if err != nil {
				errors = append(errors, fmt.Sprintf(".DS_Store temizleme hatası: %v", err))
			} else {
				totalFreed += freed
				deletedCount += int(count)
			}
			continue
		}

		freed, count, err := c.CleanTarget(id, req.MoveToTrash)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s silinemedi: %v", filepath.Base(id), err))
		} else {
			totalFreed += freed
			deletedCount += int(count)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deletedCount": deletedCount,
		"freedBytes":   totalFreed,
		"freedStr":     cleaner.FormatBytes(totalFreed),
		"errors":       errors,
	})
}

// ----------------------------------------------------
// App Uninstaller Handlers
// ----------------------------------------------------

func (s *Server) handleAppsScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanInstalledApplications(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type appUninstallRequest struct {
	AppID     string `json:"appId"`
	ResetOnly bool   `json:"resetOnly"`
	UseTrash  bool   `json:"useTrash"`
}

func (s *Server) handleAppUninstall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req appUninstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek gövdesi", http.StatusBadRequest)
		return
	}

	freed, err := cleaner.UninstallApp(req.AppID, req.ResetOnly, req.UseTrash)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	action := "kaldırıldı"
	if req.ResetOnly {
		action = "sıfırlandı"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"freedBytes": freed,
		"freedStr":   cleaner.FormatBytes(freed),
		"message":    fmt.Sprintf("Uygulama başarıyla %s. %s yer açıldı.", action, cleaner.FormatBytes(freed)),
	})
}

// ----------------------------------------------------
// Apple & System Data Handlers
// ----------------------------------------------------

func (s *Server) handleAppleScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanAppleSystem(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleAppleSnapshotsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SnapshotID string `json:"snapshotId"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if req.SnapshotID != "" && req.SnapshotID != "all" {
		if err := cleaner.DeleteAPFSSnapshot(req.SnapshotID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "APFS anlık görüntüsü silindi."})
		return
	}

	deleted, err := cleaner.DeleteAllAPFSSnapshots()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "deletedCount": deleted, "message": fmt.Sprintf("%d APFS anlık görüntüsü silindi.", deleted)})
}

func (s *Server) handleAppleBackupsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path     string `json:"path"`
		UseTrash bool   `json:"useTrash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	var err error
	if req.UseTrash {
		err = cleaner.MoveToTrash(req.Path)
	} else {
		err = os.RemoveAll(req.Path)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Aygıt yedeklemesi silindi."})
}

func (s *Server) handleAppleSimulatorsClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := cleaner.CleanSimulators(); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Kullanılmayan iOS simülatörleri ve önbellekleri temizlendi."})
}

// ----------------------------------------------------
// Media & Messages Attachments Handlers
// ----------------------------------------------------

func (s *Server) handleMediaAttachmentsScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanMediaAttachments(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleMediaAttachmentsClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req cleaner.CleanAttachmentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz parametreler", http.StatusBadRequest)
		return
	}

	freed, count, err := cleaner.CleanAttachments(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"deletedCount": count,
		"freedBytes":   freed,
		"freedStr":     cleaner.FormatBytes(freed),
		"message":      fmt.Sprintf("%d ek temizlendi, %s alan açıldı.", count, cleaner.FormatBytes(freed)),
	})
}

// ----------------------------------------------------
// Downloads Organizer Handlers
// ----------------------------------------------------

func (s *Server) handleDownloadsScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanDownloadsDirectory(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleDownloadsClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Paths    []string `json:"paths"`
		UseTrash bool     `json:"useTrash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	freed, count, err := cleaner.CleanDownloadsBatch(req.Paths, req.UseTrash)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"deletedCount": count,
		"freedBytes":   freed,
		"freedStr":     cleaner.FormatBytes(freed),
		"message":      fmt.Sprintf("%d dosya silindi, %s alan kazanıldı.", count, cleaner.FormatBytes(freed)),
	})
}

// ----------------------------------------------------
// Browser Caches Handlers
// ----------------------------------------------------

func (s *Server) handleBrowsersScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanBrowserCaches(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleBrowsersClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Browsers []string `json:"browsers"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	freed, count, err := cleaner.CleanBrowserCaches(req.Browsers)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"deletedCount": count,
		"freedBytes":   freed,
		"freedStr":     cleaner.FormatBytes(freed),
		"message":      fmt.Sprintf("Tarayıcı önbellekleri temizlendi: %s kazanıldı.", cleaner.FormatBytes(freed)),
	})
}

// ----------------------------------------------------
// Smart Care Handlers
// ----------------------------------------------------

func (s *Server) handleSmartCareScan(w http.ResponseWriter, r *http.Request) {
	result, err := cleaner.ScanSmartCare(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleSmartCareClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SelectedIDs []string `json:"selectedIds"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	freed, count, err := cleaner.ExecuteSmartCareClean(req.SelectedIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"deletedCount": count,
		"freedBytes":   freed,
		"freedStr":     cleaner.FormatBytes(freed),
		"message":      fmt.Sprintf("Akıllı Bakım tamamlandı! Toplam %s güvenle temizlendi.", cleaner.FormatBytes(freed)),
	})
}

// ----------------------------------------------------
// Authentication Handlers
// ----------------------------------------------------

func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	isAuth := s.isAuthenticated(r)
	writeJSON(w, http.StatusOK, map[string]any{
		"authEnabled":   cfg.Auth.Enabled,
		"authenticated": !cfg.Auth.Enabled || isAuth,
	})
}

func (s *Server) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	cfg := config.GetConfig()
	if !cfg.Auth.Enabled {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "Giriş başarılı"})
		return
	}

	if req.Password != cfg.Auth.Password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Hatalı şifre. Lütfen tekrar deneyin."})
		return
	}

	token := s.generateToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "dc_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"token":   token,
		"message": "Giriş başarılı.",
	})
}

func (s *Server) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("dc_token")
	if err == nil && cookie != nil {
		s.revokeToken(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "dc_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

// ----------------------------------------------------
// Maintenance & Speedup Toolkit Handlers
// ----------------------------------------------------

func (s *Server) handleMaintenanceList(w http.ResponseWriter, r *http.Request) {
	tasks := cleaner.GetAvailableMaintenanceTasks()
	writeJSON(w, http.StatusOK, tasks)
}

func (s *Server) handleMaintenanceRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TaskID string `json:"taskId"`
		ID     string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	targetID := req.TaskID
	if targetID == "" {
		targetID = req.ID
	}

	res, err := cleaner.ExecuteMaintenanceTask(r.Context(), targetID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// ----------------------------------------------------
// File Shredder Handler
// ----------------------------------------------------

func (s *Server) handleShred(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path   string `json:"path"`
		Passes int    `json:"passes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz parametreler"})
		return
	}

	if req.Passes <= 0 {
		req.Passes = 3
	}

	res, err := cleaner.ShredPath(req.Path, req.Passes)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// ----------------------------------------------------
// Hardware Resource Monitor Handler
// ----------------------------------------------------

func (s *Server) handleHardwareMonitor(w http.ResponseWriter, r *http.Request) {
	data, err := cleaner.GetHardwareMonitorData()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, data)
}

// ----------------------------------------------------
// Startup & Background Items Handlers
// ----------------------------------------------------

func (s *Server) handleStartupList(w http.ResponseWriter, r *http.Request) {
	summary, err := cleaner.ScanStartupItems(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) handleStartupToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.ToggleStartupItem(req.ID, req.Enabled); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	statusMsg := "Öğe devre dışı bırakıldı."
	if req.Enabled {
		statusMsg = "Öğe başarıyla etkinleştirildi."
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": statusMsg,
	})
}

func (s *Server) handleStartupDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID       string `json:"id"`
		UseTrash bool   `json:"useTrash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.RemoveStartupItem(req.ID, req.UseTrash); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Öğe başlangıç listesinden kaldırıldı.",
	})
}

func (s *Server) handleStartupAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path   string `json:"path"`
		Hidden bool   `json:"hidden"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.AddLoginItem(req.Path, req.Hidden); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Uygulama oturum açma öğelerine başarıyla eklendi.",
	})
}

// ----------------------------------------------------
// System Extensions Handlers
// ----------------------------------------------------

func (s *Server) handleExtensionsList(w http.ResponseWriter, r *http.Request) {
	summary, err := cleaner.ScanExtensions(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) handleExtensionsToggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID      string `json:"id"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.ToggleExtension(req.ID, req.Enabled); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	statusMsg := "Eklenti devre dışı bırakıldı."
	if req.Enabled {
		statusMsg = "Eklenti başarıyla etkinleştirildi."
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": statusMsg,
	})
}

func (s *Server) handleExtensionsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.DeleteExtension(req.ID, req.Path); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Eklenti başarıyla kaldırıldı.",
	})
}

// ----------------------------------------------------
// Network Monitor Handler
// ----------------------------------------------------

func (s *Server) handleNetworkStats(w http.ResponseWriter, r *http.Request) {
	stats, err := cleaner.GetNetworkStats(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// ----------------------------------------------------
// Privacy Protection Handlers
// ----------------------------------------------------

func (s *Server) handlePrivacyList(w http.ResponseWriter, r *http.Request) {
	summary, err := cleaner.ScanPrivacyTraces(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (s *Server) handlePrivacyClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ItemIDs []string `json:"itemIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	cleaned, err := cleaner.CleanPrivacyItems(req.ItemIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"cleaned": cleaned,
		"message": fmt.Sprintf("%d gizlilik öğesi başarıyla temizlendi.", cleaned),
	})
}

func (s *Server) handlePrivacyReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Service string `json:"service"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
		return
	}

	if err := cleaner.ResetPrivacyPermission(req.Service); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": fmt.Sprintf("%s izinleri başarıyla sıfırlandı.", req.Service),
	})
}





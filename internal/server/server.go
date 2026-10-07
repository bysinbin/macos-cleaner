package server

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"disk-cleaner/internal/config"
)

// Server handles the HTTP dashboard and API
type Server struct {
	port          int
	embeddedFS    embed.FS
	mux           *http.ServeMux
	sessions      map[string]time.Time
	sessionsMutex sync.RWMutex
}

// NewServer creates a new Server instance
func NewServer(port int, embeddedFS embed.FS) *Server {
	s := &Server{
		port:       port,
		embeddedFS: embeddedFS,
		mux:        http.NewServeMux(),
		sessions:   make(map[string]time.Time),
	}
	s.routes()
	return s
}

// Port returns the active listening port
func (s *Server) Port() int {
	return s.port
}

// Listen establishes the TCP listener on localhost, automatically finding an available port if needed
func (s *Server) Listen() (net.Listener, error) {
	cfg := config.GetConfig()
	bindHost := cfg.BindAddress
	if bindHost == "" {
		bindHost = "127.0.0.1"
	}

	addr := fmt.Sprintf("%s:%d", bindHost, s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		for p := s.port + 1; p <= s.port + 20; p++ {
			altAddr := fmt.Sprintf("%s:%d", bindHost, p)
			ln, err = net.Listen("tcp", altAddr)
			if err == nil {
				s.port = p
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("port dinlenemedi (%s): %w", addr, err)
		}
	}
	return ln, nil
}

// Serve handles incoming requests using the provided TCP listener
func (s *Server) Serve(ln net.Listener) error {
	cfg := config.GetConfig()
	bindHost := cfg.BindAddress
	if bindHost == "" {
		bindHost = "127.0.0.1"
	}

	addr := fmt.Sprintf("%s:%d", bindHost, s.port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.securityMiddleware(s.mux),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	fmt.Printf("\n🔒 Güvenlik: Yalnızca localhost (%s) üzerinden erişilebilir.\n", bindHost)
	if cfg.Auth.Enabled {
		fmt.Printf("🔑 Giriş Koruması: AKTİF (Şifreli oturum gerekiyor)\n")
	} else {
		fmt.Printf("🔓 Giriş Koruması: KAPALI\n")
	}
	fmt.Printf("🚀 Disk Cleaner Dashboard hazır: http://%s\n", addr)
	return srv.Serve(ln)
}

// Start launches the HTTP server with localhost-only check and optional auth middleware
func (s *Server) Start() error {
	ln, err := s.Listen()
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

// Security middleware enforcing localhost-only binding and authentication
func (s *Server) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Verify localhost-only access (reject any external incoming IPs)
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip := net.ParseIP(host)
			if ip != nil && !ip.IsLoopback() {
				http.Error(w, "Forbidden: Yalnızca localhost üzerinden erişilebilir.", http.StatusForbidden)
				return
			}
		}

		// 2. Auth protection check
		cfg := config.GetConfig()
		if cfg.Auth.Enabled {
			// Allow auth endpoints and static web assets without auth
			isAuthEndpoint := strings.HasPrefix(r.URL.Path, "/api/auth/")
			isAPIEndpoint := strings.HasPrefix(r.URL.Path, "/api/")

			if isAPIEndpoint && !isAuthEndpoint {
				if !s.isAuthenticated(r) {
					writeJSON(w, http.StatusUnauthorized, map[string]any{
						"error":        "Yetkilendirme gerekli. Lütfen şifrenizi girin.",
						"authRequired": true,
					})
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) getStableSecret() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "disk-cleaner-fallback-secret-2026"
	}
	secretFile := filepath.Join(home, ".disk-cleaner-secret")
	data, err := os.ReadFile(secretFile)
	if err == nil && len(data) >= 32 {
		return string(data)
	}
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	sec := hex.EncodeToString(b)
	_ = os.WriteFile(secretFile, []byte(sec), 0600)
	return sec
}

func (s *Server) generateToken() string {
	cfg := config.GetConfig()
	now := time.Now().Unix()
	payload := strconv.FormatInt(now, 10)
	mac := hmac.New(sha256.New, []byte(cfg.Auth.Password+":"+s.getStableSecret()))
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s.%s", payload, sig)
}

func (s *Server) revokeToken(token string) {
	// Stateless tokens are invalidated by client cookie clearing or password change
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	cfg := config.GetConfig()
	if !cfg.Auth.Enabled {
		return true
	}

	cookie, err := r.Cookie("dc_token")
	var token string
	if err == nil && cookie != nil {
		token = cookie.Value
	}
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		return false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	createdUnix, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}

	// 24-hour expiration window
	if time.Since(time.Unix(createdUnix, 0)) > 24*time.Hour {
		return false
	}

	mac := hmac.New(sha256.New, []byte(cfg.Auth.Password+":"+s.getStableSecret()))
	mac.Write([]byte(parts[0]))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(parts[1]), []byte(expectedSig))
}

// OpenBrowser opens the default web browser on macOS
func OpenBrowser(url string) {
	time.Sleep(200 * time.Millisecond)
	if runtime.GOOS == "darwin" {
		_ = exec.Command("open", url).Start()
	}
}

func (s *Server) routes() {
	// Authentication endpoints
	s.mux.HandleFunc("/api/auth/status", s.handleAuthStatus)
	s.mux.HandleFunc("/api/auth/login", s.handleAuthLogin)
	s.mux.HandleFunc("/api/auth/logout", s.handleAuthLogout)

	// Core API routes
	s.mux.HandleFunc("/api/system", s.handleSystemStats)
	s.mux.HandleFunc("/api/scan", s.handleScan)
	s.mux.HandleFunc("/api/clean", s.handleClean)
	s.mux.HandleFunc("/api/nodemodules", s.handleNodeModulesScan)
	s.mux.HandleFunc("/api/nodemodules/clean", s.handleNodeModulesClean)
	s.mux.HandleFunc("/api/largefiles", s.handleLargeFilesScan)
	s.mux.HandleFunc("/api/largefiles/delete", s.handleLargeFilesDelete)
	s.mux.HandleFunc("/api/tree", s.handleDirTree)
	s.mux.HandleFunc("/api/duplicates", s.handleDuplicates)
	s.mux.HandleFunc("/api/duplicates/delete", s.handleDuplicatesDelete)
	s.mux.HandleFunc("/api/leftovers", s.handleLeftoversScan)
	s.mux.HandleFunc("/api/leftovers/clean", s.handleLeftoversClean)
	s.mux.HandleFunc("/api/reveal", s.handleReveal)

	// App Uninstaller
	s.mux.HandleFunc("/api/apps", s.handleAppsScan)
	s.mux.HandleFunc("/api/apps/uninstall", s.handleAppUninstall)

	// Apple & System Data
	s.mux.HandleFunc("/api/apple", s.handleAppleScan)
	s.mux.HandleFunc("/api/apple/snapshots/delete", s.handleAppleSnapshotsDelete)
	s.mux.HandleFunc("/api/apple/backups/delete", s.handleAppleBackupsDelete)
	s.mux.HandleFunc("/api/apple/simulators/clean", s.handleAppleSimulatorsClean)
	s.mux.HandleFunc("/api/apple/systemdata/clean", s.handleAppleSystemDataClean)
	s.mux.HandleFunc("/api/apple/systemdata/clean-safe", s.handleAppleSystemDataCleanSafe)
	s.mux.HandleFunc("/api/apple/brew/cleanup", s.handleAppleBrewCleanup)
	s.mux.HandleFunc("/api/apple/purgeable/reclaim", s.handleApplePurgeableReclaim)

	// Media & Messages Attachments
	s.mux.HandleFunc("/api/media/attachments", s.handleMediaAttachmentsScan)
	s.mux.HandleFunc("/api/media/attachments/clean", s.handleMediaAttachmentsClean)

	// Downloads Organizer
	s.mux.HandleFunc("/api/downloads", s.handleDownloadsScan)
	s.mux.HandleFunc("/api/downloads/clean", s.handleDownloadsClean)

	// Browser Caches
	s.mux.HandleFunc("/api/browsers", s.handleBrowsersScan)
	s.mux.HandleFunc("/api/browsers/clean", s.handleBrowsersClean)

	// Smart Care & Health
	s.mux.HandleFunc("/api/smartcare", s.handleSmartCareScan)
	s.mux.HandleFunc("/api/smartcare/clean", s.handleSmartCareClean)

	// Maintenance Toolkit (Image 2: Free Up RAM, Flush DNS, Speed Up Mail, etc.)
	s.mux.HandleFunc("/api/maintenance", s.handleMaintenanceList)
	s.mux.HandleFunc("/api/maintenance/run", s.handleMaintenanceRun)

	// File Shredder (Image 1: Shredder)
	s.mux.HandleFunc("/api/shred", s.handleShred)

	// Hardware Resource Monitor (Image 3: CPU, RAM, Uptime)
	s.mux.HandleFunc("/api/monitor", s.handleHardwareMonitor)

	// Startup & Background Items Manager (CleanMyMac Optimization / Login Items)
	s.mux.HandleFunc("/api/startup", s.handleStartupList)
	s.mux.HandleFunc("/api/startup/toggle", s.handleStartupToggle)
	s.mux.HandleFunc("/api/startup/delete", s.handleStartupDelete)
	s.mux.HandleFunc("/api/startup/add", s.handleStartupAdd)

	// System Extensions Manager (SystemExtensions, PluginKit, QuickLook, Spotlight, PrefPanes)
	s.mux.HandleFunc("/api/extensions", s.handleExtensionsList)
	s.mux.HandleFunc("/api/extensions/toggle", s.handleExtensionsToggle)
	s.mux.HandleFunc("/api/extensions/delete", s.handleExtensionsDelete)
	s.mux.HandleFunc("/api/system/open-settings", s.handleOpenSettings)

	// Network Monitor (Interfaces, Throughput, Active TCP Sockets)
	s.mux.HandleFunc("/api/network", s.handleNetworkStats)

	// Privacy Protection (Recent Items, Terminal Histories, Browser Traces, TCC Reset)
	s.mux.HandleFunc("/api/privacy", s.handlePrivacyList)
	s.mux.HandleFunc("/api/privacy/clean", s.handlePrivacyClean)
	s.mux.HandleFunc("/api/privacy/reset", s.handlePrivacyReset)

	// Permissions & Full Disk Access
	s.mux.HandleFunc("/api/system/permissions", s.handlePermissionsStatus)
	s.mux.HandleFunc("/api/system/permissions/open-fda", s.handleOpenFDA)

	// Docker & Containers System
	s.mux.HandleFunc("/api/docker/status", s.handleDockerStatus)
	s.mux.HandleFunc("/api/docker/clean", s.handleDockerClean)

	// Live Scan Progress SSE Stream
	s.mux.HandleFunc("/api/events/progress", s.handleScanProgressStream)

	// Static UI assets from embedded FS
	sub, err := fs.Sub(s.embeddedFS, "web")
	if err != nil {
		s.mux.Handle("/", http.FileServer(http.Dir("./web")))
		return
	}
	s.mux.Handle("/", http.FileServer(http.FS(sub)))
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = jsonEncode(w, data)
	}
}

// Helper to set up SSE headers
func setupSSE(w http.ResponseWriter) (http.Flusher, bool) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return nil, false
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return flusher, true
}

func getCurrentUserHome() string {
	h, _ := os.UserHomeDir()
	return h
}

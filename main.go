package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"disk-cleaner/internal/cleaner"
	"disk-cleaner/internal/config"
	"disk-cleaner/internal/native"
	"disk-cleaner/internal/server"
	"disk-cleaner/internal/version"
)

//go:embed web/* web/css/* web/js/*
var embeddedWeb embed.FS

func init() {
	// Lock the main goroutine to OS thread 0 for Cocoa UI loop
	runtime.LockOSThread()
}

func main() {
	cfg, _ := config.LoadConfig()
	defaultPort := 8089
	if cfg != nil && cfg.Port > 0 {
		defaultPort = cfg.Port
	}

	port := flag.Int("port", defaultPort, "HTTP sunucu portu")
	browserMode := flag.Bool("browser", false, "Yerel pencere yerine varsayılan web tarayıcısında aç")
	headless := flag.Bool("headless", false, "Penceresiz / Arka plan sunucu modu")
	showVersion := flag.Bool("version", false, "Sürüm bilgisini göster")
	flag.Parse()

	if *showVersion {
		fmt.Printf("DiskCleaner Pro %s (%s, %s)\n", version.Version, version.GitCommit, version.BuildDate)
		os.Exit(0)
	}

	if cfg != nil && *port != defaultPort {
		cfg.Port = *port
		_ = config.SaveConfig(cfg)
	}

	fmt.Println("==================================================")
	fmt.Printf(" 🍏 DiskCleaner Pro %s — macOS Disk & Geliştirici Temizleyici\n", version.Version)
	fmt.Println("==================================================")

	srv := server.NewServer(*port, embeddedWeb)

	// Graceful shutdown on Interrupt / SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nDiskCleaner Pro kapatılıyor...")
		os.Exit(0)
	}()

	// Start Go HTTP Server in background
	go func() {
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Sunucu başlatılamadı: %v\n", err)
			os.Exit(1)
		}
	}()

	serverURL := fmt.Sprintf("http://127.0.0.1:%d", *port)

	// Health check to ensure local server is ready
	for i := 0; i < 30; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(serverURL + "/api/system")
		if err == nil && resp.StatusCode < 500 {
			_ = resp.Body.Close()
			break
		}
	}

	// Start Real-time Trash Watcher (AppCleaner behavior)
	trashWatcher := cleaner.NewTrashWatcher()
	trashWatcher.Start(context.Background(), func(appName, trashPath string) {
		cleaner.DefaultOnAppTrashed(appName, trashPath)
	})

	if *headless {
		fmt.Printf("DiskCleaner Pro arka planda çalışıyor: %s\n", serverURL)
		select {}
	} else if *browserMode {
		fmt.Printf("Tarayıcı modunda açılıyor: %s\n", serverURL)
		go server.OpenBrowser(serverURL)
		select {}
	} else {
		// DEFAULT: Native macOS Cocoa Window + Menubar Agent!
		fmt.Printf("macOS yerel penceresi ve menü çubuğu ajanı başlatılıyor (%s)...\n", serverURL)
		native.RunNativeApp(serverURL, fmt.Sprintf("DiskCleaner Pro %s", version.Version))
	}
}

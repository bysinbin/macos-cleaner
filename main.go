package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"disk-cleaner/internal/config"
	"disk-cleaner/internal/server"
)

//go:embed web/* web/css/* web/js/*
var embeddedWeb embed.FS

func main() {
	cfg, _ := config.LoadConfig()
	defaultPort := 8089
	if cfg != nil && cfg.Port > 0 {
		defaultPort = cfg.Port
	}

	port := flag.Int("port", defaultPort, "HTTP sunucu portu")
	noBrowser := flag.Bool("no-browser", false, "Tarayıcıyı otomatik açma")
	flag.Parse()

	if cfg != nil && *port != defaultPort {
		cfg.Port = *port
		_ = config.SaveConfig(cfg)
	}

	fmt.Println("==================================================")
	fmt.Println(" 🍏 DiskCleaner Pro — macOS Disk & Geliştirici Temizleyici")
	fmt.Println("==================================================")

	srv := server.NewServer(*port, embeddedWeb)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nDiskCleaner Pro kapatılıyor...")
		os.Exit(0)
	}()

	// Open browser automatically unless disabled
	if !*noBrowser {
		go server.OpenBrowser(fmt.Sprintf("http://127.0.0.1:%d", *port))
	}

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Sunucu başlatılamadı: %v\n", err)
		os.Exit(1)
	}
}

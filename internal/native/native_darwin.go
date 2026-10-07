//go:build darwin && !headless

package native

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework WebKit

#include "native_darwin.h"
*/
import "C"

import (
	"context"
	"fmt"
	"runtime"
	"time"
	"unsafe"

	"disk-cleaner/internal/cleaner"
)

//export goQuickFreeRAM
func goQuickFreeRAM() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		res, err := cleaner.ExecuteMaintenanceTask(ctx, "free-ram")
		msg := "Inaktif RAM blokları serbest bırakıldı."
		if err != nil {
			msg = "RAM temizliği sırasında uyarı: " + err.Error()
		} else if res != nil && res.Output != "" {
			msg = res.Output
		}
		SendNotification("⚡ Bellek Temizlendi", "DiskCleaner Pro", msg)
	}()
}

//export goQuickSmartCare
func goQuickSmartCare() {
	go func() {
		freed, count, err := cleaner.ExecuteSmartCareClean(nil)
		msg := "Akıllı Bakım tamamlandı."
		if err == nil {
			msg = fmt.Sprintf("%d dosya temizlendi, %s alan kazanıldı.", count, cleaner.FormatBytes(freed))
		}
		SendNotification("✨ Akıllı Bakım Tamamlandı", "DiskCleaner Pro", msg)
	}()
}

//export goOpenAppWindow
func goOpenAppWindow() {
	// Window brought to front by Cocoa action
}

// RunNativeApp initializes the native macOS Window with WKWebView and Menubar status item
func RunNativeApp(serverURL string, appTitle string) {
	runtime.LockOSThread()

	// Start background menubar RAM usage ticker
	go startMenubarTicker()

	cUrl := C.CString(serverURL)
	cTitle := C.CString(appTitle)
	defer C.free(unsafe.Pointer(cUrl))
	defer C.free(unsafe.Pointer(cTitle))

	C.runCocoaApp(cUrl, cTitle)
}

// ShowMainWindow brings the native window to front
func ShowMainWindow() {
	C.bringWindowToFront()
}

// SetMenubarTitle dynamically updates the tray item text
func SetMenubarTitle(title string) {
	cStr := C.CString(title)
	defer C.free(unsafe.Pointer(cStr))
	C.setMenubarTitle(cStr)
}

// SendNotification displays a native macOS User Notification
func SendNotification(title, subtitle, message string) {
	cTitle := C.CString(title)
	cSub := C.CString(subtitle)
	cMsg := C.CString(message)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cSub))
	defer C.free(unsafe.Pointer(cMsg))

	C.showNativeNotification(cTitle, cSub, cMsg)
}

func startMenubarTicker() {
	ticker := time.NewTicker(4 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		data, err := cleaner.GetHardwareMonitorData()
		if err == nil && data != nil && data.Memory.UsedPercent > 0 {
			title := fmt.Sprintf("🍏 RAM: %.0f%%", data.Memory.UsedPercent)
			SetMenubarTitle(title)
		}
	}
}

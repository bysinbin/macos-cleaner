//go:build !darwin || headless

package native

func RunNativeApp(serverURL string, appTitle string) {
	// Headless / non-darwin fallback
	select {}
}

func ShowMainWindow() {}

func SetMenubarTitle(title string) {}

func SendNotification(title, subtitle, message string) {}

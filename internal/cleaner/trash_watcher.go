package cleaner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// TrashWatcher monitors ~/.Trash for newly deleted .app applications
type TrashWatcher struct {
	knownApps map[string]bool
	mu        sync.Mutex
	stopCh    chan struct{}
}

// NewTrashWatcher creates a new watcher instance
func NewTrashWatcher() *TrashWatcher {
	return &TrashWatcher{
		knownApps: make(map[string]bool),
		stopCh:    make(chan struct{}),
	}
}

// Start begins background monitoring of ~/.Trash
func (tw *TrashWatcher) Start(ctx context.Context, onAppDeleted func(appName, trashPath string)) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	trashDir := filepath.Join(home, ".Trash")

	// Initial scan to populate already existing items in Trash
	tw.scanExisting(trashDir)

	ticker := time.NewTicker(2 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tw.stopCh:
				return
			case <-ticker.C:
				tw.checkTrash(trashDir, onAppDeleted)
			}
		}
	}()
}

func (tw *TrashWatcher) Stop() {
	select {
	case <-tw.stopCh:
	default:
		close(tw.stopCh)
	}
}

func (tw *TrashWatcher) scanExisting(trashDir string) {
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return
	}

	tw.mu.Lock()
	defer tw.mu.Unlock()
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			tw.knownApps[entry.Name()] = true
		}
	}
}

func (tw *TrashWatcher) checkTrash(trashDir string, callback func(appName, trashPath string)) {
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return
	}

	var newlyTrashed []struct {
		appName   string
		trashPath string
	}

	tw.mu.Lock()
	currentSet := make(map[string]bool)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".app") {
			currentSet[name] = true
			if !tw.knownApps[name] {
				tw.knownApps[name] = true
				appName := strings.TrimSuffix(name, ".app")
				newlyTrashed = append(newlyTrashed, struct {
					appName   string
					trashPath string
				}{
					appName:   appName,
					trashPath: filepath.Join(trashDir, name),
				})
			}
		}
	}

	// Clean up removed items from known set
	for known := range tw.knownApps {
		if !currentSet[known] {
			delete(tw.knownApps, known)
		}
	}
	tw.mu.Unlock()

	for _, item := range newlyTrashed {
		if callback != nil {
			callback(item.appName, item.trashPath)
		} else {
			DefaultOnAppTrashed(item.appName, item.trashPath)
		}
	}
}

// DefaultOnAppTrashed displays a native macOS User Notification
func DefaultOnAppTrashed(appName, trashPath string) {
	title := fmt.Sprintf("🗑️ %s Çöp Kutusuna Taşındı", appName)
	msg := fmt.Sprintf("%s uygulamasına ait artık dosyalar (Application Support, Caches) bulunabilir. Temizlemek için DiskCleaner'ı açın.", appName)

	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "Kaldırılan Uygulama Artıkları" sound name "Submarine"`,
		strings.ReplaceAll(msg, `"`, `\"`),
		strings.ReplaceAll(title, `"`, `\"`))

	_ = exec.Command("osascript", "-e", script).Run()
}

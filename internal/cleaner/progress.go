package cleaner

import (
	"sync"
	"time"
)

// ProgressEvent represents a live scanning/cleaning progress update
type ProgressEvent struct {
	Stage     string `json:"stage"`
	Message   string `json:"message"`
	Current   int64  `json:"current"`
	Total     int64  `json:"total"`
	Percent   int    `json:"percent"`
	Timestamp int64  `json:"timestamp"`
}

type progressHub struct {
	mu        sync.RWMutex
	listeners map[chan ProgressEvent]bool
}

var hub = &progressHub{
	listeners: make(map[chan ProgressEvent]bool),
}

// SubscribeProgress registers a new listener channel for progress events
func SubscribeProgress() chan ProgressEvent {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	ch := make(chan ProgressEvent, 64)
	hub.listeners[ch] = true
	return ch
}

// UnsubscribeProgress unregisters and closes a listener channel
func UnsubscribeProgress(ch chan ProgressEvent) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.listeners, ch)
	close(ch)
}

// EmitProgress broadcasts a progress event to all active SSE listeners
func EmitProgress(stage, message string, current, total int64) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	if len(hub.listeners) == 0 {
		return
	}

	var pct int
	if total > 0 {
		pct = int((float64(current) / float64(total)) * 100)
		if pct > 100 {
			pct = 100
		}
	}

	ev := ProgressEvent{
		Stage:     stage,
		Message:   message,
		Current:   current,
		Total:     total,
		Percent:   pct,
		Timestamp: time.Now().UnixMilli(),
	}

	for ch := range hub.listeners {
		select {
		case ch <- ev:
		default:
			// Non-blocking drop if channel is full
		}
	}
}

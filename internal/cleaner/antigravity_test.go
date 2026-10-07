package cleaner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsAntigravityMediaFile(t *testing.T) {
	mediaFiles := []string{
		"screenshot.png",
		"recording.webp",
		"preview.jpg",
		"capture.jpeg",
		"video.mp4",
		"clip.webm",
		"anim.gif",
	}
	for _, f := range mediaFiles {
		if !IsAntigravityMediaFile(f) {
			t.Errorf("expected %s to be recognized as media file", f)
		}
	}

	nonMediaFiles := []string{
		"transcript.jsonl",
		"transcript_full.jsonl",
		"task.md",
		"walkthrough.md",
		"metadata.json",
		"script.py",
		"main.go",
	}
	for _, f := range nonMediaFiles {
		if IsAntigravityMediaFile(f) {
			t.Errorf("expected %s NOT to be recognized as media file", f)
		}
	}
}

func TestBrainMediaScanAndClean(t *testing.T) {
	// Create a mock brain directory structure in temp
	tmpDir, err := os.MkdirTemp("", "mock-brain-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	convDir := filepath.Join(tmpDir, "mock-conv-1234")
	tempMediaDir := filepath.Join(convDir, ".tempmediaStorage")
	if err := os.MkdirAll(tempMediaDir, 0755); err != nil {
		t.Fatalf("failed to create conv dir: %v", err)
	}

	// 1. Create files that SHOULD be cleaned (older than 3 minutes)
	oldTime := time.Now().Add(-10 * time.Minute)

	media1 := filepath.Join(convDir, "phone_screenshot.png")
	if err := os.WriteFile(media1, []byte("fake-png-data-12345"), 0644); err != nil {
		t.Fatalf("failed to write media1: %v", err)
	}
	_ = os.Chtimes(media1, oldTime, oldTime)

	media2 := filepath.Join(tempMediaDir, "media_frame.jpg")
	if err := os.WriteFile(media2, []byte("fake-jpg-data-67890"), 0644); err != nil {
		t.Fatalf("failed to write media2: %v", err)
	}
	_ = os.Chtimes(media2, oldTime, oldTime)

	// 2. Create critical files that MUST NEVER be cleaned
	mdFile := filepath.Join(convDir, "task.md")
	if err := os.WriteFile(mdFile, []byte("# Task Plan\nKeep this!"), 0644); err != nil {
		t.Fatalf("failed to write mdFile: %v", err)
	}
	_ = os.Chtimes(mdFile, oldTime, oldTime)

	jsonlFile := filepath.Join(convDir, "transcript.jsonl")
	if err := os.WriteFile(jsonlFile, []byte("{\"step\":1}"), 0644); err != nil {
		t.Fatalf("failed to write jsonlFile: %v", err)
	}
	_ = os.Chtimes(jsonlFile, oldTime, oldTime)

	// Test Scan
	size, count, err := CalculateBrainMediaSize(tmpDir)
	if err != nil {
		t.Fatalf("CalculateBrainMediaSize error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 media files, got %d", count)
	}
	expectedSize := uint64(len("fake-png-data-12345") + len("fake-jpg-data-67890"))
	if size != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, size)
	}

	// Test Clean
	cleaner := NewCleaner(false)
	freed, delCount, err := cleaner.CleanBrainMedia(tmpDir, false)
	if err != nil {
		t.Fatalf("CleanBrainMedia error: %v", err)
	}
	if delCount != 2 {
		t.Errorf("expected 2 deleted files, got %d", delCount)
	}
	if freed != expectedSize {
		t.Errorf("expected freed %d, got %d", expectedSize, freed)
	}

	// Verify that media files were deleted
	if _, err := os.Stat(media1); !os.IsNotExist(err) {
		t.Errorf("media1 should have been deleted")
	}
	if _, err := os.Stat(media2); !os.IsNotExist(err) {
		t.Errorf("media2 should have been deleted")
	}

	// Verify that non-media files are STILL PRESENT and INTACT
	if _, err := os.Stat(mdFile); err != nil {
		t.Errorf("task.md MUST be preserved, but got error: %v", err)
	}
	if _, err := os.Stat(jsonlFile); err != nil {
		t.Errorf("transcript.jsonl MUST be preserved, but got error: %v", err)
	}
}

func TestAntigravityTargetsInPredefined(t *testing.T) {
	targets := GetPredefinedTargets()
	foundScreenshots := false
	foundRecordings := false
	foundBrowserCache := false
	foundScratch := false

	for _, target := range targets {
		switch target.ID {
		case "antigravity-screenshots":
			foundScreenshots = true
			if target.Category != "developer" {
				t.Errorf("expected category 'developer', got %s", target.Category)
			}
		case "antigravity-recordings":
			foundRecordings = true
		case "antigravity-browser-cache":
			foundBrowserCache = true
		case "antigravity-scratch":
			foundScratch = true
		}
	}

	if !foundScreenshots {
		t.Error("antigravity-screenshots target not found in predefined targets")
	}
	if !foundRecordings {
		t.Error("antigravity-recordings target not found in predefined targets")
	}
	if !foundBrowserCache {
		t.Error("antigravity-browser-cache target not found in predefined targets")
	}
	if !foundScratch {
		t.Error("antigravity-scratch target not found in predefined targets")
	}
}

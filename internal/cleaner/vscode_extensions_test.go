package cleaner

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	if !compareVersions("1.45.1", "1.46.0") {
		t.Errorf("expected 1.45.1 < 1.46.0")
	}
	if !compareVersions("2026.3.1", "2026.4.1") {
		t.Errorf("expected 2026.3.1 < 2026.4.1")
	}
	if !compareVersions("1.23.0-darwin-arm64", "1.24.0-darwin-arm64") {
		t.Errorf("expected 1.23.0 < 1.24.0")
	}
	if compareVersions("2.0.0", "1.9.9") {
		t.Errorf("expected 2.0.0 > 1.9.9")
	}
}

func TestFindVSCodeObsoleteExtensions(t *testing.T) {
	paths, err := FindVSCodeObsoleteExtensions()
	if err != nil {
		t.Fatalf("FindVSCodeObsoleteExtensions returned error: %v", err)
	}
	// On this user machine there are 10 obsolete extension directories (~824MB)
	if len(paths) == 0 {
		t.Logf("No obsolete extensions found (or .vscode/extensions is empty)")
	} else {
		t.Logf("Found %d obsolete extension versions", len(paths))
	}
}

package cleaner

import (
	"context"
	"testing"
)

func TestScanAppleSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow system scan in short mode")
	}
	result, err := ScanAppleSystem(context.Background())
	if err != nil {
		t.Fatalf("ScanAppleSystem failed: %v", err)
	}

	if result.MacOSInfo == nil {
		t.Fatal("expected MacOSInfo to not be nil")
	}

	t.Logf("macOS Total Size: %s, Status: %s, Version: %s", result.MacOSInfo.TotalSizeStr, result.MacOSInfo.Status, result.MacOSInfo.Version)
	t.Logf("Detected System Data Total: %s, Safe Cleanable: %s", result.TotalSystemDataStr, result.SafeCleanableStr)
	t.Logf("Categories count: %d", len(result.SystemDataCategories))

	for _, cat := range result.SystemDataCategories {
		t.Logf(" - [%s] %s: %s (%d items)", cat.ID, cat.Title, cat.TotalSizeStr, len(cat.Items))
		for _, item := range cat.Items {
			t.Logf("    * %s: %s (Risk: %s, Cleanable: %v)", item.Name, item.SizeStr, item.Risk, item.Cleanable)
		}
	}
}

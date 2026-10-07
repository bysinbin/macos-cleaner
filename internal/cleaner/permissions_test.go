package cleaner

import (
	"testing"
)

func TestCheckFullDiskAccess(t *testing.T) {
	status := GetSystemPermissions()
	t.Logf("Full Disk Access: %v, Message: %s", status.HasFullDiskAccess, status.FDAStatusMessage)
	if status.Platform != "darwin" {
		t.Errorf("expected platform darwin, got %s", status.Platform)
	}
}

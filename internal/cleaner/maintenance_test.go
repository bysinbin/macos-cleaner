package cleaner

import (
	"context"
	"strings"
	"testing"
)

func TestMaintenanceTasks(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live system maintenance tasks in short mode")
	}
	ctx := context.Background()

	// 1. Test free-ram
	resRAM, err := ExecuteMaintenanceTask(ctx, "free-ram")
	if err != nil {
		t.Fatalf("free-ram failed with error: %v", err)
	}
	if !resRAM.Success {
		t.Errorf("free-ram was expected to succeed")
	}
	if strings.Contains(resRAM.Output, "Operation not permitted") {
		t.Errorf("free-ram should not output unhandled permission errors, got: %s", resRAM.Output)
	}

	// 2. Test reindex-spotlight
	resSpotlight, err := ExecuteMaintenanceTask(ctx, "reindex-spotlight")
	if err != nil {
		t.Fatalf("reindex-spotlight failed with error: %v", err)
	}
	if !resSpotlight.Success {
		t.Errorf("reindex-spotlight was expected to succeed")
	}
	if strings.Contains(resSpotlight.Output, "Try as root") {
		t.Errorf("reindex-spotlight should gracefully handle non-root, got: %s", resSpotlight.Output)
	}

	// 3. Test run-periodic
	resPeriodic, err := ExecuteMaintenanceTask(ctx, "run-periodic")
	if err != nil {
		t.Fatalf("run-periodic failed with error: %v", err)
	}
	if !resPeriodic.Success {
		t.Errorf("run-periodic was expected to succeed")
	}

	// 4. Test speed-mail
	resMail, err := ExecuteMaintenanceTask(ctx, "speed-mail")
	if err != nil {
		t.Fatalf("speed-mail failed with error: %v", err)
	}
	if !resMail.Success {
		t.Errorf("speed-mail was expected to succeed")
	}
}

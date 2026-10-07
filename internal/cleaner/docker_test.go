package cleaner

import (
	"context"
	"testing"
)

func TestCheckDockerStatus(t *testing.T) {
	status := CheckDockerStatus(context.Background())
	t.Logf("Docker Installed: %v, Running: %v, Message: %s", status.Installed, status.Running, status.Message)
	if status.DockerRawSize > 0 {
		t.Logf("Docker.raw size: %s at %s", status.DockerRawSizeStr, status.DockerRawPath)
	}
}

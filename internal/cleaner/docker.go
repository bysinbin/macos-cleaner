package cleaner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DockerComponent represents individual container/image/volume metrics
type DockerComponent struct {
	Type        string `json:"type"`        // "Images", "Containers", "Volumes", "Build Cache"
	TotalCount  int    `json:"totalCount"`
	ActiveCount int    `json:"activeCount"`
	SizeStr     string `json:"sizeStr"`
	Reclaimable string `json:"reclaimable"`
	Size        uint64 `json:"size"`
}

// DockerStatus holds current Docker engine & storage status
type DockerStatus struct {
	Installed       bool              `json:"installed"`
	Running         bool              `json:"running"`
	Version         string            `json:"version"`
	DockerRawPath   string            `json:"dockerRawPath"`
	DockerRawSize   uint64            `json:"dockerRawSize"`
	DockerRawSizeStr string           `json:"dockerRawSizeStr"`
	Components      []DockerComponent `json:"components"`
	TotalReclaimable uint64           `json:"totalReclaimable"`
	TotalReclaimStr string            `json:"totalReclaimStr"`
	ColimaSizeStr   string            `json:"colimaSizeStr"`
	Message         string            `json:"message"`
}

// CheckDockerStatus inspects the Docker CLI, daemon, and physical storage files on macOS
func CheckDockerStatus(ctx context.Context) DockerStatus {
	status := DockerStatus{
		Components: []DockerComponent{},
	}

	// 1. Check physical Docker.raw file (Docker Desktop VM)
	home, _ := os.UserHomeDir()
	dockerRawPaths := []string{
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/data/Docker.raw"),
		filepath.Join(home, "Library/Containers/com.docker.docker/Data/vms/0/Docker.raw"),
	}

	for _, p := range dockerRawPaths {
		if fi, err := os.Stat(p); err == nil && fi.Size() > 0 {
			status.DockerRawPath = p
			status.DockerRawSize = uint64(fi.Size())
			status.DockerRawSizeStr = FormatBytes(status.DockerRawSize)
			break
		}
	}

	// 2. Check Colima VM if present
	colimaPath := filepath.Join(home, ".colima")
	if fi, err := os.Stat(colimaPath); err == nil && fi.IsDir() {
		scanner := NewScanner(ctx)
		cSz, _, _ := scanner.CalculateDirSize(colimaPath)
		if cSz > 0 {
			status.ColimaSizeStr = FormatBytes(cSz)
		}
	}

	// 3. Check docker executable
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		// Also check common macOS brew/application paths
		candidates := []string{
			"/usr/local/bin/docker",
			"/opt/homebrew/bin/docker",
			"/Applications/Docker.app/Contents/Resources/bin/docker",
		}
		for _, c := range candidates {
			if _, e := os.Stat(c); e == nil {
				dockerBin = c
				break
			}
		}
	}

	if dockerBin == "" {
		status.Installed = false
		status.Running = false
		if status.DockerRawSize > 0 {
			status.Message = fmt.Sprintf("Docker Desktop kurulu ancak CLI PATH'te değil. Sanal disk boyutu: %s", status.DockerRawSizeStr)
		} else {
			status.Message = "Sisteminizde aktif Docker kurulumu bulunamadı."
		}
		return status
	}

	status.Installed = true

	// 4. Check if Docker daemon is running
	verCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(verCtx, dockerBin, "version", "--format", "{{.Server.Version}}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil || strings.TrimSpace(out.String()) == "" {
		status.Running = false
		status.Message = "Docker kurulu fakat Docker servisi (daemon) çalışmıyor."
		return status
	}

	status.Running = true
	status.Version = strings.TrimSpace(out.String())
	status.Message = fmt.Sprintf("Docker Engine v%s aktif.", status.Version)

	// 5. Query docker system df --format json
	dfCtx, cancelDf := context.WithTimeout(ctx, 6*time.Second)
	defer cancelDf()

	dfCmd := exec.CommandContext(dfCtx, dockerBin, "system", "df", "--format", "{{json .}}")
	dfOut, err := dfCmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(dfOut)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(line), &raw); err == nil {
				compType, _ := raw["Type"].(string)
				totCountStr := fmt.Sprintf("%v", raw["TotalCount"])
				actCountStr := fmt.Sprintf("%v", raw["Active"])
				sizeStr, _ := raw["Size"].(string)
				reclaimStr, _ := raw["Reclaimable"].(string)

				var totCount, actCount int
				fmt.Sscanf(totCountStr, "%d", &totCount)
				fmt.Sscanf(actCountStr, "%d", &actCount)

				status.Components = append(status.Components, DockerComponent{
					Type:        compType,
					TotalCount:  totCount,
					ActiveCount: actCount,
					SizeStr:     sizeStr,
					Reclaimable: reclaimStr,
				})
			}
		}
	}

	return status
}

// CleanDocker executes targeted Docker cleanup
func CleanDocker(ctx context.Context, cleanType string) (string, error) {
	dockerBin, err := exec.LookPath("docker")
	if err != nil {
		candidates := []string{
			"/usr/local/bin/docker",
			"/opt/homebrew/bin/docker",
			"/Applications/Docker.app/Contents/Resources/bin/docker",
		}
		for _, c := range candidates {
			if _, e := os.Stat(c); e == nil {
				dockerBin = c
				break
			}
		}
	}

	if dockerBin == "" {
		return "", fmt.Errorf("Docker komut satırı aracı bulunamadı")
	}

	var args []string
	switch cleanType {
	case "all":
		// docker system prune -a -f --volumes
		args = []string{"system", "prune", "-a", "-f", "--volumes"}
	case "buildcache":
		// docker builder prune -a -f
		args = []string{"builder", "prune", "-a", "-f"}
	case "images":
		// docker image prune -a -f
		args = []string{"image", "prune", "-a", "-f"}
	case "containers":
		// docker container prune -f
		args = []string{"container", "prune", "-f"}
	case "volumes":
		// docker volume prune -f
		args = []string{"volume", "prune", "-f"}
	default:
		// docker system prune -f (dangling only)
		args = []string{"system", "prune", "-f"}
	}

	cmd := exec.CommandContext(ctx, dockerBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out)), fmt.Errorf("docker temizliği başarısız: %v (%s)", err, strings.TrimSpace(string(out)))
	}

	return strings.TrimSpace(string(out)), nil
}

package cleaner

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// NetworkConnection represents an active network socket connection
type NetworkConnection struct {
	Command     string `json:"command"`
	PID         int    `json:"pid"`
	User        string `json:"user"`
	Protocol    string `json:"protocol"`
	LocalAddr   string `json:"localAddr"`
	ForeignAddr string `json:"foreignAddr"`
	State       string `json:"state"`
}

// NetworkStats holds current network interfaces and traffic measurements
type NetworkStats struct {
	Interface         string              `json:"interface"`
	InterfaceName     string              `json:"interfaceName"`
	Status            string              `json:"status"`
	IPv4Address       string              `json:"ipv4Address"`
	Gateway           string              `json:"gateway"`
	DNS               []string            `json:"dns"`
	MACAddress        string              `json:"macAddress"`
	BytesIn           uint64              `json:"bytesIn"`
	BytesOut          uint64              `json:"bytesOut"`
	BytesInFormatted  string              `json:"bytesInFormatted"`
	BytesOutFormatted string              `json:"bytesOutFormatted"`
	DownloadSpeedStr  string              `json:"downloadSpeedStr"`
	UploadSpeedStr    string              `json:"uploadSpeedStr"`
	ActiveConnections []NetworkConnection `json:"activeConnections"`
	ConnectionCount   int                 `json:"connectionCount"`
}

var (
	lastNetCheckTime time.Time
	lastBytesIn      uint64
	lastBytesOut     uint64
	lastNetLock      sync.Mutex
)

// GetNetworkStats gathers all network interfaces, traffic counters, and active sockets
func GetNetworkStats(ctx context.Context) (*NetworkStats, error) {
	stats := &NetworkStats{
		Interface:     "en0",
		InterfaceName: "Wi-Fi / Yerel Ağ (en0)",
		Status:        "Bağlantı Yok",
		DNS:           make([]string, 0),
		ActiveConnections: make([]NetworkConnection, 0),
	}

	// 1. IP Address
	if out, err := exec.CommandContext(ctx, "ipconfig", "getifaddr", "en0").Output(); err == nil {
		ip := strings.TrimSpace(string(out))
		if ip != "" {
			stats.IPv4Address = ip
			stats.Status = "Bağlı"
		}
	}

	// 2. Gateway & Default Route
	if out, err := exec.CommandContext(ctx, "netstat", "-rn").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "default" {
				gw := fields[1]
				if !strings.HasPrefix(gw, "link#") && !strings.HasPrefix(gw, "fe80:") {
					stats.Gateway = gw
					break
				}
			}
		}
	}

	// 3. DNS Servers
	if out, err := exec.CommandContext(ctx, "scutil", "--dns").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		seen := make(map[string]bool)
		for _, line := range lines {
			if strings.Contains(line, "nameserver[") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					dnsIP := strings.TrimSpace(parts[1])
					if dnsIP != "" && !seen[dnsIP] {
						seen[dnsIP] = true
						stats.DNS = append(stats.DNS, dnsIP)
						if len(stats.DNS) >= 3 {
							break
						}
					}
				}
			}
		}
	}

	// 4. Traffic Bytes via `netstat -ib -n -I en0`
	if out, err := exec.CommandContext(ctx, "netstat", "-ib", "-n", "-I", "en0").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			// Header: Name Mtu Network Address Ipkts Ierrs Ibytes Opkts Oerrs Obytes Coll
			// We look for the link row which has MAC Address
			if len(fields) >= 11 && fields[0] == "en0" {
				if strings.Contains(fields[2], "<Link") {
					stats.MACAddress = fields[3]
					if bi, err := strconv.ParseUint(fields[6], 10, 64); err == nil {
						stats.BytesIn = bi
					}
					if bo, err := strconv.ParseUint(fields[9], 10, 64); err == nil {
						stats.BytesOut = bo
					}
				}
			}
		}
	}

	stats.BytesInFormatted = FormatBytes(stats.BytesIn)
	stats.BytesOutFormatted = FormatBytes(stats.BytesOut)

	// Calculate Speed
	lastNetLock.Lock()
	now := time.Now()
	if !lastNetCheckTime.IsZero() && lastBytesIn > 0 && stats.BytesIn >= lastBytesIn {
		sec := now.Sub(lastNetCheckTime).Seconds()
		if sec > 0.3 {
			deltaIn := stats.BytesIn - lastBytesIn
			deltaOut := stats.BytesOut - lastBytesOut

			speedIn := float64(deltaIn) / sec
			speedOut := float64(deltaOut) / sec

			stats.DownloadSpeedStr = formatSpeed(speedIn)
			stats.UploadSpeedStr = formatSpeed(speedOut)
		}
	} else {
		stats.DownloadSpeedStr = "0 KB/s"
		stats.UploadSpeedStr = "0 KB/s"
	}
	lastNetCheckTime = now
	lastBytesIn = stats.BytesIn
	lastBytesOut = stats.BytesOut
	lastNetLock.Unlock()

	// 5. Active TCP Connections via `lsof -iTCP -sTCP:ESTABLISHED -n -P`
	if out, err := exec.CommandContext(ctx, "lsof", "-iTCP", "-sTCP:ESTABLISHED", "-n", "-P").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		connMap := make(map[string]NetworkConnection)

		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			// COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
			if len(fields) >= 9 {
				cmd := fields[0]
				pid, _ := strconv.Atoi(fields[1])
				user := fields[2]
				proto := fields[4]
				nodeName := fields[8]

				// nodeName is like "192.168.254.95:50480->169.150.215.46:443"
				addrs := strings.Split(nodeName, "->")
				local := ""
				foreign := ""
				if len(addrs) == 2 {
					local = addrs[0]
					foreign = addrs[1]
				} else {
					local = nodeName
				}

				key := fmt.Sprintf("%d:%s:%s", pid, local, foreign)
				connMap[key] = NetworkConnection{
					Command:     cmd,
					PID:         pid,
					User:        user,
					Protocol:    proto,
					LocalAddr:   local,
					ForeignAddr: foreign,
					State:       "ESTABLISHED",
				}
			}
		}

		stats.ConnectionCount = len(connMap)
		limit := 50
		for _, conn := range connMap {
			stats.ActiveConnections = append(stats.ActiveConnections, conn)
			if len(stats.ActiveConnections) >= limit {
				break
			}
		}
	}

	return stats, nil
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 1024 {
		return fmt.Sprintf("%.0f B/s", bytesPerSec)
	} else if bytesPerSec < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSec/1024)
	} else if bytesPerSec < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB/s", bytesPerSec/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB/s", bytesPerSec/(1024*1024*1024))
}

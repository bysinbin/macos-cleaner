package cleaner

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"sync/atomic"
)

var extPattern = regexp.MustCompile(`^([a-zA-Z0-9_\-]+(?:\.[a-zA-Z0-9_\-]+)+)-(\d.*)$`)

// parseVersionChunks splits version strings into comparable segments (numbers or text)
func parseVersionChunks(ver string) []int {
	numRe := regexp.MustCompile(`\d+`)
	matches := numRe.FindAllString(ver, -1)
	chunks := make([]int, 0, len(matches))
	for _, m := range matches {
		if n, err := strconv.Atoi(m); err == nil {
			chunks = append(chunks, n)
		}
	}
	return chunks
}

// compareVersions returns true if v1 < v2
func compareVersions(v1, v2 string) bool {
	c1 := parseVersionChunks(v1)
	c2 := parseVersionChunks(v2)

	minLen := len(c1)
	if len(c2) < minLen {
		minLen = len(c2)
	}

	for i := 0; i < minLen; i++ {
		if c1[i] != c2[i] {
			return c1[i] < c2[i]
		}
	}
	return len(c1) < len(c2)
}

type extEntry struct {
	version string
	path    string
}

// FindVSCodeObsoleteExtensions finds old, unused versions of installed VS Code extensions
func FindVSCodeObsoleteExtensions() ([]string, error) {
	home, _ := os.UserHomeDir()
	extDir := filepath.Join(home, ".vscode", "extensions")

	entries, err := os.ReadDir(extDir)
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]extEntry)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		m := extPattern.FindStringSubmatch(name)
		if len(m) == 3 {
			baseID := m[1]
			ver := m[2]
			fullPath := filepath.Join(extDir, name)
			groups[baseID] = append(groups[baseID], extEntry{version: ver, path: fullPath})
		}
	}

	var obsoletePaths []string

	for _, list := range groups {
		if len(list) <= 1 {
			continue
		}

		// Sort by version ascending
		sort.Slice(list, func(i, j int) bool {
			return compareVersions(list[i].version, list[j].version)
		})

		// All except the last (highest) version are obsolete
		for i := 0; i < len(list)-1; i++ {
			obsoletePaths = append(obsoletePaths, list[i].path)
		}
	}

	return obsoletePaths, nil
}

// CalculateVSCodeObsoleteExtensionsSize computes size of old duplicate extensions
func CalculateVSCodeObsoleteExtensionsSize() (uint64, int64, error) {
	paths, err := FindVSCodeObsoleteExtensions()
	if err != nil {
		return 0, 0, err
	}

	var totalSize uint64
	var fileCount int64
	scanner := NewScanner(nil)

	for _, p := range paths {
		sz, cnt, _ := scanner.CalculateDirSize(p)
		atomic.AddUint64(&totalSize, sz)
		atomic.AddInt64(&fileCount, cnt)
	}

	return totalSize, fileCount, nil
}

// CleanVSCodeObsoleteExtensions removes old duplicate extension versions safely
func CleanVSCodeObsoleteExtensions() (uint64, int64, error) {
	paths, err := FindVSCodeObsoleteExtensions()
	if err != nil {
		return 0, 0, err
	}

	var freedBytes uint64
	var deletedFiles int64
	scanner := NewScanner(nil)

	for _, p := range paths {
		sz, cnt, _ := scanner.CalculateDirSize(p)
		if err := os.RemoveAll(p); err == nil {
			freedBytes += sz
			deletedFiles += cnt
		}
	}

	return freedBytes, deletedFiles, nil
}

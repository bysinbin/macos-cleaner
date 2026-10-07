package version

var (
	// Version is the current application release version
	Version = "v2.0.0"
	// BuildDate is populated at build time via -ldflags
	BuildDate = "2026-10-07"
	// GitCommit is populated at build time via -ldflags
	GitCommit = "HEAD"
)

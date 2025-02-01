package version

import "runtime"

var (
	Version = "0.5.0"

	// GoVersion is Go tree's version.
	GoVersion = runtime.Version()
)

// ConfigVersion is the current highest supported configuration version.
// Any configuration less than this version which has structural changes
// should migrate the configuration structures used by this version.
const ConfigVersion = 1

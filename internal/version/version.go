package version

import (
	"fmt"
	"runtime"

	gversion "github.com/arnaud-deprez/gsemver/pkg/gsemver/version"
)

var (
	// version is the current version of the gsemver.
	// Update this whenever making a new release.
	// The version is in string format and follow the semver 2 spec (Major.Minor.Patch[-Prerelease][+BuildMetadata])
	version = "0.1.0"

	// gitCommit is the git sha1
	gitCommit = ""
	// gitTreeState is the state of the git tree
	gitTreeState = ""
	// buildDate
	buildDate = ""
)

// GetVersion returns the version
func GetVersion() string {
	return version
}

// Get returns build info
func Get() gversion.BuildInfo {
	v := gversion.BuildInfo{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	return v
}

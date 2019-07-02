/*
Package version represents the current version of the project.

It is about the version of the tool and not the semver version implementation used by this tool.
*/
package version

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"gitCommit,omitempty"`
	// GitTreeState is the state of the git tree.
	// It is either clean or dirty.
	GitTreeState string `json:"gitTreeState,omitempty"`
	// BuildDate is the build date.
	BuildDate string `json:"buildDate,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"goVersion,omitempty"`
	// Compiler is the go compiler that built gsemver.
	Compiler string `json:"compiler,omitempty"`
	// Platform is the OS on which it is running.
	Platform string `json:"platform,omitempty"`
}

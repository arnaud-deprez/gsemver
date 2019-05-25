package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	/* const */ versionRegex *regexp.Regexp
	/* const */ zeroVersion = Version{}
	/* const */ versionBumperIdentity = func(v Version) Version { return v }
)

func init() {
	versionRegex = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)` +
		`(?:-([0-9A-Za-z\-]+(?:\.[0-9A-Za-z\-]+)*))?` +
		`(?:\+([0-9A-Za-z\-]+(?:\.[0-9A-Za-z\-]+)*))?$`)
}

// NewVersion creates a new Version from a string representation
func NewVersion(value string) (Version, error) {
	if value == "" {
		return zeroVersion, nil
	}

	m := versionRegex.FindStringSubmatch(value)
	if m == nil {
		return zeroVersion, Error{message: fmt.Sprintf("'%s' is not a semver compatible version", value)}
	}

	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])

	return Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    m[4],
		BuildMetadata: m[5],
	}, nil
}

// versionBumper type helper for the bump process
type versionBumper func(Version) Version

// Version object to represent a SemVer version
type Version struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
}

// String returns a string representation of a Version object.
// The format is: major.minor.patch[-pre_release_identifiers][+build_metadata]
func (v Version) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		sb.WriteString("-")
		sb.WriteString(v.PreRelease)
	}
	if v.BuildMetadata != "" {
		sb.WriteString("+")
		sb.WriteString(v.BuildMetadata)
	}
	return sb.String()
}

// IsUnstable returns true if the version is an early stage version. eg. 0.Y.Z
func (v *Version) IsUnstable() bool {
	return v.Major == 0
}

// BumpMajor bump the major number of the version
func (v Version) BumpMajor() Version {
	next := v
	// according to https://semver.org/#spec-item-11
	// Pre-release versions have a lower precedence than the associated normal version.
	// Build metadata SHOULD be ignored when determining version precedence.
	next.PreRelease = ""
	next.BuildMetadata = ""
	if v.PreRelease == "" || v.Minor != 0 || v.Patch != 0 {
		next.Major++
		next.Minor = 0
		next.Patch = 0
	}

	return next
}

// BumpMinor bumps the minor number of the version
func (v Version) BumpMinor() Version {
	next := v
	// according to https://semver.org/#spec-item-11
	// Pre-release versions have a lower precedence than the associated normal version.
	// Build metadata SHOULD be ignored when determining version precedence.
	next.PreRelease = ""
	next.BuildMetadata = ""
	if v.PreRelease == "" || v.Patch != 0 {
		next.Minor++
		next.Patch = 0
	}
	return next
}

// BumpPatch bumps the patch number of the version
func (v Version) BumpPatch() Version {
	next := v
	// according to https://semver.org/#spec-item-11
	// Pre-release versions have a lower precedence than the associated normal version.
	// Build metadata SHOULD be ignored when determining version precedence.
	next.PreRelease = ""
	next.BuildMetadata = ""
	if v.PreRelease == "" {
		next.Patch++
	}
	return next
}

// BumpPreRelease bumps the pre-release identifiers
func (v Version) BumpPreRelease(preRelease string, overwrite bool, semverBumper versionBumper) Version {
	// if no pre-release is define, just return the current version
	if preRelease == "" {
		return v
	}

	next := v

	if semverBumper == nil {
		// by default bump minor if this is not yet a pre-release
		semverBumper = Version.BumpMinor
	}
	// extract desired identifiers
	desiredIdentifiers := strings.Split(preRelease, ".")
	if !v.IsPreRelease() {
		// bump MAJOR, MINOR or PATCH if it's not yet a pre-release
		next = semverBumper(v)
	}

	if overwrite {
		next.PreRelease = preRelease
		return next
	}

	if v.IsPreRelease() {
		currentIdentifiers := strings.Split(v.PreRelease, ".")
		id, err := strconv.Atoi(currentIdentifiers[len(currentIdentifiers)-1])
		if arrayStringEqual(currentIdentifiers, desiredIdentifiers) ||
			(err == nil && arrayStringEqual(currentIdentifiers[:len(currentIdentifiers)-1], desiredIdentifiers)) {
			next.PreRelease = strings.Join(append(desiredIdentifiers, strconv.Itoa(id+1)), ".")
			return next
		}
		// TODO: else compare if pre-release name is >= than v.PreRelease
	}
	next.PreRelease = strings.Join(append(desiredIdentifiers, strconv.Itoa(0)), ".")
	return next
}

// IsPreRelease returns true if it's a pre-release version. eg 1.1.0-alpha.1
func (v Version) IsPreRelease() bool {
	return v.PreRelease != ""
}

// WithBuildMetadata return a new Version with build metadata
func (v Version) WithBuildMetadata(metadata string) Version {
	next := v
	next.BuildMetadata = metadata
	return next
}

// Error is a typical error representation that can happen during the version process
type Error struct {
	message string
	cause   error
}

// Error formats VersionError into a string
func (e Error) Error() string {
	if e.cause == nil {
		return e.message
	}
	return fmt.Sprintf("%s caused by: '%s'", e.message, e.cause)
}
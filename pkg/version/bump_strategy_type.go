package version

import (
	"encoding/json"
	"strings"
)

// BumpStrategyType represents the bump SemVer strategy to use to bump the version
type BumpStrategyType int

const (
	// PATCH means to bump the patch number
	PATCH BumpStrategyType = iota
	// MINOR means to bump the minor number
	MINOR
	// MAJOR means to bump the patch number
	MAJOR
	// AUTO means to apply the automatic strategy based on commit history
	AUTO
)

var bumpStrategyToString = []string{"PATCH", "MINOR", "MAJOR", "AUTO"}

// ParseBumpStrategyType converts string value to BumpStrategy
func ParseBumpStrategyType(value string) BumpStrategyType {
	switch strings.ToLower(value) {
	case "major":
		return MAJOR
	case "minor":
		return MINOR
	case "patch":
		return PATCH
	default:
		return AUTO
	}
}

func (b BumpStrategyType) String() string {
	return bumpStrategyToString[b]
}

// UnmarshalJSON implements unmarshall for encoding/json
func (b *BumpStrategyType) UnmarshalJSON(bs []byte) error {
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return err
	}
	*b = ParseBumpStrategyType(s)
	return nil
}

// MarshalJSON implements marshall for encoding/json
func (b BumpStrategyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

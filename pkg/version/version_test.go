package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionBumpMajor(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		data     string
		expected string
	}{
		{"0.1.0", "1.0.0"},
		{"v1.0.0-alpha.0", "1.0.0"}, // pre-release < release
		{"1.1.0-alpha.0", "2.0.0"},  // but not if we want to bump major on a minor pre-release
		{"v1.1.1-alpha.0", "2.0.0"}, // but not if we want to bump major on a patch pre-release
		{"v1.0.1-alpha.0", "2.0.0"},
	}
	for _, tc := range testData {
		v1, _ := NewVersion(tc.data)
		expected, _ := NewVersion(tc.expected)
		actual := v1.BumpMajor()
		assert.Equal(expected, actual)
	}
}

func TestVersionBumpMinor(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		data     string
		expected string
	}{
		{"0.1.0", "0.2.0"},
		{"1.0.1", "1.1.0"},
		{"v1.1.0-alpha.2", "1.1.0"}, // pre-release < release
		{"v1.0.1-alpha.2", "1.1.0"}, // but not if we want to bump minor on a patch pre-release
		{"1.1.1-alpha.2", "1.2.0"},  // same
	}
	for _, tc := range testData {
		v1, _ := NewVersion(tc.data)
		expected, _ := NewVersion(tc.expected)
		actual := v1.BumpMinor()
		assert.Equal(expected, actual)
	}
}

func TestVersionBumpPatch(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		data     string
		expected string
	}{
		{"0.1.0", "0.1.1"},
		{"0.1.0-alpha.0", "0.1.0"}, // pre-release < release
	}
	for _, tc := range testData {
		v1, _ := NewVersion(tc.data)
		expected, _ := NewVersion(tc.expected)
		actual := v1.BumpPatch()
		assert.Equal(expected, actual)
	}
}

func TestNewVersion(t *testing.T) {
	assert := assert.New(t)
	testData := []struct {
		version string
		err     bool
	}{
		{"1.2.3", false},
		{"v1.2.3", false},
		{"1.0", true},
		{"v1.0", true},
		{"1", true},
		{"v1", true},
		{"1.2.beta", true},
		{"v1.2.beta", true},
		{"foo", true},
		{"1.2-5", true},
		{"v1.2-5", true},
		{"1.2-beta.5", true},
		{"v1.2-beta.5", true},
		{"\n1.2", true},
		{"\nv1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"v1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hypen", false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", false},
		{"1.2.3-rc1-with-hypen", false},
		{"v1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", true},
		{"v1.2.3.4", true},
		{"1.2.2147483648", false},
		{"1.2147483648.3", false},
		{"2147483648.3.0", false},
	}

	for _, tc := range testData {
		_, err := NewVersion(tc.version)
		if tc.err && err == nil {
			assert.Fail("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			assert.Fail("error for version %s: %s", tc.version, err)
		}
	}
}

func TestBumpPreRelease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testData := []struct {
		version            string
		preRelease         string
		overridePreRelease bool
		expected           string
	}{
		{"1.0.0", "", false, "1.0.0"},
		{"1.0.0", "alpha", false, "1.1.0-alpha.0"},
		{"1.1.0-alpha.0", "alpha", false, "1.1.0-alpha.1"},
		{"1.1.0-alpha.1", "beta", false, "1.1.0-beta.0"},
		{"1.0.0", "", true, "1.0.0"},
		{"1.0.0", "SNAPSHOT", true, "1.1.0-SNAPSHOT"},
		{"1.1.0-SNAPSHOT", "SNAPSHOT", true, "1.1.0-SNAPSHOT"},
	}

	for _, tc := range testData {
		version, err := NewVersion(tc.version)
		assert.Nil(err)
		actual := version.BumpPreRelease(tc.preRelease, tc.overridePreRelease, nil)
		assert.Equal(tc.expected, actual.String())
	}
}

func TestWithBuildMetadata(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testData := []struct {
		version       string
		buildMetadata string
		expected      string
	}{
		{"1.0.0", "build.8", "1.0.0+build.8"},
		{"1.0.0", "3.abcdkd", "1.0.0+3.abcdkd"},
	}

	for _, tc := range testData {
		version, err := NewVersion(tc.version)
		assert.Nil(err)
		actual := version.WithBuildMetadata(tc.buildMetadata)
		assert.Equal(tc.expected, actual.String())
	}
}

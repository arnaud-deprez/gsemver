package cmd

import (
	"bytes"
	"os"
	"testing"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/stretchr/testify/assert"

	"github.com/arnaud-deprez/gsemver/internal/utils"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

func TestBumpNoFlag(t *testing.T) {
	testData := []struct {
		args             string
		expectedStrategy version.BumpStrategyType
	}{
		{"major", version.MAJOR},
		{"minor", version.MINOR},
		{"patch", version.PATCH},
		{"auto", version.AUTO},
		{"", version.AUTO},
	}

	for _, tc := range testData {
		t.Run(tc.args, func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			globalOpts := &globalOptions{
				ioStreams: newIOStreams(os.Stdin, out, errOut),
			}

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			root := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
				s := o.createBumpStrategy()

				assert.Equal(tc.expectedStrategy, s.Strategy)
				assert.Equal(version.DefaultMajorPattern, utils.RegexpToString(s.MajorPattern))
				assert.Equal(version.DefaultMinorPattern, utils.RegexpToString(s.MinorPattern))

				assert.Len(s.BumpBranchesStrategies, 1)
				assert.Equal(version.DefaultReleaseBranchesPattern, utils.RegexpToString(s.BumpBranchesStrategies[0].BranchesPattern))
				assert.False(s.BumpBranchesStrategies[0].PreRelease)
				assert.Equal("", utils.TemplateToString(s.BumpBranchesStrategies[0].PreReleaseTemplate))
				assert.False(s.BumpBranchesStrategies[0].PreReleaseOverwrite)
				assert.Equal("", utils.TemplateToString(s.BumpBranchesStrategies[0].BuildMetadataTemplate))

				assert.NotNil(s.BumpDefaultStrategy)
				assert.Equal(".*", utils.RegexpToString(s.BumpDefaultStrategy.BranchesPattern))
				assert.Equal(version.DefaultPreRelease, s.BumpDefaultStrategy.PreRelease)
				assert.Equal(version.DefaultPreReleaseTemplate, utils.TemplateToString(s.BumpDefaultStrategy.PreReleaseTemplate))
				assert.Equal(version.DefaultPreReleaseOverwrite, s.BumpDefaultStrategy.PreReleaseOverwrite)
				assert.Equal(version.DefaultBuildMetadataTemplate, utils.TemplateToString(s.BumpDefaultStrategy.BuildMetadataTemplate))

				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

func TestBumpChangePattern(t *testing.T) {
	testData := []struct {
		args                 string
		expectedMajorPattern string
		expectedMinorPattern string
	}{
		{`--major-pattern 'foo'`, "foo", version.DefaultMinorPattern},
		{`--minor-pattern 'bar'`, version.DefaultMajorPattern, "bar"},
		{`--major-pattern 'foo' --minor-pattern 'bar'`, "foo", "bar"},
	}

	for _, tc := range testData {
		t.Run(tc.args, func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			globalOpts := &globalOptions{
				ioStreams: newIOStreams(os.Stdin, out, errOut),
			}

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			root := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
				s := o.createBumpStrategy()

				assert.Equal(tc.expectedMajorPattern, utils.RegexpToString(s.MajorPattern))
				assert.Equal(tc.expectedMinorPattern, utils.RegexpToString(s.MinorPattern))
				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

func TestBumpPreRelease(t *testing.T) {
	testData := []struct {
		args                        string
		expectedPreRelease          bool
		expectedPreReleaseTemplate  string
		expectedPreReleaseOverwrite bool
	}{
		// TODO: would be nice to make this works {`--pre-release`, true, "", false},
		{`--pre-release ""`, true, "", false},
		{`--pre-release alpha`, true, "alpha", false},
		{`--pre-release SNAPSHOT --pre-release-overwrite`, true, "SNAPSHOT", true},
		// TODO: and this {`--pre-release --pre-release-overwrite`, true, "", true},
		{`--pre-release '' --pre-release-overwrite`, true, "", true},
	}

	for _, tc := range testData {
		t.Run(tc.args, func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			globalOpts := &globalOptions{
				ioStreams: newIOStreams(os.Stdin, out, errOut),
			}

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			root := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
				s := o.createBumpStrategy()

				assert.NotNil(s.BumpDefaultStrategy)
				assert.Equal(".*", utils.RegexpToString(s.BumpDefaultStrategy.BranchesPattern))
				assert.Equal(tc.expectedPreRelease, s.BumpDefaultStrategy.PreRelease)
				assert.Equal(tc.expectedPreReleaseTemplate, utils.TemplateToString(s.BumpDefaultStrategy.PreReleaseTemplate))
				assert.Equal(tc.expectedPreReleaseOverwrite, s.BumpDefaultStrategy.PreReleaseOverwrite)
				assert.Equal("{{.Commits | len}}.{{(.Commits | first).Hash.Short}}", utils.TemplateToString(s.BumpDefaultStrategy.BuildMetadataTemplate))

				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

func TestBumpBuildMetadata(t *testing.T) {
	testData := []struct {
		args                  string
		expectedBuildMetadata string
	}{
		{``, "{{.Commits | len}}.{{(.Commits | first).Hash.Short}}"},
		{`--build ""`, ""},
		{`--build "{{.Branch}}.{{(.Commits | first).Hash.Short}}"`, "{{.Branch}}.{{(.Commits | first).Hash.Short}}"},
	}

	for _, tc := range testData {
		t.Run(tc.args, func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			globalOpts := &globalOptions{
				ioStreams: newIOStreams(os.Stdin, out, errOut),
			}

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			root := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
				s := o.createBumpStrategy()

				assert.NotNil(s.BumpDefaultStrategy)
				assert.Equal(".*", utils.RegexpToString(s.BumpDefaultStrategy.BranchesPattern))
				assert.False(s.BumpDefaultStrategy.PreRelease)
				assert.Equal("", utils.TemplateToString(s.BumpDefaultStrategy.PreReleaseTemplate))
				assert.False(s.BumpDefaultStrategy.PreReleaseOverwrite)
				assert.Equal(tc.expectedBuildMetadata, utils.TemplateToString(s.BumpDefaultStrategy.BuildMetadataTemplate))

				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

func TestBumpBranchStrategy(t *testing.T) {
	testData := []struct {
		args                        string
		expectedBranchPattern       string
		expectedPreRelease          bool
		expectedPreReleaseTemplate  string
		expectedPreReleaseOverwrite bool
		expectedBuildMetadata       string
	}{
		{``, `^(master|release/.*)$`, false, "", false, ""},
		{`--branch-strategy '{"branchesPattern":".*","preRelease":true,"preReleaseTemplate":"foo","preReleaseOverwrite":true,"buildMetadataTemplate":"bar"}'`, `.*`, true, "foo", true, "bar"},
	}

	for _, tc := range testData {
		t.Run(tc.args, func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			globalOpts := &globalOptions{
				ioStreams: newIOStreams(os.Stdin, out, errOut),
			}

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			root := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
				s := o.createBumpStrategy()

				assert.Len(s.BumpBranchesStrategies, 1)
				assert.Equal(tc.expectedBranchPattern, utils.RegexpToString(s.BumpBranchesStrategies[0].BranchesPattern))
				assert.Equal(tc.expectedPreRelease, s.BumpBranchesStrategies[0].PreRelease)
				assert.Equal(tc.expectedPreReleaseTemplate, utils.TemplateToString(s.BumpBranchesStrategies[0].PreReleaseTemplate))
				assert.Equal(tc.expectedPreReleaseOverwrite, s.BumpBranchesStrategies[0].PreReleaseOverwrite)
				assert.Equal(tc.expectedBuildMetadata, utils.TemplateToString(s.BumpBranchesStrategies[0].BuildMetadataTemplate))

				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

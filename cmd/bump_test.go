package cmd

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/arnaud-deprez/gsemver/internal/utils"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

func TestBumpNoFlag(t *testing.T) {
	testData := []struct {
		args                         string
		expectedStrategy             version.BumpStrategyType
		expectedBumpBranchesStrategy []version.BumpBranchesStrategy
	}{
		{"major", version.MAJOR, []version.BumpBranchesStrategy{*version.NewBumpAllBranchesStrategy(version.MAJOR, false, "", false, "")}},
		{"minor", version.MINOR, []version.BumpBranchesStrategy{*version.NewBumpAllBranchesStrategy(version.MINOR, false, "", false, "")}},
		{"patch", version.PATCH, []version.BumpBranchesStrategy{*version.NewBumpAllBranchesStrategy(version.PATCH, false, "", false, "")}},
		{"auto", version.AUTO, []version.BumpBranchesStrategy{
			*version.NewDefaultBumpBranchesStrategy(version.DefaultReleaseBranchesPattern),
			*version.NewBuildBumpBranchesStrategy(".*", version.DefaultBuildMetadataTemplate),
		}},
		{"", version.AUTO, []version.BumpBranchesStrategy{
			*version.NewDefaultBumpBranchesStrategy(version.DefaultReleaseBranchesPattern),
			*version.NewBuildBumpBranchesStrategy(".*", version.DefaultBuildMetadataTemplate),
		}},
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

				assert.Equal(version.DefaultMajorPattern, utils.RegexpToString(s.MajorPattern))
				assert.Equal(version.DefaultMinorPattern, utils.RegexpToString(s.MinorPattern))
				assert.Equal(len(tc.expectedBumpBranchesStrategy), len(s.BumpStrategies))

				for i := range tc.expectedBumpBranchesStrategy {
					assert.Equal(tc.expectedBumpBranchesStrategy[i].GoString(), s.BumpStrategies[i].GoString())
				}

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

				assert.Len(s.BumpStrategies, 1)
				assert.Equal(".*", utils.RegexpToString(s.BumpStrategies[0].BranchesPattern))
				assert.Equal(tc.expectedPreRelease, s.BumpStrategies[0].PreRelease)
				assert.Equal(tc.expectedPreReleaseTemplate, utils.TemplateToString(s.BumpStrategies[0].PreReleaseTemplate))
				assert.Equal(tc.expectedPreReleaseOverwrite, s.BumpStrategies[0].PreReleaseOverwrite)
				assert.Equal("", utils.TemplateToString(s.BumpStrategies[0].BuildMetadataTemplate))

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
		args                         string
		expectedBumpBranchesStrategy []version.BumpBranchesStrategy
	}{
		{``, []version.BumpBranchesStrategy{
			*version.NewDefaultBumpBranchesStrategy(version.DefaultReleaseBranchesPattern),
			*version.NewBuildBumpBranchesStrategy(".*", version.DefaultBuildMetadataTemplate),
		}},
		{`--build-metadata ""`, []version.BumpBranchesStrategy{
			*version.NewBumpAllBranchesStrategy(version.AUTO, false, "", false, ""),
		}},
		{`--build-metadata "{{.Branch}}.{{(.Commits | first).Hash.Short}}"`, []version.BumpBranchesStrategy{
			*version.NewBumpAllBranchesStrategy(version.AUTO, false, "", false, "{{.Branch}}.{{(.Commits | first).Hash.Short}}"),
		}},
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

				assert.Equal(len(tc.expectedBumpBranchesStrategy), len(s.BumpStrategies))
				for i := range tc.expectedBumpBranchesStrategy {
					assert.Equal(tc.expectedBumpBranchesStrategy[i].GoString(), s.BumpStrategies[i].GoString())
				}

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

				size := 1
				if tc.args == "" {
					size = 2
				}
				assert.Len(s.BumpStrategies, size)
				assert.Equal(tc.expectedBranchPattern, utils.RegexpToString(s.BumpStrategies[0].BranchesPattern))
				assert.Equal(tc.expectedPreRelease, s.BumpStrategies[0].PreRelease)
				assert.Equal(tc.expectedPreReleaseTemplate, utils.TemplateToString(s.BumpStrategies[0].PreReleaseTemplate))
				assert.Equal(tc.expectedPreReleaseOverwrite, s.BumpStrategies[0].PreReleaseOverwrite)
				assert.Equal(tc.expectedBuildMetadata, utils.TemplateToString(s.BumpStrategies[0].BuildMetadataTemplate))

				return nil
			})
			globalOpts.addGlobalFlags(root)

			_, err = executeCommand(root, args...)
			assert.NoError(err)
		})
	}
}

func TestWithConfiguration(t *testing.T) {
	assert := assert.New(t)

	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	globalOpts := &globalOptions{
		ioStreams: newIOStreams(os.Stdin, out, errOut),
	}

	//args, err := shellquote.Split(tc.args)
	// assert.NoError(err)
	cmd := newBumpCommandsWithRun(globalOpts, func(o *bumpOptions) error {
		s := o.createBumpStrategy()

		assert.Equal("majorPatternConfig", s.MajorPattern.String(), "majorPattern does not match")
		assert.Equal("minorPatternConfig", s.MinorPattern.String(), "minorPattern does not match")
		expectedBumpBranchesStrategy := []version.BumpBranchesStrategy{
			{
				Strategy:        version.AUTO,
				BranchesPattern: regexp.MustCompile("releaseBranchesPattern"),
			},
			{
				Strategy:              version.AUTO,
				BranchesPattern:       regexp.MustCompile("all"),
				BuildMetadataTemplate: utils.NewTemplate("myBuildMetadataTemplate"),
			},
		}
		assert.Equal(len(expectedBumpBranchesStrategy), len(s.BumpStrategies))
		for i := range expectedBumpBranchesStrategy {
			assert.Equal(expectedBumpBranchesStrategy[i].GoString(), s.BumpStrategies[i].GoString())
		}

		return nil
	})
	globalOpts.addGlobalFlags(cmd)

	cobra.OnInitialize(func() {
		viper.SetConfigType("yaml")

		var yamlConfig = []byte(`
majorPattern: "majorPatternConfig"
minorPattern: "minorPatternConfig"
bumpStrategies:
- branchesPattern: "releaseBranchesPattern"
  strategy: "AUTO"
- branchesPattern: "all"
  strategy: "AUTO"
  buildMetadataTemplate: "myBuildMetadataTemplate"
`)
		err := viper.ReadConfig(bytes.NewBuffer(yamlConfig))
		assert.NoError(err, "Cannot read configuration")
	})

	_, err := executeCommand(cmd)
	assert.NoError(err)
}

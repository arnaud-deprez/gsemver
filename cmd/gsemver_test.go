package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// func emptyRun(*cobra.Command, ...string) {}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()
	return c, buf.String(), err
}

func TestConfigFile(t *testing.T) {
	assert := assert.New(t)
	for _, opt := range []string{"--config", "-c"} {
		t.Run(opt, func(_ *testing.T) {
			_, err := executeCommand(newDefaultRootCommand(), opt, "../test/data/gsemver-test-config.yaml")
			assert.NoError(err)
			assert.Equal("../test/data/gsemver-test-config.yaml", viper.ConfigFileUsed())
		})
	}
}

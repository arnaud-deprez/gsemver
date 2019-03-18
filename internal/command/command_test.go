package command

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	execCmd := New("vim --noplugin")
	assert.Equal(t, "vim", execCmd.Name)
	assert.Equal(t, 1, len(execCmd.Args))
	assert.Equal(t, "--noplugin", execCmd.Args[0])
}

func TestWithDir(t *testing.T) {
	execCmd := New("git")
	path, _ := filepath.Abs("")
	execCmd.InDir(path)
	assert.Equal(t, path, execCmd.Dir)
}

func TestWithArg(t *testing.T) {
	execCmd := New("git")
	execCmd.WithArg("command").WithArg("--amend").WithArg("-m").WithArg(`""`)
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 4, len(execCmd.Args))
}

func TestWithArgs(t *testing.T) {
	execCmd := New("git")
	//WithArgs reset the array
	execCmd.WithArg("command").WithArgs("--amend", "-m")
	assert.Equal(t, "git", execCmd.Name)
	assert.Equal(t, 2, len(execCmd.Args))
}

func TestWithEnvVariable(t *testing.T) {
	assert := assert.New(t)
	execCmd := New("git")
	execCmd.WithEnvVariable("FOO", "BAR")
	assert.Equal("git", execCmd.Name)
	assert.Contains(execCmd.Env, "FOO")
	assert.Equal("BAR", execCmd.Env["FOO"])
}

func TestWithEnv(t *testing.T) {
	assert := assert.New(t)
	execCmd := New("git")
	env := map[string]string{}
	env["TEST"] = "POC"
	//WithEnv reset the map
	execCmd.WithEnvVariable("FOO", "BAR").WithEnv(env)
	assert.Equal("git", execCmd.Name)
	assert.Equal(env, execCmd.Env)
}

func TestEcho(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	cmd := New("echo foo")
	out, err := cmd.Run()
	assert.Nil(err)
	assert.False(cmd.DidError())
	assert.Equal("foo", out)
}

func TestUnknownCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	cmd := New("unknownCommand")
	_, err := cmd.Run()
	assert.NotNil(err)
	assert.True(cmd.DidError())
	_, ok := cmd.Error().(Error)
	assert.True(ok)
}

func TestCommandTimedOut(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	cmd := New("sleep 2").WithTimeout(500 * time.Millisecond)
	out, err := cmd.Run()
	assert.NotNil(err)
	assert.True(cmd.DidError())
	assert.Equal(fmt.Sprintf("Failed to run '%s' command in directory '%s', output: '%s' caused by: 'Command timed out after %.2f seconds'", cmd.String(), cmd.Dir, out, cmd.Timeout.Seconds()), err.Error())
	_, ok := cmd.Error().(Error)
	assert.True(ok)
}

package command

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	shellquote "github.com/kballard/go-shellquote"

	"github.com/arnaud-deprez/gsemver/internal/log"
)

// Command is a struct containing the details of an external command to be executed
type Command struct {
	Name    string
	Args    []string
	Dir     string
	In      io.Reader
	Out     io.Writer
	Err     io.Writer
	Env     map[string]string
	Timeout time.Duration
	_error  error
}

// Error is the error object encapsulating an error from a Command
type Error struct {
	Command Command
	Output  string
	cause   error
}

func (c Error) Error() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Failed to run '%s %s' command in directory '%s', output: '%s'",
		c.Command.Name, strings.Join(c.Command.SanitisedArgs(), " "), c.Command.Dir, c.Output)
	if c.cause != nil {
		fmt.Fprintf(&sb, " caused by: '%s'", c.cause.Error())
	}
	return sb.String()
}

// New construct new command based on string
func New(cmd string) *Command {
	cmds, err := shellquote.Split(cmd)
	if err != nil {
		log.Error("Failed to parse command %s due to %s", cmd, err)
		// if we cannot parse the command, we should stop here and panic
		os.Exit(1)
	}

	return NewWithVarArgs(cmds...)
}

// NewWithVarArgs construct new command based on a string array
func NewWithVarArgs(cmd ...string) *Command {
	if len(cmd) == 0 {
		log.Error("Cannot instantiate an empty command!")
		// if we cannot parse the command, we should stop here and panic
		os.Exit(1)
	}
	return &Command{Name: cmd[0], Args: cmd[1:]}
}

// InDir Setter method for Dir to enable use of interface instead of Command struct
func (c *Command) InDir(dir string) *Command {
	c.Dir = dir
	return c
}

// WithArg sets an argument into the args
func (c *Command) WithArg(arg string) *Command {
	c.Args = append(c.Args, arg)
	return c
}

// WithArgs Setter method for Args to enable use of interface instead of Command struct
func (c *Command) WithArgs(args ...string) *Command {
	c.Args = args
	return c
}

// SanitisedArgs sanitises any password arguments before printing the error string.
// The actual sensitive argument is still present in the Command object
func (c *Command) SanitisedArgs() []string {
	sanitisedArgs := make([]string, len(c.Args))
	copy(sanitisedArgs, c.Args)
	for i, arg := range sanitisedArgs {
		if strings.Contains(strings.ToLower(arg), "password") && i < len(sanitisedArgs)-1 {
			// sanitise the subsequent argument to any 'password' fields
			sanitisedArgs[i+1] = "*****"
		}
	}
	return sanitisedArgs
}

// WithTimeout Setter method for Timeout to enable use of interface instead of Command struct
func (c *Command) WithTimeout(timeout time.Duration) *Command {
	c.Timeout = timeout
	return c
}

// getTimeout private getter method that returns the current timeout duration or 3 minutes by default
func (c *Command) getOrDefaultTimeout() time.Duration {
	// configure timeout, default is 3 minutes
	if c.Timeout == 0 {
		c.Timeout = 3 * time.Minute
	}
	return c.Timeout
}

// WithEnv Setter method for Env to enable use of interface instead of Command struct
func (c *Command) WithEnv(env map[string]string) *Command {
	c.Env = env
	return c
}

// WithEnvVariable sets an environment variable into the environment
func (c *Command) WithEnvVariable(name string, value string) *Command {
	if c.Env == nil {
		c.Env = map[string]string{}
	}
	c.Env[name] = value
	return c
}

// DidError returns a boolean if any error occurred in any execution of the command
func (c *Command) DidError() bool {
	return c._error != nil
}

// Error returns the last error
func (c *Command) Error() error {
	return c._error
}

// Run Execute the command without retrying on failure and block waiting for return values
func (c *Command) Run() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.getOrDefaultTimeout())
	// The cancel should be deferred so resources are cleaned up
	defer cancel()

	r, e := c.RunWithContext(&ctx)

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		err := Error{
			Command: *c,
			cause:   fmt.Errorf("Command timed out after %.2f seconds", c.Timeout.Seconds()),
		}
		c._error = err
		return "", err
	}

	if e != nil {
		c._error = e
	}
	return r, e
}

// String method returns a string representation of the Command
func (c Command) String() string {
	var builder strings.Builder
	builder.WriteString(c.Name)
	for _, arg := range c.Args {
		builder.WriteString(" ")
		builder.WriteString(arg)
	}
	return builder.String()
}

// RunWithContext private method executes the command and wait for the result
func (c *Command) RunWithContext(ctx *context.Context) (string, error) {
	e := exec.CommandContext(*ctx, c.Name, c.Args...)
	if c.Dir != "" {
		e.Dir = c.Dir
	}
	// merge env in e *Cmd
	if len(c.Env) > 0 {
		m := map[string]string{}
		environ := os.Environ()
		for _, kv := range environ {
			paths := strings.SplitN(kv, "=", 2)
			if len(paths) == 2 {
				m[paths[0]] = paths[1]
			}
		}
		for k, v := range c.Env {
			m[k] = v
		}
		envVars := []string{}
		for k, v := range m {
			envVars = append(envVars, k+"="+v)
		}
		e.Env = envVars
	}

	if c.Out != nil {
		e.Stdout = c.Out
	}

	if c.Err != nil {
		e.Stderr = c.Err
	}

	// zero value of string is ""
	var text string
	var err error

	if c.Out != nil {
		err := e.Run()
		if err != nil {
			return text, Error{
				Command: *c,
				cause:   err,
			}
		}
	} else {
		data, err := e.CombinedOutput()
		output := string(data)
		text = strings.TrimSpace(output)
		if err != nil {
			return text, Error{
				Command: *c,
				Output:  text,
				cause:   err,
			}
		}
	}

	return text, err
}

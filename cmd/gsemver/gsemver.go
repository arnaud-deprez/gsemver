package gsemver

// Run runs the command
func Run() error {
	cmd := newDefaultRootCommand()
	return cmd.Execute()
}

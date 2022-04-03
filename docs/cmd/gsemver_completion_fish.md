## gsemver completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	gsemver completion fish | source

To load completions for every new session, execute once:

	gsemver completion fish > ~/.config/fish/completions/gsemver.fish

You will need to start a new shell for this setup to take effect.


```
gsemver completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config string      config file (default is .gsemver.yaml)
      --log-level string   Sets the logging level (fatal, error, warning, info, debug, trace) (default "info")
  -v, --verbose            Enables verbose output by setting log level to debug. This is a shortland to --log-level debug.
```

### SEE ALSO

* [gsemver completion](gsemver_completion.md)	 - Generate the autocompletion script for the specified shell


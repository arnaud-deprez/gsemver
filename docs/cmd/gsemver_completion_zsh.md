## gsemver completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(gsemver completion zsh)

To load completions for every new session, execute once:

#### Linux:

	gsemver completion zsh > "${fpath[1]}/_gsemver"

#### macOS:

	gsemver completion zsh > $(brew --prefix)/share/zsh/site-functions/_gsemver

You will need to start a new shell for this setup to take effect.


```
gsemver completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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


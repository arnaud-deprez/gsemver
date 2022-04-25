## gsemver completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	gsemver completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
gsemver completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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


## gsemver completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(gsemver completion bash)

To load completions for every new session, execute once:

#### Linux:

	gsemver completion bash > /etc/bash_completion.d/gsemver

#### macOS:

	gsemver completion bash > /usr/local/etc/bash_completion.d/gsemver

You will need to start a new shell for this setup to take effect.


```
gsemver completion bash
```

### Options

```
  -h, --help              help for bash
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


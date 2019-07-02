## gsemver version

Print the CLI version information

### Synopsis


Show the version for gsemver.

This will print a representation the version of gsemver.
The output will look something like this:

version.BuildInfo{Version:"v0.1.0", GitCommit:"ff52399e51bb880526e9cd0ed8386f6433b74da1", GitTreeState:"clean"}

- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was
  built, and "dirty" if the binary was built from locally modified code.


```
gsemver version [flags]
```

### Options

```
  -h, --help    help for version
      --short   print the version number
```

### Options inherited from parent commands

```
      --log-level string   Sets the logging level (panic, fatal, error, warning, info, debug) (default "INFO")
      --verbose            Enables verbose output
```

### SEE ALSO

* [gsemver](gsemver.md)	 - CLI to manage semver compliant version from your git tags


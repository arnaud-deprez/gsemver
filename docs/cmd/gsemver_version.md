## gsemver version

Print the CLI version information

### Synopsis


Show the version for gsemver.

This will print a representation the version of gsemver.
The output will look something like this:

version.BuildInfo{Version:"0.1.0", GitCommit:"acfe51b15f9a1f12d47a20f88c29e5364916ae57", GitTreeState:"clean", BuildDate:"2019-07-02T07:44:00Z", GoVersion:"go1.12.6", Compiler:"gc", Platform:"darwin/amd64"}

- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was
  built, and "dirty" if the binary was built from locally modified code.
- BuildDate is the build date in ISO-8601 format at UTC.
- GoVersion is the go version with which it has been built.
- Compiler is the go compiler with which it has been built.
- Platform is the current OS platform on which it is running and for which it has been built.


```
gsemver version [flags]
```

### Examples

```

# Print version of gsemver
$ gsemver version

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


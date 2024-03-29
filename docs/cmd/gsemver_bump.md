## gsemver bump

Bump to next version

### Synopsis


This will compute and print the next semver compatible version of your project based on commits logs, tags and current branch.

The version will look like <X>.<Y>.<Z>[-<pre-release>][+<metadata>] where:
- X is the Major number
- Y is the Minor number
- Z is the Patch number
- pre-release is the pre-release identifiers (optional)
- metadata is the build metadata identifiers (optional)

More info on the semver spec https://semver.org/spec/v2.0.0.html.

It can work in 2 fashions, the automatic or manual.

Automatic way assumes: 
- your previous tags are semver compatible.
- you follow some conventions in your commit and ideally https://www.conventionalcommits.org
- you follow some branch convention for your releases (eg. a release should be done on main, master or release/* branches) 

Base on this information, it is able to compute the next version.

The manual way is less restrictive and just assumes your previous tags are semver compatible.


```
gsemver bump [strategy] [flags]
```

### Examples

```

# To bump automatically:
gsemver bump

# Or more explicitly
gsemver bump auto

# To bump manually the major number:
gsemver bump major

# To bump manually the minor number:
gsemver bump minor

# To bump manually the patch number:
gsemver bump patch

# To use a pre-release version
gsemver bump --pre-release alpha
# Or with go-template
gsemver bump --pre-release "alpha-{{.Branch}}"

# To use a pre-release version without indexation (maven like SNAPSHOT)
gsemver bump minor --pre-release SNAPSHOT --pre-release-overwrite true

# To use version with build metadata
gsemver bump --build-metadata "issue-1.build.1"
# Or with go-template
gsemver bump --build-metadata "{{(.Commits | first).Hash.Short}}"

# To use bump auto with one or many branch strategies
gsemver bump --branch-strategy='{"branchesPattern":"^miletone-1.1$","preReleaseTemplate":"beta"}' --branch-strategy='{"branchesPattern":"^miletone-2.0$","preReleaseTemplate":"alpha"}'

```

### Options

```
      --branch-strategy stringArray            Use branch-strategy will set a strategy for a set of branches. 
                                               The strategy is defined in json and looks like {"branchesPattern":"^milestone-.*$", "preReleaseTemplate":"alpha"} for example.
                                               This will use pre-release alpha version for every milestone-* branches. 
                                               You can find all available options https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#BumpBranchesStrategy
      --build-metadata string                  Use build metadata template which will give something like X.Y.Z+<build>.
                                               You can also use go-template expression with context https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#Context and http://masterminds.github.io/sprig functions.
                                               This flag cannot be used with --pre-release* flags and take precedence over them.
  -h, --help                                   help for bump
      --major-pattern string                   Use major-pattern option to define your regular expression to match a breaking change commit message
      --minor-pattern string                   Use major-pattern option to define your regular expression to match a minor change commit message
      --pre-release string                     Use pre-release template version such as 'alpha' which will give a version like 'X.Y.Z-alpha.N'.
                                               If pre-release flag is present but does not contain template value, it will give a version like 'X.Y.Z-N' where 'N' is the next pre-release increment for the version 'X.Y.Z'.
                                               You can also use go-template expression with context https://godoc.org/github.com/arnaud-deprez/gsemver/pkg/version#Context and http://masterminds.github.io/sprig functions.
                                               This flag is not taken into account if --build-metadata is set.
      --pre-release-overwrite X.Y.Z-SNAPSHOT   Use pre-release overwrite option to remove the pre-release identifier suffix which will give a version like X.Y.Z-SNAPSHOT if pre-release=SNAPSHOT
```

### Options inherited from parent commands

```
  -c, --config string      config file (default is .gsemver.yaml)
      --log-level string   Sets the logging level (fatal, error, warning, info, debug, trace) (default "info")
  -v, --verbose            Enables verbose output by setting log level to debug. This is a shortland to --log-level debug.
```

### SEE ALSO

* [gsemver](gsemver.md)	 - CLI to manage semver compliant version from your git tags


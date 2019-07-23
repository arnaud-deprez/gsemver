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
- you follow some branch convention for your releases (eg. a release should be done on master or release/* branches) 

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

# To use a pre-release version without indexation (maven like SNAPSHOT)
gsemver bump minor --pre-release SNAPSHOT --pre-release-overwrite true

# To use version with build metadata
gsemver bump --build "issue-1.build.1"

```

### Options

```
      --build string            Use build metadata which will give something like X.Y.Z+<build>
  -h, --help                    help for bump
      --pre-release string      Use pre-release version such as alpha which will give a version like X.Y.Z-alpha.N
      --pre-release-overwrite   Use pre-release overwrite option to remove the pre-release identifier suffix which will give a version like X.Y.Z-SNAPSHOT if pre-release=SNAPSHOT
```

### Options inherited from parent commands

```
      --log-level string   Sets the logging level (fatal, error, warning, info, debug, trace) (default "info")
  -v, --verbose            Enables verbose output by setting log level to debug. This is a shortland to --log-level debug.
```

### SEE ALSO

* [gsemver](gsemver.md)	 - CLI to manage semver compliant version from your git tags


# Contributing `gsemver`

Thank you for contributing `gsemver` :tada:

## Issue templates

Please use issue/PR templates for bugs or feature request which are inserted automatically.

If you have a question or need to raise another kind of issue, then choose custom.

Once you have raised an issue, you can also submit a Pull Request with a resolution.

## Commit Message Format

A format influenced by [Conventional commits](https://www.conventionalcommits.org).

```
<type>: <subject>
<BLANK LINE>
[body]
<BLANK LINE>
[footer]
```

### Type

Must be one of the following:

* **docs:** Documention only changes
* **ci:** Changes to our CI configuration files and scripts
* **chore:** Updating Makefile etc, no production code changes
* **feat:** A new feature
* **fix:** A bug fix
* **perf:** A code change that improves performance
* **refactor:** A code change that neither fixes a bug nor adds a feature
* **style:** Changes that do not affect the meaning of the code
* **test:** Adding missing tests or correcting existing tests

### Footer

The footer should contain a [closing reference to an issue](https://help.github.com/articles/closing-issues-via-commit-messages/) if any.

The **footer** should contain any information about **Breaking Changes** and is also the place to reference GitHub issues that this commit **Closes**.

**Breaking Changes** must start with the word `BREAKING CHANGE:` followed by a space and a description of it.

## Development

1. Fork (https://github.com/arnaud-deprez/gsemver) :tada:
1. Create a feature branch :coffee:
1. Run test suite with the `$ make test test-integration` command and confirm that it passes :zap:
1. Ensure the doc is up to date with your changes with the `$ make docs`command :+1:
1. Commit your changes :memo:
1. Rebase your local changes against the `master` branch and squash your commits if necessary :bulb:
1. Create new Pull Request :love_letter:

Bugs, feature requests and comments are more than welcome in the [issues](https://github.com/arnaud-deprez/gsemver/issues).
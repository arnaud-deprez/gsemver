# gsemver

gsemver is a command line tool developed in go that uses git commit convention to automate the generation of your next version compliant with [semver 2.0.0 spec](https://semver.org/spec/v2.0.0.html).

[![Build Status](https://travis-ci.com/arnaud-deprez/gsemver.svg?token=uYA1Qbgs8qpfX4ajYZsn&branch=master)](https://travis-ci.com/arnaud-deprez/gsemver)
[![GoDoc](https://godoc.org/github.com/arnaud-deprez/gsemver?status.svg)](https://godoc.org/github.com/arnaud-deprez/gsemver)
[![Downloads](https://img.shields.io/github/downloads/arnaud-deprez/gsemver/total.svg)](https://github.com/arnaud-deprez/gsemver/releases)
[![GoReport](https://goreportcard.com/badge/github.com/arnaud-deprez/gsemver)](https://goreportcard.com/report/github.com/arnaud-deprez/gsemver)
[![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/arnaud-deprez/gsemver/blob/master/LICENSE)

## Table of Contents

* [Motivations](#motivations)
* [Thanks](#thanks)
* [Getting Started](#getting-started)
  * [Installation](#installation)
    * [Go users](#go-users)
    * [Manual](#manual)
  * [Test Installation](#test-installation)
* [Usage](#usage)
  * [Pre-requisites](#pre-requisites)
  * [CLI](#cli)
    * [Automatic version bump](#automatic-version-bump)
    * [Manual version bump](#manual-version-bump)
  * [API](#api)
* [Contributing](#contributing)
  * [Feedback](#feedback)
* [License](#license)

## Motivations

Why yet another git version tool ?

When you try to implement DevOps pipeline for applications and libraries from different horizons (java, go, javascript, etc.), you always need to deal with versions from the moment you want to release your application/library to the deployment in production.

As DevOps is all about automation, you need a way to automate the generation of your next version.

Then, you have 2 choices:

1. you can use no human meaningful information:
    * forever increment a number
    * use git commit hash
    * use build number injected by your CI server
    * etc.
2. you can use a human meaningful convention such as [semver](https://semver.org/spec/v2.0.0.html).

The first option is easy and does not required any tool.

However some tools/tech require you to use a [semver](https://semver.org/spec/v2.0.0.html) compatible format version (eg. [go modules](https://github.com/golang/go/wiki/Modules), [helm](https://helm.sh/), etc.).
You can still decide to always bump the major, minor or patch number but then your version is not meaningful in you are just doing a hack to be compliant with the spec format but not with spec semantic.

So for the second option, in order to provide human meaningful information by following the spec semantic, you need to rely on some conventions.

You can find some git convention such as:

* [conventional commits](https://www.conventionalcommits.org): generalization of angular commit convention to other projects
* [angular commit convention](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#-commit-message-guidelines)
* [gitflow](https://datasift.github.io/gitflow/IntroducingGitFlow.html)

Then I looked for existing tools and here is a non exhaustive list of what I've found so far:

* [GitVersion](https://gitversion.readthedocs.io/en/latest/): tool written in .NET.
* [semantic-release](https://github.com/semantic-release/semantic-release): tool for npm
* [standard-version](https://github.com/conventional-changelog/standard-version): tool for npm
* [jgitver](https://github.com/jgitver/jgitver): CLI running on java, maven and gradle plugins.
* [hartym/git-semver](https://github.com/hartym/git-semver): git plugin written in python.
* [markchalloner/git-semver](https://github.com/markchalloner/git-semver): another git plugin written in bash
* [semver-maven-plugin](https://github.com/sidohaakma/semver-maven-plugin)

All these tools have at least one of these problems:

* They rely on a runtime environment (nodejs, python, java). But what if I want to build an application on another runtime ? On a VM, this is probably not a big deal but in a container where we try keep them as small as possible, this can be annoying.
* They are not designed to automatically generate a new version based on a convention. Instead, you have to specify what number you want to bump (major, minor, patch) and/or what type of version you want to generate (alpha, beta, build, etc.)
* They manage the full release lifecycle and so they are tightly coupled to some build tools like `npm`, `maven` or `gradle`.

I've found some libraries written in [go](https://golang.org/) but they don't deal with git commits/tags convention:

* [hashicorp/go-version](https://github.com/hashicorp/go-version)
* [coreos/go-semver](https://github.com/coreos/go-semver)
* [Masterminds/semver](https://github.com/Masterminds/semver)
* [blang/semver](https://github.com/blang/semver)

I needed a tool to generate the next release semver compatible version based on previous git tag that I could use on every type of application/library and so that is not relying on a specific runtime environment.

That's why I decided to build this tool using [go](https://golang.org/) with inspirations and credits from the tools I've found.

## Thanks

Thank you all for the inspirations!

I'd like also to thanks 2 projects that are used in combination with gsemver to better automate the release of this tool:

* [conventional commits](https://www.conventionalcommits.org) a commit convention I've decided to adopt in all my commits.
* [git-chglog](https://github.com/git-chglog/git-chglog) is a customizable CHANGELOG generator implemented in go based on commits log.
* [GoGeleaser](https://goreleaser.com) is a release automation tool for Go projects.

With these 3 tools and `gsemver`, it gets easier to automate the release your projects.

## Getting Started

### Installation

Please install `gsemver` in a way that matches your environment.

#### Go users

```sh
go get -u github.com/arnaud-deprez/gsemver/cmd
```

#### Manual

For a manual installation, you can download binary from [release page](https://github.com/arnaud-deprez/gsemver/releases) and place it in `$PATH` directory.

### Test Installation

You can check with the following command whether the `gsemver` command was included in a valid `$PATH`.

```bash
$ gsemver version
# output the gsemver version
```

## Usage

### Pre-requisites

Most of CI server uses - by default - [shallow git clone](https://git-scm.com/docs/git-clone#Documentation/git-clone.txt---depthltdepthgt) when cloning your git repository.

When performing such a clone, the local copy of your git repository will contain a _truncated history_ and most probably will be _detached from HEAD_.
As `gsemver` is currently using `git describe` and `git symbolic-ref`, it needs to have access to the whole history and HEAD needs to be attached.

This can cause error like `fatal: ref HEAD is not a symbolic ref` in case you need to google it for your CI server.

To overcome the situation, you can execute `git fetch --unshallow` before using `gsemver` or use custom options for your CI server: 

* [Travis options](https://docs.travis-ci.com/user/customizing-the-build#git-clone-depth) (which is used in this project)

### CLI

#### Automatic version bump

```sh
gsemver bump
```

This will use the git commits convention to generate the next version.

The only current supported convention is [conventional commits](https://www.conventionalcommits.org).
It also uses by default `master` and `release/*` branches by default as release branches and it generates version with build metadata for any branch that does not match.
This is a current limitation but the [roadmap](https://github.com/arnaud-deprez/gsemver/issues/4) is to make more configurable.

The [conventional commits integration tests](test/integration/gsemver_bump_auto_conventionalcommits_test.go) shows you in depth how version is generated.
For a more comprehension view, here an example of the logs graph these tests generate:

```git
*   34385d9 (HEAD -> master, tag: v1.2.2) Merge from feature/merge2-release-1.1.x
|\  
| *   b884197 Merge from release/1.1.x
| |\  
|/ /  
| *   869c83f (tag: v1.1.3, release/1.1.x) Merge from fix/fix-3
| |\  
| | * 22eabaf fix: my bug fix 3 on release/1.1.x
| |/  
* |   704fde4 (tag: v1.2.1) Merge from feature/merge-release-1.1.x
|\ \  
| * \   61b6a7c Merge from release/1.1.x
| |\ \  
|/ / /  
| | _   
| *   f2d9b5e (tag: v1.1.2) Merge from fix/fix-2
| |\  
| | * f95ccbe fix: my bug fix 2 on release/1.1.x
| |/  
* |   99a3662 (tag: v1.2.0) Merge from feature/awesome-3
|\ \  
| |/  
|/|   
| * cc6c1ed feat: my awesome 3rd change
|/  
*   145cbff (tag: v1.1.1) Merge from bug/fix-1
|\  
| * 681a11b fix: my bug fix on master
|/  
*   e9e7644 (tag: v1.1.0) Merge from feature/awesome-2
|\  
| * f30042e feat: my awesome 2nd change
|/  
*   fba50a2 (tag: v1.0.0, tag: v0.2.0) Merge from feature/awesome-1
|\  
| * bf05218 feat: my awesome change
|/  
* c619bff (tag: v0.1.1) fix(doc): fix documentation
* 128a5d9 (tag: v0.1.0) feat: add README.md
```

#### Manual version bump

```sh
gsemver bump major
gsemver bump minor
gsemver bump patch
```

All the CLI options are documented [here](docs/cmd/gsemver.md).

### API

For the API usage, you can check the [godoc](https://godoc.org/github.com/arnaud-deprez/gsemver) where there are some examples.

You can also check [version bumper release](internal/release/main.go) which is used to release gsemver itself.

## Contributing

We are always welcoming your contribution :clap:

But to make everyone's work easier, please read the [CONTRIBUTING guide](CONTRIBUTING.md) first.

### Feedback

I would like to make `gsemver` a better tool and take more scenario into account and eventually non conventional commits log.

Therefore, your feedback is very useful.  
I am very happy to hear your opinions on Issues and PR :heart:

## License

[MIT © Arnaud Deprez](./LICENSE)

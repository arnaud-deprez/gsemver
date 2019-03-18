# gsemver

gsemver is a command line tool developed in go that uses git commit convention to automate the generation of your next version compliant with [semver 2.0.0 spec](https://semver.org/spec/v2.0.0.html).

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

The first option is easy and does not need any tool but it will also prevent you using a lot of well known tools - such as [go modules](https://github.com/golang/go/wiki/Modules), [helm](https://helm.sh/), etc. - that are using [semver](https://semver.org/spec/v2.0.0.html) as a convention.

So for the second option, in order to provide human meaningful information, you need to rely on some conventions.

You can find some around and at the end there all based on git commit message convention:

* [angular commit convention](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#-commit-message-guidelines) 
* [conventional commits](https://www.conventionalcommits.org): generalization of angular commit convention to other projects
* [gitflow](https://datasift.github.io/gitflow/IntroducingGitFlow.html)

Then I looked for existing tools and here is what I've found so far:

* [GitVersion](https://gitversion.readthedocs.io/en/latest/): tool written in .NET.
* [semantic-release](https://github.com/semantic-release/semantic-release): tool for npm
* [standard-version](https://github.com/conventional-changelog/standard-version): tool for npm
* [jgitver](https://github.com/jgitver/jgitver): CLI running on java, maven and gradle plugins.
* [hartym/git-semver](https://github.com/hartym/git-semver): git plugin written in python.
* [markchalloner/git-semver](https://github.com/markchalloner/git-semver): another git plugin written in bash
* [semver-maven-plugin](https://github.com/sidohaakma/semver-maven-plugin)

All these tools have at least one of these problems:

* They rely on a runtime environment (nodejs, python, java). But what if I want to build an application on another runtime ?
* They are not designed to automatically generate a new version based on a convention. Instead, you have to specify what number you want to bump (major, minor, patch) or what type of version you want to generate (alpha, beta, build, etc.)
* They manage the full release lifecycle and so they are tightly coupled to some build tools like `npm`, `maven` or `gradle`.

I've found some libraries written in [go](https://golang.org/) but they don't deal with git tags:

* [hashicorp/go-version](https://github.com/hashicorp/go-version)
* [coreos/go-semver](https://github.com/coreos/go-semver)
* [Masterminds/semver](https://github.com/Masterminds/semver)
* [blang/semver](https://github.com/blang/semver)

I needed a tool to generate the next release semver compatible version based on previous git tag that I could use on every type of application/library and so that is not relying on a specific runtime environment.

That's why I decided to build this tool using [go](https://golang.org/) with inspirations from the tools I've found.

## Usage

The CLI usage is documented [here](docs/cmd/gsemver.md).

You can also use the api:

* [version bumper example](internal/release/main.go)
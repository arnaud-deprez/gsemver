package version_test

import (
	"fmt"
	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

func ExampleBumpStrategyOptions_Bump() {
	gitRepo := git.NewVersionGitRepo("dir")
	bumpStrategy := version.NewConventionalCommitBumpStrategyOptions(gitRepo)
	v, err := bumpStrategy.Bump()
	if err != nil {
		panic(err)
	}
	fmt.Println(v.String())
	// Use v like you want
}

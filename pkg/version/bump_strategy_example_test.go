package version_test

import (
	"fmt"

	"github.com/arnaud-deprez/gsemver/internal/git"
	"github.com/arnaud-deprez/gsemver/pkg/version"
)

func ExampleBumpStrategy_Bump() {
	gitRepo := git.NewVersionGitRepo("dir")
	bumpStrategy := version.NewConventionalCommitBumpStrategy(gitRepo)
	v, err := bumpStrategy.Bump()
	if err != nil {
		panic(err)
	}
	fmt.Println(v.String())
	// Use v like you want
}

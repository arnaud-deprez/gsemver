package gsemver

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/arnaud-deprez/gsemver/internal/version"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		args, expected string
	}{
		{"version", fmt.Sprintf("%#v\n", version.Get())},
		{"version --short", fmt.Sprintf("%s\n", version.GetVersion())},
	}

	for id, tc := range testCases {
		t.Run(fmt.Sprintf("TestVersion-%d", id), func(t *testing.T) {
			assert := assert.New(t)
			out, errOut := new(bytes.Buffer), new(bytes.Buffer)
			root := newRootCommand(os.Stdin, out, errOut)

			args, err := shellquote.Split(tc.args)
			assert.NoError(err)
			_, err = executeCommand(root, args...)
			assert.NoError(err)
			assert.Equal(tc.expected, out.String())
		})
	}
}

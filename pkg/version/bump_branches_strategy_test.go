package version

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBumpBranchesStrategyEncodingJson(t *testing.T) {
	assert := assert.New(t)

	testData := []struct {
		jsonVal string
		objVal  *BumpBranchesStrategy
	}{
		{
			`{"branchesPattern":"master","preRelease":false,"preReleaseOverwrite":false,"strategy":"AUTO"}`,
			NewDefaultBumpBranchesStrategy("master"),
		},
		{
			`{"branchesPattern":"milestone-1.2","preRelease":true,"preReleaseTemplate":"{{.Branch}}-foo","preReleaseOverwrite":false,"strategy":"AUTO"}`,
			NewPreReleaseBumpBranchesStrategy("milestone-1.2", "{{.Branch}}-foo", false),
		},
		{
			`{"branchesPattern":".*","preRelease":true,"preReleaseOverwrite":true,"buildMetadataTemplate":"{{.Branch}}.{{.Commits | len}}","strategy":"AUTO"}`,
			NewBumpAllBranchesStrategy(AUTO, true, "", true, "{{.Branch}}.{{.Commits | len}}"),
		},
	}

	for idx, tc := range testData {
		t.Run(fmt.Sprintf("Case %d Marshal", idx), func(t *testing.T) {
			out, err := json.Marshal(tc.objVal)
			assert.NoError(err)
			assert.JSONEq(tc.jsonVal, string(out))
		})

		t.Run(fmt.Sprintf("Case %d Unmarshal", idx), func(t *testing.T) {
			var out BumpBranchesStrategy
			err := json.Unmarshal([]byte(tc.jsonVal), &out)
			assert.NoError(err)

			if tc.objVal.BranchesPattern != nil {
				assert.Equal(tc.objVal.BranchesPattern.String(), out.BranchesPattern.String())
			}
			if tc.objVal.PreReleaseTemplate != nil {
				assert.Equal(tc.objVal.PreReleaseTemplate.Root.String(), out.PreReleaseTemplate.Root.String())
			}
			assert.Equal(tc.objVal.PreReleaseOverwrite, out.PreReleaseOverwrite)
			if tc.objVal.BuildMetadataTemplate != nil {
				assert.Equal(tc.objVal.BuildMetadataTemplate.Root.String(), out.BuildMetadataTemplate.Root.String())
			}
		})
	}
}

func ExampleBumpBranchesStrategy_GoString() {
	s := NewBumpBranchesStrategy(AUTO, ".*", true, "foo", true, "bar")
	fmt.Printf("%#v\n", s)
	// Output: version.BumpBranchesStrategy{Strategy: AUTO, BranchesPattern: &regexp.Regexp{expr: ".*"}, PreRelease: true, PreReleaseTemplate: &template.Template{text: "foo"}, PreReleaseOverwrite: true, BuildMetadataTemplate: &template.Template{text: "bar"}}
}

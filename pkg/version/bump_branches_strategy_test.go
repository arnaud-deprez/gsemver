package version

import (
	"encoding/json"
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
			`{"branchesPattern":".*","preRelease":true,"preReleaseTemplate":"{{.Branch}}-foo","preReleaseOverwrite":false}`,
			NewBumpBranchesStrategy(".*", true, "{{.Branch}}-foo", false, ""),
		},
		{
			`{"branchesPattern":".*","preRelease":true,"preReleaseOverwrite":true,"buildMetadataTemplate":"{{.Branch}}.{{.Commits | len}}"}`,
			NewBumpBranchesStrategy(".*", true, "", true, "{{.Branch}}.{{.Commits | len}}"),
		},
	}

	for _, tc := range testData {
		t.Run("Marshal", func(t *testing.T) {
			out, err := json.Marshal(tc.objVal)
			assert.NoError(err)
			assert.JSONEq(tc.jsonVal, string(out))
		})

		t.Run("Unmarshal", func(t *testing.T) {
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

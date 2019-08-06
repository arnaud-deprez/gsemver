package version

import (
	"encoding/json"
	"regexp"
	"text/template"
)

// NewBumpBranchesStrategy creates a new BumpBranchesStrategy
func NewBumpBranchesStrategy(pattern string, preReleaseTemplate string, preReleaseOverwrite bool, buildMetadataTemplate string) *BumpBranchesStrategy {
	return &BumpBranchesStrategy{
		BranchesPattern:       regexp.MustCompile(pattern),
		PreReleaseTemplate:    NewTemplate(preReleaseTemplate),
		PreReleaseOverwrite:   preReleaseOverwrite,
		BuildMetadataTemplate: NewTemplate(buildMetadataTemplate),
	}
}

// NewDefaultBumpBranchesStrategy creates a new BumpBranchesStrategy that matches all non matching branches.
func NewDefaultBumpBranchesStrategy(preReleaseTemplate string, preReleaseOverwrite bool, buildMetadataTemplate string) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(".*", preReleaseTemplate, preReleaseOverwrite, buildMetadataTemplate)
}

// NewBumpBranchesPreReleaseStrategy creates a new BumpBranchesStrategy for pre-release version strategy.
func NewBumpBranchesPreReleaseStrategy(pattern string, preReleaseTemplate string, preReleaseOverwrite bool) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(pattern, preReleaseTemplate, preReleaseOverwrite, "")
}

// NewBumpBranchesBuildStrategy creates a new BumpBranchesStrategy for build version strategy.
func NewBumpBranchesBuildStrategy(pattern string, buildMetadataTemplate string) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(pattern, "", false, buildMetadataTemplate)
}

// BumpBranchesStrategy allows you to configure the bump strategy option for a matching set of branches.
type BumpBranchesStrategy struct {
	// BranchesPattern is the regex used to match against the current branch
	BranchesPattern *regexp.Regexp `json:"branchesPattern,omitempty"`
	// PreReleaseTemplate defines the pre-release template for the next version
	// It can be alpha, beta, or a go-template expression
	PreReleaseTemplate *template.Template `json:"preReleaseTemplate,omitempty"`
	// PreReleaseOverwrite defines if a pre-release can be overwritten
	// If true, it will not append an index to the next version
	// If false, it will append an incremented index based on the previous same version of same class if any and 0 otherwise
	PreReleaseOverwrite bool `json:"preReleaseOverwrite"`
	// BuildMetadataTemplate defines the build metadata for the next version.
	// It can be a static value but it will usually be a go-template expression to guarantee uniqueness of each built version.
	BuildMetadataTemplate *template.Template `json:"buildMetadataTemplate,omitempty"`
}

// createVersionBumperFrom is an implementation for BumpBranchStrategy
func (s *BumpBranchesStrategy) createVersionBumperFrom(bumper versionBumper, ctx *Context) versionBumper {
	return func(v Version) Version {
		// build-metadata and pre-release are exclusives
		if s != nil && s.BuildMetadataTemplate != nil {
			return v.WithBuildMetadata(ctx.EvalTemplate(s.BuildMetadataTemplate))
		} else if s != nil && s.PreReleaseTemplate != nil {
			return v.BumpPreRelease(ctx.EvalTemplate(s.PreReleaseTemplate), s.PreReleaseOverwrite, bumper)
		}
		return bumper(v)
	}
}

// MarshalJSON implements json encoding
func (s *BumpBranchesStrategy) MarshalJSON() ([]byte, error) {
	type Alias BumpBranchesStrategy
	return json.Marshal(&struct {
		BranchesPattern       string `json:"branchesPattern,omitempty"`
		PreReleaseTemplate    string `json:"preReleaseTemplate,omitempty"`
		BuildMetadataTemplate string `json:"buildMetadataTemplate,omitempty"`
		*Alias
	}{
		BranchesPattern:       regexpToString(s.BranchesPattern),
		PreReleaseTemplate:    templateToString(s.PreReleaseTemplate),
		BuildMetadataTemplate: templateToString(s.BuildMetadataTemplate),
		Alias:                 (*Alias)(s),
	})
}

func regexpToString(r *regexp.Regexp) string {
	if r != nil {
		return r.String()
	}
	return ""
}

func templateToString(t *template.Template) string {
	if t != nil {
		return t.Root.String()
	}
	return ""
}

// UnmarshalJSON implements json decoding
func (s *BumpBranchesStrategy) UnmarshalJSON(data []byte) error {
	type Alias BumpBranchesStrategy
	aux := struct {
		BranchesPattern       string `json:"branchesPattern,omitempty"`
		PreReleaseTemplate    string `json:"preReleaseTemplate,omitempty"`
		BuildMetadataTemplate string `json:"buildMetadataTemplate,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.BranchesPattern = regexp.MustCompile(aux.BranchesPattern)
	s.PreReleaseTemplate = NewTemplate(aux.PreReleaseTemplate)
	s.BuildMetadataTemplate = NewTemplate(aux.BuildMetadataTemplate)
	return nil
}

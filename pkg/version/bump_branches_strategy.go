package version

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/arnaud-deprez/gsemver/internal/utils"
)

// NewBumpBranchesStrategy creates a new BumpBranchesStrategy
func NewBumpBranchesStrategy(strategy BumpStrategyType, pattern string, preRelease bool, preReleaseTemplate string, preReleaseOverwrite bool, buildMetadataTemplate string) *BumpBranchesStrategy {
	return &BumpBranchesStrategy{
		Strategy:              strategy,
		BranchesPattern:       regexp.MustCompile(pattern),
		PreRelease:            preRelease,
		PreReleaseTemplate:    utils.NewTemplate(preReleaseTemplate),
		PreReleaseOverwrite:   preReleaseOverwrite,
		BuildMetadataTemplate: utils.NewTemplate(buildMetadataTemplate),
	}
}

// NewBumpAllBranchesStrategy creates a new BumpBranchesStrategy that matches all branches.
func NewBumpAllBranchesStrategy(strategy BumpStrategyType, preRelease bool, preReleaseTemplate string, preReleaseOverwrite bool, buildMetadataTemplate string) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(strategy, ".*", preRelease, preReleaseTemplate, preReleaseOverwrite, buildMetadataTemplate)
}

// NewDefaultBumpBranchesStrategy creates a new BumpBranchesStrategy for pre-release version strategy.
func NewDefaultBumpBranchesStrategy(pattern string) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(AUTO, pattern, false, "", false, "")
}

// NewPreReleaseBumpBranchesStrategy creates a new BumpBranchesStrategy for pre-release version strategy.
func NewPreReleaseBumpBranchesStrategy(pattern string, preReleaseTemplate string, preReleaseOverwrite bool) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(AUTO, pattern, true, preReleaseTemplate, preReleaseOverwrite, "")
}

// NewBuildBumpBranchesStrategy creates a new BumpBranchesStrategy for build version strategy.
func NewBuildBumpBranchesStrategy(pattern string, buildMetadataTemplate string) *BumpBranchesStrategy {
	return NewBumpBranchesStrategy(AUTO, pattern, false, "", false, buildMetadataTemplate)
}

// BumpBranchesStrategy allows you to configure the bump strategy option for a matching set of branches.
type BumpBranchesStrategy struct {
	// Strategy defines the strategy to use to bump the version.
	// It can be automatic (AUTO) or manual (MAJOR, MINOR, PATCH)
	Strategy BumpStrategyType `json:"strategy"`
	// BranchesPattern is the regex used to match against the current branch
	BranchesPattern *regexp.Regexp `json:"branchesPattern,omitempty"`
	// PreRelease defines if the bump strategy should generate a pre-release version
	PreRelease bool `json:"preRelease"`
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
		}
		if s != nil && s.PreRelease {
			return v.BumpPreRelease(ctx.EvalTemplate(s.PreReleaseTemplate), s.PreReleaseOverwrite, bumper)
		}
		return bumper(v)
	}
}

// GoString makes BumpBranchesStrategy satisfy the GoStringer interface.
func (s BumpBranchesStrategy) GoString() string {
	var sb strings.Builder
	sb.WriteString("version.BumpBranchesStrategy{")
	sb.WriteString(fmt.Sprintf("Strategy: %v, ", s.Strategy))
	sb.WriteString(fmt.Sprintf("BranchesPattern: &regexp.Regexp{expr: %q}, ", s.BranchesPattern))
	sb.WriteString(fmt.Sprintf("PreRelease: %v, PreReleaseTemplate: &template.Template{text: %q}, PreReleaseOverwrite: %v, ", s.PreRelease, utils.TemplateToString(s.PreReleaseTemplate), s.PreReleaseOverwrite))
	sb.WriteString(fmt.Sprintf("BuildMetadataTemplate: &template.Template{text: %q}", utils.TemplateToString(s.BuildMetadataTemplate)))
	sb.WriteString("}")
	return sb.String()
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
		BranchesPattern:       utils.RegexpToString(s.BranchesPattern),
		PreReleaseTemplate:    utils.TemplateToString(s.PreReleaseTemplate),
		BuildMetadataTemplate: utils.TemplateToString(s.BuildMetadataTemplate),
		Alias:                 (*Alias)(s),
	})
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
	s.PreReleaseTemplate = utils.NewTemplate(aux.PreReleaseTemplate)
	s.BuildMetadataTemplate = utils.NewTemplate(aux.BuildMetadataTemplate)
	return nil
}

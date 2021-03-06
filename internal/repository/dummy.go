package repository

import (
	"github.com/drone/ff-mock-server/internal"
	"github.com/drone/ff-mock-server/pkg/api"
)

const (
	falseValue        = "false"
	segmentIdentifier = "demo"
)

var (
	trueValue           = "true"
	version       int64 = 1
	trueVariation       = api.Variation{
		Identifier: "true",
		Value:      trueValue,
	}
	falseVariation = api.Variation{
		Identifier: "false",
		Value:      falseValue,
	}
	featureConfig = api.FeatureConfig{
		Project: internal.Project,
		DefaultServe: api.Serve{
			Variation: &trueValue,
		},
		Environment:          internal.Environment,
		Feature:              "bool-flag",
		Kind:                 "boolean",
		OffVariation:         falseValue,
		Prerequisites:        nil,
		Rules:                nil,
		State:                api.FeatureStateOn,
		VariationToTargetMap: nil,
		Variations:           []api.Variation{trueVariation, falseVariation},
		Version:              &version,
	}
)

var (
	env     = internal.Environment
	segment = api.Segment{
		Environment: &env,
		Identifier:  segmentIdentifier,
		Excluded:    nil,
		Included:    nil,
		Name:        "Demo segment",
		Rules:       nil,
		Version:     &version,
	}
)

var (
	evaluation = api.Evaluation{
		Flag:       featureConfig.Feature,
		Identifier: &featureConfig.Feature,
		Kind:       string(featureConfig.Kind),
		Value:      "true",
	}
)

// DummyRepository contains mocked configurations, target groups and
// evaluations
type DummyRepository struct {
	featureConfigs map[string]api.FeatureConfig
	targetGroups   map[string]api.Segment
	evaluations    map[string]api.Evaluation
}

var _ Repository = DummyRepository{}

// NewDummyRepository returns new DummyRepository with initialized
// dummy data
func NewDummyRepository() *DummyRepository {
	return &DummyRepository{
		featureConfigs: map[string]api.FeatureConfig{
			featureConfig.Feature: featureConfig,
		},
		targetGroups: map[string]api.Segment{
			segment.Identifier: segment,
		},
		evaluations: map[string]api.Evaluation{
			evaluation.Flag: evaluation,
		},
	}
}

// GetFlagConfigurations returns all mocked configurations
// there is no need for env because environment is also mocked
func (r DummyRepository) GetFlagConfigurations() []api.FeatureConfig {
	slice := make([]api.FeatureConfig, 0, len(r.featureConfigs))
	for _, val := range r.featureConfigs {
		slice = append(slice, val)
	}
	return slice
}

// GetFlagConfiguration returns mocked configurations with identifier specified
// there is no need for env because environment is also mocked
func (r DummyRepository) GetFlagConfiguration(identifier string) (fc api.FeatureConfig, exists bool) {
	fc, exists = r.featureConfigs[identifier]
	return
}

// GetTargetGroups returns all mocked target groups
// there is no need for env because environment is also mocked
func (r DummyRepository) GetTargetGroups() []api.Segment {
	slice := make([]api.Segment, 0, len(r.targetGroups))
	for _, val := range r.targetGroups {
		slice = append(slice, val)
	}
	return slice
}

// GetTargetGroup returns mocked target group with identifier specified
// there is no need for env because environment is also mocked
func (r DummyRepository) GetTargetGroup(identifier string) (segment api.Segment, exists bool) {
	segment, exists = r.targetGroups[identifier]
	return
}

// GetEvaluations returns all mocked evaluations
// there is no need for env because environment is also mocked
func (r DummyRepository) GetEvaluations() api.Evaluations {
	slice := make([]api.Evaluation, 0, len(r.evaluations))
	for _, val := range r.evaluations {
		slice = append(slice, val)
	}
	return slice
}

// GetEvaluation returns mocked evaluation with identifier specified
// there is no need for env because environment is also mocked
func (r DummyRepository) GetEvaluation(identifier string) (evaluation api.Evaluation, exists bool) {
	evaluation, exists = r.evaluations[identifier]
	return
}

package repository

import "github.com/drone/ff-mock-server/pkg/api"

// Repository is used as an interface to access data
type Repository interface {
	GetFlagConfigurations() []api.FeatureConfig
	GetFlagConfiguration(identifier string) (fc api.FeatureConfig, exists bool)
	GetTargetGroups() []api.Segment
	GetTargetGroup(identifier string) (segment api.Segment, exists bool)
	GetEvaluations() api.Evaluations
	GetEvaluation(identifier string) (evaluation api.Evaluation, exists bool)
}

package router

import (
	"net/http"

	"github.com/drone/ff-mock-server/internal/repository"
	"github.com/drone/ff-mock-server/internal/service"
	"github.com/drone/ff-mock-server/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
)

// EventSource interface
type EventSource interface {
	CreateStream(id string) *sse.Stream
	StreamExists(id string) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Handler struct for implementing api methods
type Handler struct {
	eventSource EventSource
	repo        repository.Repository
}

// NewHandler returns new Handler struct with Repository and EventSource initialized
// using DIP
func NewHandler(repo repository.Repository, eventSource EventSource) *Handler {
	return &Handler{
		eventSource,
		repo,
	}
}

// Authenticate just check the mocked key and type of key
// and returns JWT token
func (h Handler) Authenticate(ctx echo.Context) error {
	authenticationRequest := api.AuthenticationRequest{}
	err := ctx.Bind(&authenticationRequest)
	if err != nil {
		return err
	}

	token, err := service.Authenticate(authenticationRequest.ApiKey)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, api.AuthenticationResponse{
		AuthToken: token,
	})
}

// GetFeatureConfig serve configuration array as JSON response
// environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetFeatureConfig(ctx echo.Context, environmentUUID string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetFlagConfigurations())
}

// GetFeatureConfigByIdentifier serve configuration specified with identifier
// environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetFeatureConfigByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	featureConfig, ok := h.repo.GetFlagConfiguration(identifier)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "feature not found",
		})
	}
	return ctx.JSON(http.StatusOK, featureConfig)
}

// GetAllSegments serve mocked target groups as JSON response
// environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetAllSegments(ctx echo.Context, environmentUUID string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetTargetGroups())
}

// GetSegmentByIdentifier serve mocked target group specified by identifier as JSON response
// environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetSegmentByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	segment, ok := h.repo.GetTargetGroup(identifier)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "segment not found",
		})
	}
	return ctx.JSON(http.StatusOK, segment)
}

// GetEvaluations serve evaluations as JSON response
// target and environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetEvaluations(ctx echo.Context, environmentUUID string, target string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetEvaluations())
}

// GetEvaluationByIdentifier serve evaluation as JSON response with specified feature
// target and environmentUUID not used because we are serving mocks for single environment
func (h Handler) GetEvaluationByIdentifier(ctx echo.Context, environmentUUID string, target string, feature string) error {
	evaluation, ok := h.repo.GetEvaluation(feature)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "segment not found",
		})
	}
	return ctx.JSON(http.StatusOK, evaluation)
}

// Stream is used to notify SDK instances using SSE
func (h Handler) Stream(ctx echo.Context, params api.StreamParams) error {
	req := ctx.Request()
	req.URL.RawQuery = "stream=" + params.APIKey
	if !h.eventSource.StreamExists(params.APIKey) {
		h.eventSource.CreateStream(params.APIKey)
	}
	h.eventSource.ServeHTTP(ctx.Response().Writer, req)
	return nil
}

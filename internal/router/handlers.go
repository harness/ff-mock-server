package router

import (
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/drone/ff-mock-server/internal/config"
	"github.com/drone/ff-mock-server/internal/repository"
	"github.com/drone/ff-mock-server/internal/service"
	"github.com/drone/ff-mock-server/pkg/api"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/r3labs/sse/v2"
)

// ErrAuthTokenNilOrInvalid ...
var ErrAuthTokenNilOrInvalid = errors.New("authorization token is either nil or incorrect type")

// EventSource interface
type EventSource interface {
	CreateStream(id string) *sse.Stream
	StreamExists(id string) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Close()
}

// Handler struct for implementing api methods
type Handler struct {
	eventSource        EventSource
	repo               repository.Repository
	targetDataReceived bool
	sseSeq             uint32
	sseTimeout         uint32
}

// NewHandler returns new Handler struct with Repository and EventSource initialized
// using DIP
func NewHandler(repo repository.Repository, eventSource EventSource) *Handler {
	return &Handler{
		eventSource: eventSource,
		repo:        repo,
	}
}

// Authenticate just check the mocked key and type of key
// and returns JWT token
func (h *Handler) Authenticate(ctx echo.Context) error {
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
func (h *Handler) GetFeatureConfig(ctx echo.Context, environmentUUID string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ServerKeyType, token); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, h.repo.GetFlagConfigurations())
}

// GetFeatureConfigByIdentifier serve configuration specified with identifier
// environmentUUID not used because we are serving mocks for single environment
func (h *Handler) GetFeatureConfigByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ServerKeyType, token); err != nil {
		return err
	}

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
func (h *Handler) GetAllSegments(ctx echo.Context, environmentUUID string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ServerKeyType, token); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, h.repo.GetTargetGroups())
}

// GetSegmentByIdentifier serve mocked target group specified by identifier as JSON response
// environmentUUID not used because we are serving mocks for single environment
func (h *Handler) GetSegmentByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ServerKeyType, token); err != nil {
		return err
	}
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
func (h *Handler) GetEvaluations(ctx echo.Context, environmentUUID string, target string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ClientKeyType, token); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, h.repo.GetEvaluations())
}

// GetEvaluationByIdentifier serve evaluation as JSON response with specified feature
// target and environmentUUID not used because we are serving mocks for single environment
func (h *Handler) GetEvaluationByIdentifier(ctx echo.Context, environmentUUID string, target string, feature string) error {
	token, ok := ctx.Get("user").(*jwt.Token)
	if !ok {
		return ErrAuthTokenNilOrInvalid
	}
	if err := service.CheckAPIKeyType(service.ClientKeyType, token); err != nil {
		return err
	}
	evaluation, ok := h.repo.GetEvaluation(feature)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "segment not found",
		})
	}
	return ctx.JSON(http.StatusOK, evaluation)
}

// Stream is used to notify SDK instances using SSEOffSequence
func (h *Handler) Stream(ctx echo.Context, params api.StreamParams) error {
	if timeout := atomic.LoadUint32(&h.sseSeq); timeout == 0 {
		return echo.NewHTTPError(500, "sse is in offline state")
	}
	log.Infof("connecting key %s on stream", params.APIKey)
	req := ctx.Request()
	req.URL.RawQuery = "stream=" + params.APIKey
	if !h.eventSource.StreamExists(params.APIKey) {
		h.eventSource.CreateStream(params.APIKey)
	}
	seq := atomic.LoadUint32(&h.sseSeq)
	timer := time.Tick(time.Duration(config.Options.SSEOffSequence[seq]) * time.Second)
	go func() {
		<-timer
		h.eventSource.Close()
	}()
	// blocking operation
	h.eventSource.ServeHTTP(ctx.Response().Writer, req)

	atomic.StoreUint32(&h.sseSeq, 0)
	if int(seq) < len(config.Options.SSEOffSequence)-1 {
		atomic.StoreUint32(&h.sseSeq, seq+1)
	}

	if config.Options.SSEOffDuration != nil {
		timeout := uint32(*config.Options.SSEOffDuration)
		atomic.StoreUint32(&h.sseTimeout, 1)
		time.Sleep(time.Duration(timeout) * time.Second)
		atomic.StoreUint32(&h.sseTimeout, 0)
	}
	return nil
}

// PostMetrics accept metrics data and do validation checks
func (h *Handler) PostMetrics(ctx echo.Context, environment api.EnvironmentPathParam) error {
	metricsData := &api.Metrics{}
	err := ctx.Bind(metricsData)
	if err != nil {
		return errors.New("metrics not present")
	}

	if !h.targetDataReceived && (metricsData.TargetData == nil || len(*metricsData.TargetData) == 0) {
		return errors.New("target data cannot be empty")
	}

	h.targetDataReceived = true

	return ctx.NoContent(http.StatusOK)
}

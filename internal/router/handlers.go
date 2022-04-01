package router

import (
	"github.com/drone/ff-mock-server/internal/config"
	"github.com/drone/ff-mock-server/internal/repository"
	"github.com/drone/ff-mock-server/internal/service"
	"github.com/drone/ff-mock-server/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/r3labs/sse/v2"
	"net/http"
	"sync/atomic"
	"time"
)

type EventSource interface {
	CreateStream(id string) *sse.Stream
	StreamExists(id string) bool
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Close()
}

type Handler struct {
	eventSource EventSource
	repo        repository.Repository
	sseSeq      uint32
}

func NewHandler(repo repository.Repository, eventSource EventSource) *Handler {
	return &Handler{
		eventSource: eventSource,
		repo:        repo,
	}
}

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

func (h *Handler) GetFeatureConfig(ctx echo.Context, environmentUUID string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetFlagConfigurations())
}

func (h *Handler) GetFeatureConfigByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	featureConfig, ok := h.repo.GetFlagConfiguration(identifier)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "feature not found",
		})
	}
	return ctx.JSON(http.StatusOK, featureConfig)
}

func (h *Handler) GetAllSegments(ctx echo.Context, environmentUUID string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetTargetGroups())
}

func (h *Handler) GetSegmentByIdentifier(ctx echo.Context, environmentUUID string, identifier string) error {
	segment, ok := h.repo.GetTargetGroup(identifier)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "segment not found",
		})
	}
	return ctx.JSON(http.StatusOK, segment)
}

func (h *Handler) GetEvaluations(ctx echo.Context, environmentUUID string, target string) error {
	return ctx.JSON(http.StatusOK, h.repo.GetEvaluations())
}

func (h *Handler) GetEvaluationByIdentifier(ctx echo.Context, environmentUUID string, target string, feature string) error {
	evaluation, ok := h.repo.GetEvaluation(feature)
	if !ok {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "segment not found",
		})
	}
	return ctx.JSON(http.StatusOK, evaluation)
}

func (h *Handler) Stream(ctx echo.Context, params api.StreamParams) error {
	log.Infof("connecting key %s on stream", params.APIKey)
	req := ctx.Request()
	req.URL.RawQuery = "stream=" + params.APIKey
	if !h.eventSource.StreamExists(params.APIKey) {
		h.eventSource.CreateStream(params.APIKey)
	}
	seq := atomic.LoadUint32(&h.sseSeq)
	timer := time.Tick(time.Duration(config.Options.SSE[seq]) * time.Second)
	go func() {
		<-timer
		h.eventSource.Close()
	}()
	h.eventSource.ServeHTTP(ctx.Response().Writer, req)
	atomic.StoreUint32(&h.sseSeq, 0)
	if int(seq) < len(config.Options.SSE)-1 {
		atomic.StoreUint32(&h.sseSeq, seq+1)
	}
	return nil
}

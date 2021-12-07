package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	oapimdl "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/drone/ff-mock-server/internal"
	"github.com/drone/ff-mock-server/internal/dto"
	"github.com/drone/ff-mock-server/internal/repository"
	"github.com/drone/ff-mock-server/internal/router"
	"github.com/drone/ff-mock-server/pkg/api"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
)

func main() {

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Cache-Control",
			echo.HeaderAuthorization, "api-key", "Pragma"},
	}))

	e.GET("/health", HealthCheck)

	clientSwagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
	}

	jwtConfig := middleware.JWTConfig{
		Skipper: func(e echo.Context) bool {
			jwt, ok := e.Get(internal.JWTKey).(bool)
			if ok && jwt {
				return false
			}
			return true
		},
		Claims:     &dto.JWTCustomClaims{},
		SigningKey: []byte(internal.GetAuthSecret()), // SDK_AUTH_TOKEN change on next deploy
	}

	clientGroup := e.Group("api/1.0")
	clientGroup.Use(oapimdl.OapiRequestValidatorWithOptions(clientSwagger, &oapimdl.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: router.JWTValidation,
		},
	}))
	clientGroup.Use(middleware.JWTWithConfig(jwtConfig))

	clientGroup.Use(router.ValidateEnvironment())

	server := sse.New()
	repo := repository.NewDummyRepository()
	handler := router.NewHandler(repo, server)
	api.RegisterHandlers(clientGroup, handler)

	// Start server
	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// HealthCheck returns the health of the service
func HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "healthy")
}

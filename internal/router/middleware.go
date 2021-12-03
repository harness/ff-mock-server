package router

import (
	"context"
	"fmt"
	"strings"

	oapimdl "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/drone/ff-mock-server/internal"
	"github.com/drone/ff-mock-server/internal/dto"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	IdentityServiceKey = "identityservice"
	APIKey             = "apikey"
)

// ValidateEnvironment determines that the environment UUID is present in request, and that
// it is valid.   The middleware will return unauthorized if we do not have a valid environment.
func ValidateEnvironment() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			environmentUUID := c.Param("environmentUUID")
			log.Debugf("Validating auth token claim - environment ID with request param environmentUUID: %s",
				environmentUUID)
			jwtPresent, ok := c.Get(internal.JWTKey).(bool)
			log.Infof("JwtPresent %v with status %v", jwtPresent, ok)
			if ok && jwtPresent && environmentUUID != "" {
				// validate jwt
				user, ok := c.Get("user").(*jwt.Token)
				if !ok {
					return echo.ErrUnauthorized
				}

				claims, ok := user.Claims.(*dto.JWTCustomClaims)
				if !ok {
					return echo.ErrUnauthorized
				}
				if claims.Environment != environmentUUID {
					log.Errorf("Environment %s mismatch with requested %s", claims.Environment,
						environmentUUID)
					return echo.NewHTTPError(403, fmt.Sprintf("Environment ID %s mismatch with requested %s",
						claims.Environment, environmentUUID))
				}
			}
			return next(c)
		}
	}
}

// JWTValidation validates that values have been provided for the SecurityScheme
// and sets the type of authorization in the context.  For example if we are authorization with
// BearerAuth then either a standard auth token, or an IdentityService token can be provided.
//
// The SecurityScheme is ApiKeyAuth then we check for a api-key and if enable APIKey as true in the context.
// If there is an invalid scheme or the token/key is missing this function will return an error.
func JWTValidation(c context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName == "BearerAuth" {
		if input.RequestValidationInput.Request.Header.Get("Authorization") == "" {
			return echo.ErrUnauthorized
		}

		// work with JWT
		ctx, ok := c.Value(oapimdl.EchoContextKey).(echo.Context)
		if !ok {
			return echo.ErrUnauthorized
		}

		token := input.RequestValidationInput.Request.Header.Get("Authorization")
		// our gateway service can use an IdentityService token in the authorization header, so if it is present
		// we want to use it instead of bearer
		if strings.Contains(token, "IdentityService") {
			ctx.Set(IdentityServiceKey, true)
		} else {
			ctx.Set(internal.JWTKey, true)
		}

		return nil

	} else if input.SecuritySchemeName == "ApiKeyAuth" {
		if input.RequestValidationInput.Request.Header.Get("api-key") == "" {
			return echo.ErrUnauthorized
		}

		// work with X-API-Key
		ctx, ok := c.Value(oapimdl.EchoContextKey).(echo.Context)
		if !ok {
			return echo.ErrUnauthorized
		}
		ctx.Set(APIKey, true)
		return nil
	}
	return echo.ErrUnauthorized
}

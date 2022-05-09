package config

import (
	"os"

	"github.com/drone/ff-mock-server/internal"
)

// GetAuthSecret get secret for JWT token from environment variable
func GetAuthSecret() string {
	authJwtSecret := os.Getenv("AUTH_SECRET")
	if len(authJwtSecret) == 0 {
		authJwtSecret = internal.DefaultAuthSecret
	}
	return authJwtSecret
}

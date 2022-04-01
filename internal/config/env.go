package config

import (
	"os"

	"github.com/drone/ff-mock-server/internal"
)

func GetAuthSecret() string {
	authJwtSecret := os.Getenv("AUTH_SECRET")
	if len(authJwtSecret) == 0 {
		authJwtSecret = internal.DefaultAuthSecret
	}
	return authJwtSecret
}

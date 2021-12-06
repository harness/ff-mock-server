package internal

import "os"

// GetAuthSecret get secret for JWT token from environment variable
func GetAuthSecret() string {
	authJwtSecret := os.Getenv("AUTH_SECRET")
	if len(authJwtSecret) == 0 {
		authJwtSecret = DefaultAuthSecret
	}
	return authJwtSecret
}

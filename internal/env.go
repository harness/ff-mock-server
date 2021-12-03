package internal

import "os"

func GetAuthSecret() string {
	authJwtSecret := os.Getenv("AUTH_SECRET")
	if len(authJwtSecret) == 0 {
		authJwtSecret = DefaultAuthSecret
	}
	return authJwtSecret
}

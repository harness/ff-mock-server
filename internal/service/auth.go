package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/drone/ff-mock-server/internal"
	"github.com/drone/ff-mock-server/internal/dto"
	"github.com/golang-jwt/jwt"
)

const (
	// ServerKeyType ...
	ServerKeyType = "Server"
	// ClientKeyType ...
	ClientKeyType = "Client"
)

var (
	apiKeyTypes = map[string]string{
		internal.ServerKey: ServerKeyType,
		internal.ClientKey: ClientKeyType,
	}
)

// Authenticate with apiKey and return JWT signed token
func Authenticate(apiKey string) (string, error) {
	apiKeyType, ok := apiKeyTypes[apiKey]
	if !ok {
		return "", fmt.Errorf("api key '%s' not found", apiKey)
	}

	clusterIdentifier := os.Getenv("CLUSTER_IDENTIFIER")
	if len(clusterIdentifier) == 0 {
		clusterIdentifier = internal.DefaultClusterIdentifier
	}
	var jwtKey = []byte(internal.GetAuthSecret())
	claims := &dto.JWTCustomClaims{
		ClusterIdentifier:      clusterIdentifier,
		Account:                "Harness account",
		Organization:           "Harness",
		OrganizationIdentifier: "harness",
		Project:                internal.Project,
		ProjectIdentifier:      internal.Project,
		Environment:            internal.EnvironmentUUID,
		EnvironmentIdentifier:  internal.Environment,
		KeyType:                apiKeyType,
		StandardClaims:         jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

// CheckAPIKeyType ...
func CheckAPIKeyType(expectedTokenType string, token *jwt.Token) error {
	claims, ok := token.Claims.(*dto.JWTCustomClaims)
	if !ok {
		return errors.New("can't decode api key type from request")
	}

	if claims.KeyType != string(expectedTokenType) {
		return errors.New("you cannot access resource with this api key type")
	}
	return nil
}

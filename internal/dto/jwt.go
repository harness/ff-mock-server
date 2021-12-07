package dto

import (
	"github.com/golang-jwt/jwt"
)

// JWTCustomClaims contains all fields for interacting with the FF backend
type JWTCustomClaims struct {
	Environment            string `json:"environment"`
	EnvironmentIdentifier  string `json:"environmentIdentifier"`
	Project                string `json:"project"`
	ProjectIdentifier      string `json:"projectIdentifier"`
	Account                string `json:"accountID"`
	Organization           string `json:"organization"` /*OrgID*/
	OrganizationIdentifier string `json:"organizationIdentifier"`
	ClusterIdentifier      string `json:"clusterIdentifier"`
	KeyType                string `json:"key_type"`
	jwt.StandardClaims
}

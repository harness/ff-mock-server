// Client Service

//go:generate oapi-codegen -generate server,spec -package=api -o services.gen.go ../../api.yaml
//go:generate oapi-codegen -generate types -package=api -o types.gen.go ../../api.yaml

package api

// Client Service

//go:generate oapi-codegen -exclude-tags=metrics --exclude-schemas=KeyValue,Metrics,MetricsData,TargetData -generate server,spec -package=api -o services.gen.go ../../api.yaml
//go:generate oapi-codegen -exclude-tags=metrics --exclude-schemas=KeyValue,Metrics,MetricsData,TargetData -generate types -package=api -o types.gen.go ../../api.yaml

package api

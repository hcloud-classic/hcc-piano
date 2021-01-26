package config

type grpc struct {
	Port             int64 `goconf:"grpc:port"`               // Port : Port number for listening graphql request via http server
	RequestTimeoutMs int64 `goconf:"grpc:request_timeout_ms"` // RequestTimeoutMs : Timeout for HTTP request
}

// HTTP : http config structure
var Grpc grpc

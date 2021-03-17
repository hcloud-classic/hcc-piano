package config

type grpc struct {
	Port             int64 `goconf:"grpc:port"`               // Port : Port number for listening graphql request via http server
	RequestTimeoutMs int64 `goconf:"grpc:request_timeout_ms"` // RequestTimeoutMs : Timeout for HTTP request
}

var Grpc grpc

type harp struct {
	Address          string `goconf:"harp:harp_server_address"`
	Port             int64  `goconf:"harp:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"harp:harp_request_timeout_ms"`
}

var Harp harp

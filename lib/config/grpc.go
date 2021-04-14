package config

type grpc struct {
	Port             int64 `goconf:"grpc:port"`               // Port : Port number for listening graphql request via http server
	RequestTimeoutMs int64 `goconf:"grpc:request_timeout_ms"` // RequestTimeoutMs : Timeout for HTTP request
}

type harp struct {
	Address          string `goconf:"harp:harp_server_address"`
	Port             int64  `goconf:"harp:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"harp:harp_request_timeout_ms"`
}

type cello struct {
	Address          string `goconf:"cello:harp_server_address"`
	Port             int64  `goconf:"cello:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"cello:harp_request_timeout_ms"`
}

type flute struct {
	Address          string `goconf:"flute:harp_server_address"`
	Port             int64  `goconf:"flute:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"flute:harp_request_timeout_ms"`
}

type violin struct {
	Address          string `goconf:"violin:harp_server_address"`
	Port             int64  `goconf:"violin:harp_server_port"`
	RequestTimeoutMs int64  `goconf:"violin:harp_request_timeout_ms"`
}

var Grpc grpc
var Harp harp
var Cello cello
var Flute flute
var Violin violin

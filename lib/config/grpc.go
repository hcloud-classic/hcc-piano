package config

type grpc struct {
	Port string `goconf:"grpc:port"`
}

var Grpc grpc

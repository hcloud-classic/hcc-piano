package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/piano/piano.conf"

type pianoConfig struct {
	GrpcConfig     *goconf.Section
	InfluxdbConfig *goconf.Section
}

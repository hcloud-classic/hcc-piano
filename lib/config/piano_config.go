package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/piano/piano.conf"

type pianoConfig struct {
	InfluxdbConfig *goconf.Section
	GrpcConfig     *goconf.Section
	MysqlConfig    *goconf.Section
	HarpConfig     *goconf.Section
	CelloConfig    *goconf.Section
	FluteConfig    *goconf.Section
	ViolinConfig   *goconf.Section
	PiccoloConfig  *goconf.Section
	BillingConfig  *goconf.Section
}

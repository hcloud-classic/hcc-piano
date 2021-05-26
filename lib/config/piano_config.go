wpackage config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/piano/piano.conf"

type pianoConfig struct {
	InfluxdbConfig *goconf.Section
	GrpcConfig     *goconf.Section
	MysqlConfig    *goconf.Section
	FluteConfig    *goconf.Section
	CelloConfig    *goconf.Section
	HarpConfig     *goconf.Section
	ViolinConfig   *goconf.Section
	PiccoloConfig  *goconf.Section
	BillingConfig  *goconf.Section
}

wpackage config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/piano/piano.conf"

type pianoConfig struct {
	RsakeyConfig   *goconf.Section
	MysqlConfig    *goconf.Section
	InfluxdbConfig *goconf.Section
	GrpcConfig     *goconf.Section
	HornConfig     *goconf.Section
	FluteConfig    *goconf.Section
	CelloConfig    *goconf.Section
	HarpConfig     *goconf.Section
	ViolinConfig   *goconf.Section
	PiccoloConfig  *goconf.Section
	BillingConfig  *goconf.Section
}

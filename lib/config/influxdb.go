package config

type influxdb struct {
	Id       string `goconf:"influxdb:id"`
	Password string `goconf:"influxdb:password"`
	Address  string `goconf:"influxdb:address"`
	Port     int64  `goconf:"influxdb:port"` // 8086
	Db       string `goconf:"influxdb:db"`   // telegraf
}

var Influxdb influxdb

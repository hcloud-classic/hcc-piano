package config

type influxdb struct {
	Id       string `goconf:"influxdb:id"`
	Password string `goconf:"influxdb:password"`
	Host     string `goconf:"influxdb:host"`
	Port     string `goconf:"influxdb:port"` // 8086
	Db       string `goconf:"influxdb:db"`   // telegraf
}

var Influxdb influxdb

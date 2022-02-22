package config

type influxdb struct {
	ID      string `goconf:"influxdb:id"`
	Address string `goconf:"influxdb:address"`
	Port    int64  `goconf:"influxdb:port"` // 8086
	Db      string `goconf:"influxdb:db"`   // telegraf
}

// Influxdb : InfluxDB config structure
var Influxdb influxdb

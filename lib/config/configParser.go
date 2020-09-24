package config

import (
	"github.com/Terry-Mao/goconf"
	"hcc/piano/lib/logger"
)

var conf = goconf.New()
var config = pianoConfig{}
var err error

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		logger.Logger.Panicln("no grpc section")
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseInfluxdb() {

	config.InfluxdbConfig = conf.Get("influxdb")
	if config.InfluxdbConfig == nil {
		logger.Logger.Panicln("no influxdb section")
	}

	Influxdb = influxdb{}

	Influxdb.Id, err = config.InfluxdbConfig.String("id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Influxdb.Password, err = config.InfluxdbConfig.String("password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Influxdb.Address, err = config.InfluxdbConfig.String("address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Influxdb.Port, err = config.InfluxdbConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Influxdb.Db, err = config.InfluxdbConfig.String("db")
	if err != nil {
		logger.Logger.Panicln(err)
	}

}

// Init
func Init() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseGrpc()
	parseInfluxdb()
}

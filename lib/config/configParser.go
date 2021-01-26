package config

import (
	"github.com/Terry-Mao/goconf"
	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var conf = goconf.New()
var config = pianoConfig{}
var err error

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		errors.NewHccError(errors.PianoInternalParsingError, "grpc config").Fatal()
	}

	Grpc = grpc{}
	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "grpc port").Fatal()
	}
}

func parseInfluxdb() {

	config.InfluxdbConfig = conf.Get("influxdb")
	if config.InfluxdbConfig == nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb config").Fatal()
	}

	Influxdb = influxdb{}

	Influxdb.ID, err = config.InfluxdbConfig.String("id")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb id").Fatal()
	}

	Influxdb.Password, err = config.InfluxdbConfig.String("password")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb password").Fatal()
	}

	Influxdb.Address, err = config.InfluxdbConfig.String("address")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb address").Fatal()
	}

	Influxdb.Port, err = config.InfluxdbConfig.Int("port")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb port").Fatal()
	}

	Influxdb.Db, err = config.InfluxdbConfig.String("database")
	if err != nil {
		errors.NewHccError(errors.PianoInternalParsingError, "influxdb database").Fatal()
	}

}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, err.Error()).Fatal()
	}

	parseGrpc()
	parseInfluxdb()
}

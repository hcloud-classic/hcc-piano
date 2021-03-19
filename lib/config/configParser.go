package config

import (
	"github.com/Terry-Mao/goconf"
	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var conf = goconf.New()
var config = pianoConfig{}
var err error

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp config").Fatal()
	}

	Harp = harp{}
	Harp.Address, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp server address").Fatal()
	}
	Harp.Port, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp server port").Fatal()
	}
	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "harp timeout").Fatal()
	}
}

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

	parseHarp()
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

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql config").Fatal()
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql id").Fatal()
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql password").Fatal()
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql address").Fatal()
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql port").Fatal()
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, "mysql database").Fatal()
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		errors.NewHccError(errors.ViolinNoVNCInternalParsingError, err.Error()).Fatal()
	}

	parseGrpc()
	parseInfluxdb()
	parseMysql()
}

package config

import (
	"github.com/Terry-Mao/goconf"
	"innogrid.com/hcloud-classic/hcc_errors"
)

var conf = goconf.New()
var config = pianoConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no mysql section").Fatal()
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
	Mysql.ConnectionRetryCount, err = config.MysqlConfig.Int("connection_retry_count")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Mysql.ConnectionRetryIntervalMs, err = config.MysqlConfig.Int("connection_retry_interval_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

}

func parseFlute() {
	config.FluteConfig = conf.Get("flute")
	if config.FluteConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no flute section").Fatal()
	}

	Flute = flute{}
	Flute.ServerAddress, err = config.FluteConfig.String("flute_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Flute.ServerPort, err = config.FluteConfig.Int("flute_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Flute.RequestTimeoutMs, err = config.FluteConfig.Int("flute_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parseCello() {
	config.CelloConfig = conf.Get("cello")
	if config.CelloConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no cello section").Fatal()
	}

	Cello = cello{}
	Cello.ServerAddress, err = config.CelloConfig.String("cello_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Cello.ServerPort, err = config.CelloConfig.Int("cello_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Cello.RequestTimeoutMs, err = config.CelloConfig.Int("cello_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parseHarp() {
	config.HarpConfig = conf.Get("harp")
	if config.HarpConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no harp section").Fatal()
	}

	Harp = harp{}
	Harp.ServerAddress, err = config.HarpConfig.String("harp_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Harp.ServerPort, err = config.HarpConfig.Int("harp_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Harp.RequestTimeoutMs, err = config.HarpConfig.Int("harp_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parseViolin() {
	config.ViolinConfig = conf.Get("violin")
	if config.ViolinConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no violin section").Fatal()
	}

	Violin = violin{}
	Violin.ServerAddress, err = config.ViolinConfig.String("violin_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Violin.ServerPort, err = config.ViolinConfig.Int("violin_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Violin.RequestTimeoutMs, err = config.ViolinConfig.Int("violin_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parsePiccolo() {
	config.PiccoloConfig = conf.Get("piccolo")
	if config.PiccoloConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no piccolo section").Fatal()
	}

	Piccolo = piccolo{}
	Piccolo.ServerAddress, err = config.PiccoloConfig.String("piccolo_server_address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Piccolo.ServerPort, err = config.PiccoloConfig.Int("piccolo_server_port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Piccolo.RequestTimeoutMs, err = config.PiccoloConfig.Int("piccolo_request_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "no grpc section").Fatal()
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Grpc.ClientPingIntervalMs, err = config.GrpcConfig.Int("client_ping_interval_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	Grpc.ClientPingTimeoutMs, err = config.GrpcConfig.Int("client_ping_timeout_ms")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}
}

func parseInfluxdb() {
	config.InfluxdbConfig = conf.Get("influxdb")
	if config.InfluxdbConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb config").Fatal()
	}

	Influxdb.ID, err = config.InfluxdbConfig.String("id")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb id").Fatal()
	}

	Influxdb.Password, err = config.InfluxdbConfig.String("password")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb password").Fatal()
	}

	Influxdb.Address, err = config.InfluxdbConfig.String("address")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb address").Fatal()
	}

	Influxdb.Port, err = config.InfluxdbConfig.Int("port")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb port").Fatal()
	}

	Influxdb.Db, err = config.InfluxdbConfig.String("database")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "influxdb database").Fatal()
	}
}

func parseBilling() {
	config.BillingConfig = conf.Get("billing")
	if config.BillingConfig == nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "billing config").Fatal()
	}

	Billing.UpdateInterval, err = config.BillingConfig.Int("billing_update_interval_sec")
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, "billing billing_update_interval_sec").Fatal()
	}
}

// Init : Parse config file and initialize config structure
func Init() {
	if err = conf.Parse(configLocation); err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalParsingError, err.Error()).Fatal()
	}

	parseMysql()
	parseGrpc()
	parseInfluxdb()
	parseFlute()
	parseCello()
	parseHarp()
	parseViolin()
	parsePiccolo()
	parseBilling()
}

package main

import (
	"hcc/piano/action/grpc/server"
	"hcc/piano/lib/config"
	"hcc/piano/lib/influxdb"
	"hcc/piano/lib/logger"
	"hcc/piano/lib/syscheck"
)

func main() {

	if !logger.Prepare() {
		logger.Logger.Fatalf("Failed to prepare logger.")
	}
	defer logger.FpLog.Close()

	err := syscheck.CheckRoot()
	if err != nil {
		logger.Logger.Fatalf("Failed to run piano : %v", err)
	}

	err = influxdb.Prepare()
	if err != nil {
		logger.Logger.Fatalf("Failed to prepare InfluxDB : %v", err)
	}
	logger.Logger.Println("InfluxDB is listening on port " + config.InfluxPort)

	server.InitGRPCServer()
}

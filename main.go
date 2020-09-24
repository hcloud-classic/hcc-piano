package main

import (
	gServer "hcc/piano/action/grpc/server"
	"hcc/piano/driver/influxdb"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"hcc/piano/lib/syscheck"
	"log"
)

func init() {
	err := syscheck.CheckRoot()
	if err != nil {
		log.Fatalf("syscheck.CheckRoot(): %v", err.Error())
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("logger.Init(): %v", err.Error())
	}

	config.Init()

	err = influxdb.Init()
	if err != nil {
		logger.Logger.Fatalf("influxdb.Init(): %v", err.Error())
	}
	logger.Logger.Println("InfluxDB is connected to " + config.Influxdb.Port)

}

func end() {
	logger.End()
	influxdb.End()
	gServer.End()
}

func main() {
	gServer.Init()
}

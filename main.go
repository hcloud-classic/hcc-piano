package main

import (
	"hcc/piano/action/graphql"
	"hcc/piano/lib/config"
	"hcc/piano/lib/influxdb"
	"hcc/piano/lib/logger"
	"hcc/piano/lib/syscheck"
	"net/http"
)

func main() {
	if !syscheck.CheckRoot() {
		return
	}
	if !logger.Prepare() {
		return
	}
	defer logger.FpLog.Close()

	//hostInfo := influxdb.HostInfo{URL:"http://"+config.InfluxAddress+":"+config.InfluxPort, Username:config.InfluxID, Password:config.InfluxPassword}
	//influxInfo := influxdb.InfluxInfo{HostInfo: hostInfo, Database: config.InfluxDatabase}
	//err := influxInfo.InitInfluxDB()

	err := influxdb.Prepare()
	if err != nil {
		return
	}
	logger.Logger.Println("InfluxDB is listening on port " + config.InfluxPort)

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + config.HTTPPort)
	err = http.ListenAndServe(":"+config.HTTPPort, nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}

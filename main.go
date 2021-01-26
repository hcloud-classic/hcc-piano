package main

import (
	"fmt"
	"github.com/hcloud-classic/hcc_errors"
	"hcc/piano/action/grpc/server"
	"hcc/piano/driver/influxdb"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	err := logger.Init()
	if err != nil {
		err.Fatal()
	}

	config.Parser()

	errors := influxdb.Init()
	if errors != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "influxdb.Init(): "+errors.Error()).Fatal()
	}

	logger.Logger.Println("Connected to InfluxDB (" + config.Influxdb.Address + ":" +
		strconv.FormatInt(config.Influxdb.Port, 10) + ")")

}

func end() {
	logger.End()
}

func main() {
	// Catch the exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		end()
		fmt.Println("Exiting piano module...")
		os.Exit(0)
	}()

	server.Init()
}

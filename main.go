package main

import (
	"fmt"
	"hcc/piano/action/grpc/server"
	"hcc/piano/driver/influxdb"
	"hcc/piano/lib/config"
	"hcc/piano/lib/errors"
	"hcc/piano/lib/logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func init() {
	err := logger.Init()
	if err != nil {
		errors.SetErrLogger(logger.Logger)
		errors.NewHccError(errors.PianoInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	errors.SetErrLogger(logger.Logger)

	config.Init()

	err = influxdb.Init()
	if err != nil {
		errors.NewHccError(errors.PianoInternalInitFail, "influxdb.Init(): "+err.Error()).Fatal()
	}

	logger.Logger.Println("InfluxDB is connected to " + config.Influxdb.Address + ":" +
		strconv.FormatInt(config.Influxdb.Port, 10))

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

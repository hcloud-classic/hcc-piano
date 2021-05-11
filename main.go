package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hcc/piano/action/grpc/client"
	"hcc/piano/action/grpc/server"
	"hcc/piano/driver/billing"
	"hcc/piano/driver/influxdb"

	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"hcc/piano/lib/mysql"

	"innogrid.com/hcloud-classic/hcc_errors"
)

func init() {
	err := logger.Init()
	if err != nil {
		hcc_errors.SetErrLogger(logger.Logger)
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "logger.Init(): "+err.Error()).Fatal()
	}
	hcc_errors.SetErrLogger(logger.Logger)

	config.Init()

	err = mysql.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "mysql.Init(): "+err.Error()).Fatal()
	}

	err = influxdb.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "influxdb.Init(): "+err.Error()).Fatal()
	}

	err = client.Init()
	if err != nil {
		hcc_errors.NewHccError(hcc_errors.PianoInternalInitFail, "client.Init(): "+err.Error()).Fatal()
	}

	billing.Init()
}

func end() {
	billing.End()
	logger.End()
	client.End()
	mysql.End()
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

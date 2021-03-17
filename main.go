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
)

func init() {
	err := logger.Init()
	if err != nil {
		err.Fatal()
	}

	config.Parser()

	err = influxdb.Init()
	if err != nil {
		err.Fatal()
	}

	err = mysql.Init()
	if err != nil {
		err.Fatal()
	}

	client.InitGRPCClient()

	err = billing.Init()
	if err != nil {
		err.Fatal()
	}
}

func end() {
	logger.End()
	mysql.End()
	client.CleanGRPCClient()
	server.CleanGRPCServer()
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

	server.InitGRPCServer()
}

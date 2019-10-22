package main

import (
	"hcc/piano/action/rabbitmq"
	"hcc/piano/lib/checkroot"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
)

func main() {
	if !checkroot.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := rabbitmq.PrepareChannel()
	if err != nil {
		logger.Logger.Panic(err)
	}
	defer func() {
		_ = rabbitmq.Channel.Close()
	}()
	defer func() {
		_ = rabbitmq.Connection.Close()
	}()
	err = rabbitmq.XXX()
	if err != nil {
		logger.Logger.Panic(err)
	}

	forever := make(chan bool)

	logger.Logger.Println(" [*] Waiting for messages. To exit press Ctrl+C")
	<-forever
}

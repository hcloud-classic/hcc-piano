package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/driver/influxdb"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"log"
	"net"
)

type server struct {
	rpcpiano.UnimplementedPianoServer
}

var srv *grpc.Server

func (s *server) Telegraph(ctx context.Context, in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {
	series, err := influxdb.GetInfluxData(in)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func Init() error {
	lis, err := net.Listen("tcp", ":"+config.Grpc.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + config.Grpc.Port + "...")

	srv = grpc.NewServer()
	rpcpiano.RegisterPianoServer(srv, &server{})
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func End() {
	srv.Stop()
}

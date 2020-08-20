package grpcsrv

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"hcc/piano/action/grpc/rpcpiano"
	"hcc/piano/driver"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"log"
	"net"
	"strconv"
)

type server struct {
	rpcpiano.UnimplementedPianoServer
}

var srv *grpc.Server

func (s *server) Telegraph(ctx context.Context, in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {
	series, err := driver.GetInfluxData(in)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func InitGRPCServer() error {
	lis, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(config.HttpPort), 10))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening server on port " + strconv.FormatInt(int64(config.HttpPort), 10) + "...")

	srv = grpc.NewServer()
	rpcpiano.RegisterPianoServer(srv, &server{})
	reflection.Register(srv)

	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return err
}

func CleanGRPCServer() {
	srv.Stop()
}

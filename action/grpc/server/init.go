package server

import (
	"net"
	"strconv"

	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"innogrid.com/hcloud-classic/pb"
)

type server struct {
	pb.UnimplementedPianoServer
}

var srv *grpc.Server

func InitGRPCServer() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(config.Grpc.Port)))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	logger.Logger.Println("Opening gRPC server on port " + strconv.Itoa(int(config.Grpc.Port)) + "...")

	srv = grpc.NewServer()
	pb.RegisterPianoServer(srv, &pianoServer{})
	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}

func CleanGRPCServer() {
	srv.Stop()
}

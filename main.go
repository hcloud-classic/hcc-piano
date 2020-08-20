package main

import (
	pb "hcc/piano/action/grpc"
	"hcc/piano/driver/grpcsrv"
	"hcc/piano/lib/config"
	"hcc/piano/lib/influxdb"
	"hcc/piano/lib/logger"
	"hcc/piano/lib/syscheck"
)

type server struct {
	pb.UnimplementedGreeterServer
}

//func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
//	log.Printf("Received: %v", in.GetName())
//	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
//}

func main() {

	if !logger.Prepare() {
		logger.Logger.Fatalf("Failed to prepare logger.")
	}
	defer logger.FpLog.Close()

	err := syscheck.CheckRoot()
	if err != nil {
		logger.Logger.Fatalf("Failed to run piano : %v", err)
	}

	err = influxdb.Prepare()
	if err != nil {
		logger.Logger.Fatalf("Failed to prepare InfluxDB : %v", err)
	}
	logger.Logger.Println("InfluxDB is listening on port " + config.InfluxPort)

	//lis, err := net.Listen("tcp", config.GrpcPort)
	//if err != nil {
	//	log.Fatalf("Failed to listen: %v", err)
	//}
	//s := grpc.NewServer()
	//logger.Logger.Println("GRPC server is listening on port " + config.GrpcPort)
	//
	//pb.RegisterGreeterServer(s, &server{})
	//if err := s.Serve(lis); err != nil {
	//	logger.Logger.Fatalf("failed to serve: %v", err)
	//}

	grpcsrv.InitGRPCServer()
}

package client

import (
	"context"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"net"
	"strconv"
	"time"

	"innogrid.com/hcloud-classic/pb"

	"google.golang.org/grpc"
)

var violinConn *grpc.ClientConn

func initViolin() error {
	var err error

	addr := config.Violin.ServerAddress + ":" + strconv.FormatInt(config.Violin.ServerPort, 10)
	violinConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.violin = pb.NewViolinClient(violinConn)
	logger.Logger.Println("gRPC violin client ready")

	return nil
}

func closeViolin() {
	_ = violinConn.Close()
}

func pingViolin() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Violin.ServerAddress,
		strconv.FormatInt(config.Violin.ServerPort, 10)),
		time.Duration(config.Grpc.ClientPingTimeoutMs)*time.Millisecond)
	if err != nil {
		return false
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
		return true
	}

	return false
}

func checkViolin() {
	ticker := time.NewTicker(time.Duration(config.Grpc.ClientPingIntervalMs) * time.Millisecond)
	go func() {
		connOk := true
		for range ticker.C {
			pingOk := pingViolin()
			if pingOk {
				if !connOk {
					logger.Logger.Println("checkViolin(): Ping Ok! Resetting connection...")
					closeViolin()
					err := initViolin()
					if err != nil {
						logger.Logger.Println("checkViolin(): " + err.Error())
						continue
					}
					connOk = true
				}
			} else {
				if connOk {
					logger.Logger.Println("checkViolin(): violin module seems dead. Pinging...")
				}
				connOk = false
			}
		}
	}()
}

// GetServer : Get infos of the server
func (rc *RPCClient) GetServer(uuid string) (*pb.ResGetServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServer, err := rc.violin.GetServer(ctx, &pb.ReqGetServer{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resGetServer, nil
}

// GetServerList : Get list of the server
func (rc *RPCClient) GetServerList(in *pb.ReqGetServerList) (*pb.ResGetServerList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServerList, err := rc.violin.GetServerList(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetServerList, nil
}

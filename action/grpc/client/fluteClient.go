package client

import (
	"context"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"innogrid.com/hcloud-classic/pb"
)

var fluteConn *grpc.ClientConn

func initFlute() error {
	var err error

	addr := config.Flute.ServerAddress + ":" + strconv.FormatInt(config.Flute.ServerPort, 10)
	fluteConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.flute = pb.NewFluteClient(fluteConn)
	logger.Logger.Println("gRPC flute client ready")

	return nil
}

func closeFlute() {
	_ = fluteConn.Close()
}

func pingFlute() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Flute.ServerAddress,
		strconv.FormatInt(config.Flute.ServerPort, 10)),
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

func checkFlute() {
	ticker := time.NewTicker(time.Duration(config.Grpc.ClientPingIntervalMs) * time.Millisecond)
	go func() {
		connOk := true
		for range ticker.C {
			pingOk := pingFlute()
			if pingOk {
				if !connOk {
					logger.Logger.Println("checkFlute(): Ping Ok! Resetting connection...")
					closeFlute()
					err := initFlute()
					if err != nil {
						logger.Logger.Println("checkFlute(): " + err.Error())
						continue
					}
					connOk = true
				}
			} else {
				if connOk {
					logger.Logger.Println("checkFlute(): Flute module seems dead. Pinging...")
				}
				connOk = false
			}
		}
	}()
}

// GetNodeList : Get the list of nodes
func (rc *RPCClient) GetNodeList(in *pb.ReqGetNodeList) (*pb.ResGetNodeList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetNodeList, err := rc.flute.GetNodeList(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetNodeList, nil
}

// GetNodeUptime : Get the uptime of the node
func (rc *RPCClient) GetNodeUptime(in *pb.ReqGetNodeUptime) (*pb.ResGetNodeUptime, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Flute.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetNodeUptime, err := rc.flute.GetNodeUptime(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetNodeUptime, nil
}

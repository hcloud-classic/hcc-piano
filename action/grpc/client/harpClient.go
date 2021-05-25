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

var harpConn *grpc.ClientConn

func initHarp() error {
	var err error

	addr := config.Harp.ServerAddress + ":" + strconv.FormatInt(config.Harp.ServerPort, 10)
	harpConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.harp = pb.NewHarpClient(harpConn)
	logger.Logger.Println("gRPC harp client ready")

	return nil
}

func closeHarp() {
	_ = harpConn.Close()
}

func pingHarp() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Harp.ServerAddress,
		strconv.FormatInt(config.Harp.ServerPort, 10)),
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

func checkHarp() {
	ticker := time.NewTicker(time.Duration(config.Grpc.ClientPingIntervalMs) * time.Millisecond)
	go func() {
		connOk := true
		for range ticker.C {
			pingOk := pingHarp()
			if pingOk {
				if !connOk {
					logger.Logger.Println("checkHarp(): Ping Ok! Resetting connection...")
					closeHarp()
					err := initHarp()
					if err != nil {
						logger.Logger.Println("checkHarp(): " + err.Error())
						continue
					}
					connOk = true
				}
			} else {
				if connOk {
					logger.Logger.Println("checkHarp(): Harp module seems dead. Pinging...")
				}
				connOk = false
			}
		}
	}()
}

// GetSubnetList : Get the list of subnets
func (rc *RPCClient) GetSubnetList(in *pb.ReqGetSubnetList) (*pb.ResGetSubnetList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	subnetList, err := rc.harp.GetSubnetList(ctx, in)
	if err != nil {
		return nil, err
	}

	return subnetList, nil
}

// GetAdaptiveIPServerList : Get list of the adaptiveIP server
func (rc *RPCClient) GetAdaptiveIPServerList(in *pb.ReqGetAdaptiveIPServerList) (*pb.ResGetAdaptiveIPServerList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	adaptiveIPServerList, err := rc.harp.GetAdaptiveIPServerList(ctx, in)
	if err != nil {
		return nil, err
	}

	return adaptiveIPServerList, nil
}

// GetTraffic : Get the traffic of the server
func (rc *RPCClient) GetTraffic(serverUUID string, day string) (*pb.ResGetTraffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetTraffic, err := rc.harp.GetTraffic(ctx, &pb.ReqGetTraffic{
		ServerUUID: serverUUID,
		Day:        day,
	})
	if err != nil {
		return nil, err
	}

	return resGetTraffic, nil
}

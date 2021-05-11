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

var piccoloConn *grpc.ClientConn

func initPiccolo() error {
	var err error

	addr := config.Piccolo.ServerAddress + ":" + strconv.FormatInt(config.Piccolo.ServerPort, 10)
	piccoloConn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	RC.piccolo = pb.NewPiccoloClient(piccoloConn)
	logger.Logger.Println("gRPC piccolo client ready")

	return nil
}

func closePiccolo() {
	_ = piccoloConn.Close()
}

func pingPiccolo() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(config.Piccolo.ServerAddress,
		strconv.FormatInt(config.Piccolo.ServerPort, 10)),
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

func checkPiccolo() {
	ticker := time.NewTicker(time.Duration(config.Grpc.ClientPingIntervalMs) * time.Millisecond)
	go func() {
		connOk := true
		for range ticker.C {
			pingOk := pingPiccolo()
			if pingOk {
				if !connOk {
					logger.Logger.Println("checkPiccolo(): Ping Ok! Resetting connection...")
					closePiccolo()
					err := initPiccolo()
					if err != nil {
						logger.Logger.Println("checkPiccolo(): " + err.Error())
						continue
					}
					connOk = true
				}
			} else {
				if connOk {
					logger.Logger.Println("checkPiccolo(): Piccolo module seems dead. Pinging...")
				}
				connOk = false
			}
		}
	}()
}

// GetGroupList : Get the group list
func (rc *RPCClient) GetGroupList() (*pb.ResGetGroupList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Piccolo.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetGroupList, err := rc.piccolo.GetGroupList(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return resGetGroupList, nil
}

// GetCharge : Get the charge info of the group
func (rc *RPCClient) GetCharge(groupID int64) (*pb.ResGetCharge, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Piccolo.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteNode, err := rc.piccolo.GetCharge(ctx, &pb.ReqGetCharge{GroupID: groupID})
	if err != nil {
		return nil, err
	}

	return resDeleteNode, nil
}

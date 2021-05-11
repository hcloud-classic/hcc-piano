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

// CreateServer : Create a server
func (rc *RPCClient) CreateServer(in *pb.ReqCreateServer) (*pb.ResCreateServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateServer, err := rc.violin.CreateServer(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreateServer, nil
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

// GetServerNum : Get the number of servers
func (rc *RPCClient) GetServerNum(in *pb.ReqGetServerNum) (*pb.ResGetServerNum, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServerNum, err := rc.violin.GetServerNum(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetServerNum, nil
}

// UpdateServer : Update infos of the server
func (rc *RPCClient) UpdateServer(in *pb.ReqUpdateServer) (*pb.ResUpdateServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resUpdateServer, err := rc.violin.UpdateServer(ctx, in)
	if err != nil {
		return nil, err
	}

	return resUpdateServer, nil
}

// DeleteServer : Delete of the server
func (rc *RPCClient) DeleteServer(uuid string) (*pb.ResDeleteServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteServer, err := rc.violin.DeleteServer(ctx, &pb.ReqDeleteServer{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resDeleteServer, nil
}

// CreateServerNode : Create a server node
func (rc *RPCClient) CreateServerNode(in *pb.ReqCreateServerNode) (*pb.ResCreateServerNode, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateServerNode, err := rc.violin.CreateServerNode(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreateServerNode, nil
}

// GetServerNode : Get infos of the server
func (rc *RPCClient) GetServerNode(uuid string) (*pb.ResGetServerNode, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServerNode, err := rc.violin.GetServerNode(ctx, &pb.ReqGetServerNode{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resGetServerNode, nil
}

// GetServerNodeList : Get list of the server
func (rc *RPCClient) GetServerNodeList(in *pb.ReqGetServerNodeList) (*pb.ResGetServerNodeList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	serverList, err := rc.violin.GetServerNodeList(ctx, in)
	if err != nil {
		return nil, err
	}

	return serverList, nil
}

// GetServerNodeNum : Get the number of servers
func (rc *RPCClient) GetServerNodeNum(serverUUID string) (*pb.ResGetServerNodeNum, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetServerNodeNum, err := rc.violin.GetServerNodeNum(ctx, &pb.ReqGetServerNodeNum{ServerUUID: serverUUID})
	if err != nil {
		return nil, err
	}

	return resGetServerNodeNum, nil
}

// DeleteServerNode : Delete of the serverNode
func (rc *RPCClient) DeleteServerNode(uuid string) (*pb.ResDeleteServerNode, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteServerNode, err := rc.violin.DeleteServerNode(ctx, &pb.ReqDeleteServerNode{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resDeleteServerNode, nil
}

// DeleteServerNodeByServerUUID : Delete of the server
func (rc *RPCClient) DeleteServerNodeByServerUUID(serverUUID string) (*pb.ResDeleteServerNodeByServerUUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Violin.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteServerNodeByServerUUID, err := rc.violin.DeleteServerNodeByServerUUID(ctx, &pb.ReqDeleteServerNodeByServerUUID{ServerUUID: serverUUID})
	if err != nil {
		return nil, err
	}

	return resDeleteServerNodeByServerUUID, nil
}

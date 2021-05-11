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

// CreateSubnet : Create a subnet
func (rc *RPCClient) CreateSubnet(in *pb.ReqCreateSubnet) (*pb.ResCreateSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateSubnet, err := rc.harp.CreateSubnet(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreateSubnet, nil
}

// GetSubnet : Get infos of the subnet
func (rc *RPCClient) GetSubnet(uuid string) (*pb.ResGetSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetSubnet, err := rc.harp.GetSubnet(ctx, &pb.ReqGetSubnet{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resGetSubnet, nil
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

// GetAvailableSubnetList : Get the list of available subnets
func (rc *RPCClient) GetAvailableSubnetList(in *pb.ReqGetAvailableSubnetList) (*pb.ResGetAvailableSubnetList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	subnetList, err := rc.harp.GetAvailableSubnetList(ctx, in)
	if err != nil {
		return nil, err
	}

	return subnetList, nil
}

// GetSubnetNum : Get the number of subnets
func (rc *RPCClient) GetSubnetNum(in *pb.ReqGetSubnetNum) (*pb.ResGetSubnetNum, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetSubnetNum, err := rc.harp.GetSubnetNum(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetSubnetNum, nil
}

// ValidCheckSubnet : Check if we can create the subnet
func (rc *RPCClient) ValidCheckSubnet(in *pb.ReqValidCheckSubnet) (*pb.ResValidCheckSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resValidCheckSubnet, err := rc.harp.ValidCheckSubnet(ctx, in)
	if err != nil {
		return nil, err
	}

	return resValidCheckSubnet, nil
}

// UpdateSubnet : Update infos of the subnet
func (rc *RPCClient) UpdateSubnet(in *pb.ReqUpdateSubnet) (*pb.ResUpdateSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resUpdateSubnet, err := rc.harp.UpdateSubnet(ctx, in)
	if err != nil {
		return nil, err
	}

	return resUpdateSubnet, nil
}

// DeleteSubnet : Delete of the subnet
func (rc *RPCClient) DeleteSubnet(uuid string) (*pb.ResDeleteSubnet, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteSubnet, err := rc.harp.DeleteSubnet(ctx, &pb.ReqDeleteSubnet{UUID: uuid})
	if err != nil {
		return nil, err
	}

	return resDeleteSubnet, nil
}

// CreateAdaptiveIPSetting : Create settings of AdaptiveIP
func (rc *RPCClient) CreateAdaptiveIPSetting(in *pb.ReqCreateAdaptiveIPSetting) (*pb.ResCreateAdaptiveIPSetting, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateAdaptiveIPSetting, err := rc.harp.CreateAdaptiveIPSetting(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreateAdaptiveIPSetting, nil
}

// GetAdaptiveIPAvailableIPList : Get available IP list of AdaptiveIP
func (rc *RPCClient) GetAdaptiveIPAvailableIPList() (*pb.ResGetAdaptiveIPAvailableIPList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetAdaptiveIPAvailableIPList, err := rc.harp.GetAdaptiveIPAvailableIPList(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return resGetAdaptiveIPAvailableIPList, nil
}

// GetAdaptiveIPSetting : Get settings of AdaptiveIP
func (rc *RPCClient) GetAdaptiveIPSetting() (*pb.ResGetAdaptiveIPSetting, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resAdaptiveIPSetting, err := rc.harp.GetAdaptiveIPSetting(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	return resAdaptiveIPSetting, nil
}

// CreateAdaptiveIPServer : Create AdaptiveIP server
func (rc *RPCClient) CreateAdaptiveIPServer(in *pb.ReqCreateAdaptiveIPServer) (*pb.ResCreateAdaptiveIPServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateAdaptiveIPServer, err := rc.harp.CreateAdaptiveIPServer(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreateAdaptiveIPServer, nil
}

// GetAdaptiveIPServer : Get infos of the adaptiveIP server
func (rc *RPCClient) GetAdaptiveIPServer(serverUUID string) (*pb.ResGetAdaptiveIPServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetAdaptiveIPServer, err := rc.harp.GetAdaptiveIPServer(ctx, &pb.ReqGetAdaptiveIPServer{ServerUUID: serverUUID})
	if err != nil {
		return nil, err
	}

	return resGetAdaptiveIPServer, nil
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

// GetAdaptiveIPServerNum : Get the number of adaptiveIP server
func (rc *RPCClient) GetAdaptiveIPServerNum(in *pb.ReqGetAdaptiveIPServerNum) (*pb.ResGetAdaptiveIPServerNum, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resGetAdaptiveIPServerNum, err := rc.harp.GetAdaptiveIPServerNum(ctx, in)
	if err != nil {
		return nil, err
	}

	return resGetAdaptiveIPServerNum, nil
}

// DeleteAdaptiveIPServer : Delete of the adaptiveIP server
func (rc *RPCClient) DeleteAdaptiveIPServer(serverUUID string) (*pb.ResDeleteAdaptiveIPServer, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeleteAdaptiveIPServer, err := rc.harp.DeleteAdaptiveIPServer(ctx, &pb.ReqDeleteAdaptiveIPServer{ServerUUID: serverUUID})
	if err != nil {
		return nil, err
	}

	return resDeleteAdaptiveIPServer, nil
}

// GetPortForwardingList : Get list of the AdaptiveIP Port Forwarding
func (rc *RPCClient) GetPortForwardingList(in *pb.ReqGetPortForwardingList) (*pb.ResGetPortForwardingList, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	adaptiveIPServerList, err := rc.harp.GetPortForwardingList(ctx, in)
	if err != nil {
		return nil, err
	}

	return adaptiveIPServerList, nil
}

// CreatePortForwarding : Create the AdaptiveIP Port Forwarding
func (rc *RPCClient) CreatePortForwarding(in *pb.ReqCreatePortForwarding) (*pb.ResCreatePortForwarding, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreatePortForwarding, err := rc.harp.CreatePortForwarding(ctx, in)
	if err != nil {
		return nil, err
	}

	return resCreatePortForwarding, nil
}

// DeletePortForwarding : Delete the AdaptiveIP Port Forwarding
func (rc *RPCClient) DeletePortForwarding(in *pb.ReqDeletePortForwarding) (*pb.ResDeletePortForwarding, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resDeletePortForwarding, err := rc.harp.DeletePortForwarding(ctx, in)
	if err != nil {
		return nil, err
	}

	return resDeletePortForwarding, nil
}

// CreateDHCPDConfig : Do dhcpd config file creation works
func (rc *RPCClient) CreateDHCPDConfig(subnetUUID string) (*pb.ResCreateDHCPDConf, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.Harp.RequestTimeoutMs)*time.Millisecond)
	defer cancel()
	resCreateDHCPDConf, err := rc.harp.CreateDHCPDConf(ctx, &pb.ReqCreateDHCPDConf{
		SubnetUUID: subnetUUID,
	})
	if err != nil {
		return nil, err
	}

	return resCreateDHCPDConf, nil
}

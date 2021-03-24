package client

import (
	"sync"

	"innogrid.com/hcloud-classic/pb"
)

type RpcClient struct {
	harp   pb.HarpClient
	cello  pb.CelloClient
	flute  pb.FluteClient
	violin pb.ViolinClient
}

var RC = &RpcClient{}

func InitGRPCClient() {
	var wg sync.WaitGroup

	wg.Add(3)
	go initHarp(&wg)
	//go initCello(&wg)
	go initFlute(&wg)
	go initViolin(&wg)

	wg.Wait()
}

func CleanGRPCClient() {
	cleanHarp()
	//cleanCello()
	cleanFlute()
	cleanViolin()
}

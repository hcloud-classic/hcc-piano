package client

import (
	"innogrid.com/hcloud-classic/pb"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	flute   pb.FluteClient
	harp    pb.HarpClient
	violin  pb.ViolinClient
	cello   pb.CelloClient
	piccolo pb.PiccoloClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {
	err := initFlute()
	if err != nil {
		return err
	}
	checkFlute()

	err = initHarp()
	if err != nil {
		return err
	}
	checkHarp()

	err = initViolin()
	if err != nil {
		return err
	}
	checkViolin()

	err = initCello()
	if err != nil {
		return err
	}
	checkCello()

	err = initPiccolo()
	if err != nil {
		return err
	}
	checkPiccolo()

	return nil
}

// End : Close connections of gRPC clients
func End() {
	closePiccolo()
	closeCello()
	closeViolin()
	closeHarp()
	closeFlute()
}

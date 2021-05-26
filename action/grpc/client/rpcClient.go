package client

import (
	"innogrid.com/hcloud-classic/pb"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	flute   pb.FluteClient
	harp    pb.HarpClient
	cello   pb.CelloClient
	violin  pb.ViolinClient
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

	err = initCello()
	if err != nil {
		return err
	}
	checkCello()

	err = initViolin()
	if err != nil {
		return err
	}
	checkViolin()

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
	closeViolin()
	closeCello()
	closeHarp()
	closeFlute()
}

package client

import (
	//"context"
	"strconv"
	//"strings"
	"sync"
	//"time"

	"google.golang.org/grpc"
	errors "innogrid.com/hcloud-classic/hcc_errors"
	//"innogrid.com/hcloud-classic/hcc_errors/errconv"
	rpcharp "innogrid.com/hcloud-classic/pb"

	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
)

var harpconn *grpc.ClientConn

func initHarp(wg *sync.WaitGroup) *errors.HccError {
	var err error
	addr := config.Harp.Address + ":" + strconv.FormatInt(config.Harp.Port, 10)
	logger.Logger.Println("Try connect to harp " + addr)
	harpconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.NewHccError(errors.PianoGrpcConnectionFail, "harp : "+err.Error())
	}

	RC.harp = rpcharp.NewHarpClient(harpconn)
	logger.Logger.Println("GRPC connection to harp created")

	wg.Done()
	return nil
}

func cleanHarp() {
	harpconn.Close()
}

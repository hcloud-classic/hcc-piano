package client

import (
	//"context"
	"strconv"
	//"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	errors "innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"

	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
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

	RC.harp = pb.NewHarpClient(harpconn)
	logger.Logger.Println("GRPC connection to harp created")

	wg.Done()
	return nil
}

func cleanHarp() {
	logger.Logger.Println("Close Connection to Harp")
	harpconn.Close()
}

func (rc *RpcClient) GetNetworkBillingInfo(groupList *[]int32) (*[]model.NetworkBill, *errors.HccErrorStack) {
	var billList []model.NetworkBill

	now := time.Now()

	for _, gid := range *groupList {
		billList = append(billList, model.NetworkBill{
			GroupID:            int(gid),
			Date:               strconv.Itoa(now.Year()%100*10000 + int(now.Month())*100 + now.Day()),
			SubnetCount:        now.Day(),
			AIPCount:           int(now.Month()),
			SubnetChargePerCnt: 5000,
			AIPChargePerCnt:    10000,
			DiscountRate:       0, // Not working now
		})
	}

	return &billList, nil
}

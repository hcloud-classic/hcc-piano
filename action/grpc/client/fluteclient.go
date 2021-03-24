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

var fluteconn *grpc.ClientConn

func initFlute(wg *sync.WaitGroup) *errors.HccError {
	var err error
	addr := config.Flute.Address + ":" + strconv.FormatInt(config.Flute.Port, 10)
	logger.Logger.Println("Try connect to flute " + addr)
	fluteconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.NewHccError(errors.PianoGrpcConnectionFail, "flute : "+err.Error())
	}

	RC.flute = pb.NewFluteClient(fluteconn)
	logger.Logger.Println("GRPC connection to flute created")

	wg.Done()
	return nil
}

func cleanFlute() {
	logger.Logger.Println("Close Connection to Flute")
	fluteconn.Close()
}

func (rc *RpcClient) GetNodeBillingInfo(groupList *[]int32) (*[]model.NodeBill, *errors.HccErrorStack) {
	var billList []model.NodeBill

	now := time.Now()

	for _, gid := range *groupList {
		billList = append(billList, model.NodeBill{
			GroupID:      int(gid),
			Date:         strconv.Itoa(now.Year()%100*10000 + int(now.Month())*100 + now.Day()),
			NodeUUID:     `"Node-UUID-Group-` + strconv.Itoa(int(gid)) + `"`,
			DefChargeCPU: 15000,
			DefChargeMEM: 15000,
			DefChargeNIC: 15000,
			DiscountRate: 0, // Not working now
		})
	}

	return &billList, nil
}

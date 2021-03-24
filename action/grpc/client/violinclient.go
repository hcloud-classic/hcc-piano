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

var violinconn *grpc.ClientConn

func initViolin(wg *sync.WaitGroup) *errors.HccError {
	var err error
	addr := config.Violin.Address + ":" + strconv.FormatInt(config.Violin.Port, 10)
	logger.Logger.Println("Try connect to violin " + addr)
	violinconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.NewHccError(errors.PianoGrpcConnectionFail, "violin : "+err.Error())
	}

	RC.violin = pb.NewViolinClient(violinconn)
	logger.Logger.Println("GRPC connection to violin created")

	wg.Done()
	return nil
}

func cleanViolin() {
	logger.Logger.Println("Close Connection to Violin")
	violinconn.Close()
}

func (rc *RpcClient) GetServerBillingInfo(groupList *[]int32) (*[]model.ServerBill, *errors.HccErrorStack) {
	var billList []model.ServerBill

	now := time.Now()

	for _, gid := range *groupList {
		billList = append(billList, model.ServerBill{
			GroupID:            int(gid),
			Date:               strconv.Itoa(now.Year()%100*10000 + int(now.Month())*100 + now.Day()),
			ServerUUID:         `"Server-UUID-Group-` + strconv.Itoa(int(gid)) + `"`,
			NetworkTraffic:     1500000,
			TrafficChargePerKB: 1.5,
			DiscountRate:       0, // Not working now
		})
	}

	return &billList, nil
}

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

var celloconn *grpc.ClientConn

func initCello(wg *sync.WaitGroup) *errors.HccError {
	var err error
	addr := config.Cello.Address + ":" + strconv.FormatInt(config.Cello.Port, 10)
	logger.Logger.Println("Try connect to cello " + addr)
	celloconn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.NewHccError(errors.PianoGrpcConnectionFail, "cello : "+err.Error())
	}

	RC.cello = pb.NewCelloClient(celloconn)
	logger.Logger.Println("GRPC connection to cello created")

	wg.Done()
	return nil
}

func cleanCello() {
	logger.Logger.Println("Close Connection to Cello")
	celloconn.Close()
}

func (rc *RpcClient) GetVolumeBillingInfo(groupList *[]int32) (*[]model.VolumeBill, *errors.HccErrorStack) {
	var billList []model.VolumeBill

	now := time.Now()

	for _, gid := range *groupList {
		billList = append(billList, model.VolumeBill{
			GroupID:         int(gid),
			Date:            strconv.Itoa(now.Year()%100*10000 + int(now.Month())*100 + now.Day()),
			HDDSize:         16000,
			SSDSize:         1000,
			NVMESize:        250,
			HDDChargePerGB:  150,
			SSDChargePerGB:  1500,
			NVMEChargePerGB: 15000,
			DiscountRate:    0, // Not working now
		})
	}

	return &billList, nil
}

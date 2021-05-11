package billing

import (
	"time"

	"innogrid.com/hcloud-classic/hcc_errors"

	"hcc/piano/action/grpc/client"
	"hcc/piano/dao"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
)

type Billing struct {
	lastUpdate  time.Time
	updateTimer *time.Ticker
	StopTimer   func()
}

func (bill *Billing) RunUpdateTimer() {
	if bill.updateTimer == nil {
		bill.updateTimer = time.NewTicker(1 * time.Hour)
	} else {
		// upper go v1.15
		// bill.updateTimer.Reset(duration)
		bill.updateTimer.Stop()
		bill.updateTimer = time.NewTicker(1 * time.Hour)

		return
	}
	done := make(chan bool)
	bill.StopTimer = func() {
		done <- true
		bill.updateTimer.Stop()
	}

	go func() {
		defer func() {
			bill.updateTimer.Stop()
			bill.updateTimer = nil
		}()

		for true {
			select {
			case <-done:
				return
			case <-bill.updateTimer.C:
				logger.Logger.Println("RunUpdateTimer(): Updating billing information")
				DriverBilling.UpdateBillingInfo()
				break
			}
		}
	}()
}

func (bill *Billing) UpdateBillingInfo() *hcc_errors.HccErrorStack {
	var hccErr *hcc_errors.HccError
	errStack := hcc_errors.NewHccErrorStack()

	resGetGroupList, err := client.RC.GetGroupList()
	if err != nil {
		_ = errStack.Push(hcc_errors.NewHccError(hcc_errors.PianoInternalOperationFail,
			"UpdateBillingInfo(): GetGroupList(): "+err.Error()))
		return errStack
	}

	nodeBillList, err := getNodeBillingInfo(resGetGroupList.Group)
	if err != nil {
		_ = errStack.Push(hcc_errors.NewHccError(hcc_errors.PianoInternalOperationFail,
			"UpdateBillingInfo(): getNodeBillingInfo(): "+err.Error()))
	}

	// TODO: Need to implement getServerBillingInfo()
	//serverBillList, err := getServerBillingInfo(resGetGroupList.Group)
	//if err != nil {
	//	_ = errStack.Push(hcc_errors.NewHccError(hcc_errors.PianoInternalOperationFail,
	//		"UpdateBillingInfo(): getServerBillingInfo(): "+err.Error()))
	//}

	// TODO: Need to implement getNetworkBillingInfo()
	//networkBillList, err := getNetworkBillingInfo(resGetGroupList.Group)
	//if err != nil {
	//	_ = errStack.Push(hcc_errors.NewHccError(hcc_errors.PianoInternalOperationFail,
	//		"UpdateBillingInfo(): getNetworkBillingInfo(): "+err.Error()))
	//}

	// TODO: Need to implement getVolumeBillingInfo()
	//volumeBillList, err := getVolumeBillingInfo(resGetGroupList.Group)
	//if err != nil {
	//	_ = errStack.Push(hcc_errors.NewHccError(hcc_errors.PianoInternalOperationFail,
	//		"UpdateBillingInfo(): getVolumeBillingInfo(): "+err.Error()))
	//}

	hccErr = dao.InsertNodeBillingInfo(nodeBillList)
	if hccErr != nil {
		_ = errStack.Push(hccErr)
	}

	// TODO: Need to implement getServerBillingInfo()
	//hccErr = dao.InsertServerBillingInfo(serverBillList)
	//if hccErr != nil {
	//	_ = errStack.Push(hccErr)
	//}

	// TODO: Need to implement getNetworkBillingInfo()
	//hccErr = dao.InsertNetworkBillingInfo(networkBillList)
	//if hccErr != nil {
	//	_ = errStack.Push(hccErr)
	//}

	// TODO: Need to implement getVolumeBillingInfo()
	//hccErr = dao.InsertVolumeBillingInfo(volumeBillList)
	//if hccErr != nil {
	//	_ = errStack.Push(hccErr)
	//}

	if bill.lastUpdate.Day() != time.Now().Day() {
		logger.Logger.Println("Update Daily Billing Info")
		hccErr = dao.InsertDailyInfo()
		if hccErr != nil {
			_ = errStack.Push(hccErr)
		} else {
			bill.lastUpdate = time.Now()
		}
	}

	if errStack.Len() > 0 {
		_ = errStack.Dump()
		return errStack
	}

	return nil
}

func (bill *Billing) readNetworkBillingInfo(groupID int, date, billType string) (*model.NetworkBill, *hcc_errors.HccError) {
	var billInfo model.NetworkBill

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err == nil {
		for res.Next() {
			_ = res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.SubnetCharge,
				&billInfo.AdaptiveIPCharge)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readNodeBillingInfo(groupID int, date, billType string) (*model.NodeBill, *hcc_errors.HccError) {
	var billInfo model.NodeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "node")
	if err == nil {
		for res.Next() {
			_ = res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.NodeUUID,
				&billInfo.ChargeCPU,
				&billInfo.ChargeMEM,
				&billInfo.ChargeNIC)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readServerBillingInfo(groupID int, date, billType string) (*model.ServerBill, *hcc_errors.HccError) {
	var billInfo model.ServerBill

	res, err := dao.GetBillInfo(groupID, date, billType, "server")
	if err == nil {
		for res.Next() {
			_ = res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.ServerUUID,
				&billInfo.NetworkTraffic,
				&billInfo.TrafficChargePerKB)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int, date, billType string) (*model.VolumeBill, *hcc_errors.HccError) {
	var billInfo model.VolumeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err == nil {
		for res.Next() {
			_ = res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.HDDCharge,
				&billInfo.SSDCharge)
		}
	}
	return &billInfo, err
}

func (bill *Billing) ReadBillingData(groupID *[]int32, dateStart, dateEnd, billType string, row, page int) (*[][]model.Bill, *hcc_errors.HccErrorStack) {
	var billList [][]model.Bill
	errStack := hcc_errors.NewHccErrorStack()

	for _, gid := range *groupID {
		res, err := dao.GetBill(int(gid), dateStart, dateEnd, billType, row, page)
		if err != nil {
			_ = errStack.Push(err)
			continue
		}
		var list []model.Bill
		for res.Next() {
			bill := model.Bill{}
			_ = res.Scan(&bill.BillID,
				&bill.ChargeNode,
				&bill.ChargeServer,
				&bill.ChargeNetwork,
				&bill.ChargeVolume)
			list = append(list, bill)
		}
		billList = append(billList, list)
		_ = res.Close()
	}

	return &billList, errStack
}

func (bill *Billing) ReadBillingDetail(groupID int32, date, billType string) (*model.BillDetail, *hcc_errors.HccErrorStack) {
	var billingDetail model.BillDetail
	var err *hcc_errors.HccError
	errStack := hcc_errors.NewHccErrorStack()

	billingDetail.DetailNode, err = bill.readNodeBillingInfo(int(groupID), date, billType)
	if err != nil {
		_ = errStack.Push(err)
	}
	billingDetail.DetailServer, err = bill.readServerBillingInfo(int(groupID), date, billType)
	if err != nil {
		_ = errStack.Push(err)
	}
	billingDetail.DetailNetwork, err = bill.readNetworkBillingInfo(int(groupID), date, billType)
	if err != nil {
		_ = errStack.Push(err)
	}
	billingDetail.DetailVolume, err = bill.readVolumeBillingInfo(int(groupID), date, billType)
	if err != nil {
		_ = errStack.Push(err)
	}

	logger.Logger.Println(*billingDetail.DetailNode)
	logger.Logger.Println(*billingDetail.DetailServer)
	logger.Logger.Println(*billingDetail.DetailNetwork)
	logger.Logger.Println(*billingDetail.DetailVolume)
	return &billingDetail, errStack
}

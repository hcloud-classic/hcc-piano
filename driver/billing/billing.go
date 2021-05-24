package billing

import (
	"errors"
	"time"

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

func (bill *Billing) UpdateBillingInfo() {
	resGetGroupList, err := client.RC.GetGroupList()
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): GetGroupList(): " + err.Error())
		return
	}

	nodeBillList, err := getNodeBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getNodeBillingInfo(): " + err.Error())
	} else {
		err = dao.InsertNodeBillingInfo(nodeBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertNodeBillingInfo(): " + err.Error())
		}
	}

	serverBillList, err := getServerBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getServerBillingInfo(): " + err.Error())
	} else {
		err = dao.InsertServerBillingInfo(serverBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertServerBillingInfo(): " + err.Error())
		}
	}

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
		logger.Logger.Println("Updating Daily Billing Info")
		err = dao.InsertDailyInfo()
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertDailyInfo(): " + err.Error())
		} else {
			bill.lastUpdate = time.Now()
		}
	}
}

func (bill *Billing) readNetworkBillingInfo(groupID int, date, billType string) (*[]model.NetworkBill, error) {
	var billList []model.NetworkBill

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.NetworkBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.Date,
			&billInfo.ChargeSubnet,
			&billInfo.ChargeAdaptiveIP)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) readNodeBillingInfo(groupID int, date, billType string) (*[]model.NodeBill, error) {
	var billList []model.NodeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "node")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.NodeBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.NodeUUID,
			&billInfo.ChargeCPU,
			&billInfo.ChargeMEM,
			&billInfo.ChargeNIC)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) readServerBillingInfo(groupID int, date, billType string) (*[]model.ServerBill, error) {
	var billList []model.ServerBill

	res, err := dao.GetBillInfo(groupID, date, billType, "server")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.ServerBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.ServerUUID,
			&billInfo.ChargeTraffic)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int, date, billType string) (*[]model.VolumeBill, error) {
	var billList []model.VolumeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.VolumeBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.Date,
			&billInfo.HDDCharge,
			&billInfo.SSDCharge)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) ReadBillingData(groupID *[]int32, dateStart, dateEnd, billType string, row, page int) (*[][]model.Bill, error) {
	var billList [][]model.Bill

	for _, gid := range *groupID {
		res, err := dao.GetBill(int(gid), dateStart, dateEnd, billType, row, page)
		if err != nil {
			logger.Logger.Println("ReadBillingData(): dao.GetBill(): " + err.Error())
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

	return &billList, nil
}

func (bill *Billing) ReadBillingDetail(groupID int32, date, billType string) (*model.BillDetail, error) {
	var err error
	var returnErr error = nil
	var billingDetail model.BillDetail

	billingDetail.DetailNode, err = bill.readNodeBillingInfo(int(groupID), date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNodeBillingInfo(): " + err.Error())
		returnErr = errors.New(err.Error())
	}
	billingDetail.DetailServer, err = bill.readServerBillingInfo(int(groupID), date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readServerBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}
	billingDetail.DetailNetwork, err = bill.readNetworkBillingInfo(int(groupID), date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNetworkBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}
	billingDetail.DetailVolume, err = bill.readVolumeBillingInfo(int(groupID), date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readVolumeBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	return &billingDetail, returnErr
}

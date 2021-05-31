package billing

import (
	"errors"
	"hcc/piano/lib/config"
	"time"

	"hcc/piano/action/grpc/client"
	"hcc/piano/dao"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
)

type Billing struct {
	updateTimer *time.Ticker
	StopTimer   func()
	IsRunning   bool
}

func (bill *Billing) RunUpdateTimer() {
	if bill.updateTimer == nil {
		bill.updateTimer = time.NewTicker(time.Duration(config.Billing.UpdateInterval) * time.Second)
	} else {
		// upper go v1.15
		// bill.updateTimer.Reset(duration)
		bill.updateTimer.Stop()
		bill.updateTimer = time.NewTicker(time.Duration(config.Billing.UpdateInterval) * time.Second)

		return
	}
	done := make(chan bool)
	bill.StopTimer = func() {
		done <- true
	}

	go func() {
		for true {
			select {
			case <-done:
				logger.Logger.Println("RunUpdateTimer(): Stopping billing update timer")

				bill.updateTimer.Stop()
				bill.updateTimer = nil

				return
			case <-bill.updateTimer.C:
				if config.Billing.Debug == "on" {
					logger.Logger.Println("RunUpdateTimer(): Updating billing information")
				}
				DriverBilling.UpdateBillingInfo()
				break
			}
		}
	}()
}

func (bill *Billing) UpdateBillingInfo() {
	bill.IsRunning = true

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting group list")
	}
	resGetGroupList, err := client.RC.GetGroupList()
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): GetGroupList(): " + err.Error())
		bill.IsRunning = false
		return
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting node_billing_info")
	}
	nodeBillList, err := getNodeBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getNodeBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting node_billing_info")
		}
		err = dao.InsertNodeBillingInfo(nodeBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertNodeBillingInfo(): " + err.Error())
		}
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting server_billing_info")
	}
	serverBillList, err := getServerBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getServerBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting server_billing_info")
		}
		err = dao.InsertServerBillingInfo(serverBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertServerBillingInfo(): " + err.Error())
		}
	}

	networkBillList, err := getNetworkBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getNetworkBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting network_billing_info")
		}
		err = dao.InsertNetworkBillingInfo(networkBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertNetworkBillingInfo(): " + err.Error())
		}
	}

	volumeBillList, err := getVolumeBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getVolumeBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting volume_billing_info")
		}
		err = dao.InsertVolumeBillingInfo(volumeBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertVolumeBillingInfo(): " + err.Error())
		}
	}

	if err == nil {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Getting daily_info")
		}
		dailyBillList := dao.GetDailyInfo(resGetGroupList.Group, nodeBillList, serverBillList, networkBillList, volumeBillList)
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting daily_info")
		}
		err = dao.InsertDailyInfo(dailyBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertDailyInfo(): " + err.Error())
		}
	}

	bill.IsRunning = false
}

func (bill *Billing) readNetworkBillingInfo(groupID int64, date, billType string) (*[]model.NetworkBill, error) {
	var billList []model.NetworkBill

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.NetworkBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.ChargeSubnet,
			&billInfo.ChargeAdaptiveIP)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) readNodeBillingInfo(groupID int64, date, billType string) (*[]model.NodeBill, error) {
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

func (bill *Billing) readServerBillingInfo(groupID int64, date, billType string) (*[]model.ServerBill, error) {
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

func (bill *Billing) readVolumeBillingInfo(groupID int64, date, billType string) (*[]model.VolumeBill, error) {
	var billList []model.VolumeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var billInfo model.VolumeBill
		_ = res.Scan(&billInfo.GroupID,
			&billInfo.ChargeSSD,
			&billInfo.ChargeHDD)
		billList = append(billList, billInfo)
	}

	return &billList, err
}

func (bill *Billing) ReadBillingData(groupID *[]int64, dateStart, dateEnd, billType string, row, page int) (*[]model.Bill, error) {
	var billList []model.Bill
	var groupIDAll []int64

	if len(*groupID) == 0 {
		resGetGroupList, err := client.RC.GetGroupList()
		if err != nil {
			return &billList, err
		}

		for _, group := range resGetGroupList.Group {
			if group.Id == 1 {
				continue
			}
			groupIDAll = append(groupIDAll, group.Id)
		}

		groupID = &groupIDAll
	}

	res, err := dao.GetBill(groupID, dateStart, dateEnd, billType, row, page)
	if err != nil {
		logger.Logger.Println("ReadBillingData(): dao.GetBill(): " + err.Error())
		return &billList, err
	}

	for res.Next() {
		var bill model.Bill
		_ = res.Scan(&bill.Date,
			&bill.GroupID,
			&bill.ChargeNode,
			&bill.ChargeServer,
			&bill.ChargeNetwork,
			&bill.ChargeVolume)
		billList = append(billList, bill)
	}

	return &billList, nil
}

func (bill *Billing) ReadBillingDetail(groupID int64, date, billType string) (*model.BillDetail, error) {
	var err error
	var returnErr error = nil
	var billingDetail model.BillDetail

	billingDetail.DetailNode, err = bill.readNodeBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNodeBillingInfo(): " + err.Error())
		returnErr = errors.New(err.Error())
	}
	billingDetail.DetailServer, err = bill.readServerBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readServerBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}
	billingDetail.DetailNetwork, err = bill.readNetworkBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNetworkBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}
	billingDetail.DetailVolume, err = bill.readVolumeBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readVolumeBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	return &billingDetail, returnErr
}

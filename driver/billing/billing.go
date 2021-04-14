package billing

import (
	"time"

	errors "innogrid.com/hcloud-classic/hcc_errors"

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
				logger.Logger.Println("UPDATE Billing information", []int32{1000, 1001, 1002})
				BillingDriver.UpdateBillingInfo(&[]int32{1000, 1001, 1002})
				break
			}
		}
	}()
}

func (bill *Billing) UpdateBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	var err *errors.HccError
	errStack := errors.NewHccErrorStack()

	nodeBillList, es := client.RC.GetNodeBillingInfo(groupID)
	if es != nil {
		errStack.Merge(es)
	}

	serverBillList, es := client.RC.GetServerBillingInfo(groupID)
	if es != nil {
		errStack.Merge(es)
	}

	networkBillList, es := client.RC.GetNetworkBillingInfo(groupID)
	if es != nil {
		errStack.Merge(es)
	}

	volumeBillList, es := client.RC.GetVolumeBillingInfo(groupID)
	if es != nil {
		errStack.Merge(es)
	}

	err = dao.InsertNodeBillingInfo(nodeBillList)
	if err != nil {
		errStack.Push(err)
	}
	err = dao.InsertServerBillingInfo(serverBillList)
	if err != nil {
		errStack.Push(err)
	}
	err = dao.InsertNetworkBillingInfo(networkBillList)
	if err != nil {
		errStack.Push(err)
	}
	err = dao.InsertVolumeBillingInfo(volumeBillList)
	if err != nil {
		errStack.Push(err)
	}

	if bill.lastUpdate.Day() != time.Now().Day() {
		logger.Logger.Println("Update Daily Billing Info")
		err = dao.InsertDailyInfo()
		if err != nil {
			errStack.Push(err)
		} else {
			bill.lastUpdate = time.Now()
		}
	}

	if errStack.Len() > 0 {
		errStack.Dump()
		return errStack
	}

	return nil
}

func (bill *Billing) readNetworkBillingInfo(groupID int, date, billType string) (*[]model.NetworkBill, *errors.HccError) {
	var billList []model.NetworkBill

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err == nil {
		for res.Next() {
			var billInfo model.NetworkBill
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.SubnetCount,
				&billInfo.AIPCount,
				&billInfo.SubnetChargePerCnt,
				&billInfo.AIPChargePerCnt,
				&billInfo.DiscountRate)
			billList = append(billList, billInfo)
		}
	}
	return &billList, err
}

func (bill *Billing) readNodeBillingInfo(groupID int, date, billType string) (*[]model.NodeBill, *errors.HccError) {
	var billList []model.NodeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "node")
	if err == nil {
		for res.Next() {
			var billInfo model.NodeBill
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.NodeUUID,
				&billInfo.DefChargeCPU,
				&billInfo.DefChargeMEM,
				&billInfo.DefChargeNIC,
				&billInfo.DiscountRate)
			billList = append(billList, billInfo)
		}
	}
	return &billList, err
}

func (bill *Billing) readServerBillingInfo(groupID int, date, billType string) (*[]model.ServerBill, *errors.HccError) {
	var billList []model.ServerBill

	res, err := dao.GetBillInfo(groupID, date, billType, "server")
	if err == nil {
		for res.Next() {
			var billInfo model.ServerBill
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.ServerUUID,
				&billInfo.NetworkTraffic,
				&billInfo.TrafficChargePerKB,
				&billInfo.DiscountRate)
			billList = append(billList, billInfo)
		}
	}
	return &billList, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int, date, billType string) (*[]model.VolumeBill, *errors.HccError) {
	var billList []model.VolumeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err == nil {
		for res.Next() {
			var billInfo model.VolumeBill
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.HDDSize,
				&billInfo.SSDSize,
				&billInfo.NVMESize,
				&billInfo.HDDChargePerGB,
				&billInfo.SSDChargePerGB,
				&billInfo.NVMEChargePerGB,
				&billInfo.DiscountRate)
			billList = append(billList, billInfo)
		}
	}
	return &billList, err
}

func (bill *Billing) ReadBillingData(groupID *[]int32, dateStart, dateEnd, billType string, row, page int) (*[][]model.Bill, *errors.HccErrorStack) {
	var billList [][]model.Bill
	errStack := errors.NewHccErrorStack()

	for _, gid := range *groupID {
		res, err := dao.GetBill(int(gid), dateStart, dateEnd, billType, row, page)
		if err != nil {
			errStack.Push(err)
			continue
		}
		var list []model.Bill
		for res.Next() {
			bill := model.Bill{}
			res.Scan(&bill.BillID,
				&bill.ChargeNode,
				&bill.ChargeServer,
				&bill.ChargeNetwork,
				&bill.ChargeVolume)
			list = append(list, bill)
		}
		billList = append(billList, list)
		res.Close()
	}

	return &billList, errStack
}

func (bill *Billing) ReadBillingDetail(groupID int32, date, billType string) (*model.BillDetail, *errors.HccErrorStack) {
	var billingDetail model.BillDetail
	var err *errors.HccError
	errStack := errors.NewHccErrorStack()

	billingDetail.DetailNode, err = bill.readNodeBillingInfo(int(groupID), date, billType)
	if err != nil {
		errStack.Push(err)
	}
	billingDetail.DetailServer, err = bill.readServerBillingInfo(int(groupID), date, billType)
	if err != nil {
		errStack.Push(err)
	}
	billingDetail.DetailNetwork, err = bill.readNetworkBillingInfo(int(groupID), date, billType)
	if err != nil {
		errStack.Push(err)
	}
	billingDetail.DetailVolume, err = bill.readVolumeBillingInfo(int(groupID), date, billType)
	if err != nil {
		errStack.Push(err)
	}

	return &billingDetail, errStack
}

package billing

import (
	"time"

	"hcc/piano/dao"
	"hcc/piano/lib/logger"
	"hcc/piano/model"

	errors "innogrid.com/hcloud-classic/hcc_errors"
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
				logger.Logger.Println("UPDATE Billing information")
				break
			}
		}
	}()
}

func (bil *Billing) updateNetworkBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) updateNodeBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) updateServerBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) updateVolumeBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateBillingData(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateAllBillingData() *errors.HccErrorStack {
	return nil
}

func (bill *Billing) readNetworkBillingInfo(groupID int, date, billType string) (*model.NetworkBill, *errors.HccError) {
	var billInfo model.NetworkBill

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err == nil {
		for res.Next() {
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.SubnetCount,
				&billInfo.AIPCount,
				&billInfo.SubnetChargePerCnt,
				&billInfo.AIPChargePerCnt,
				&billInfo.DiscountRate)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readNodeBillingInfo(groupID int, date, billType string) (*model.NodeBill, *errors.HccError) {
	var billInfo model.NodeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "node")
	if err == nil {
		for res.Next() {
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.NodeUUID,
				&billInfo.DefChargeCPU,
				&billInfo.DefChargeMEM,
				&billInfo.DefChargeNIC,
				&billInfo.DiscountRate)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readServerBillingInfo(groupID int, date, billType string) (*model.ServerBill, *errors.HccError) {
	var billInfo model.ServerBill

	res, err := dao.GetBillInfo(groupID, date, billType, "server")
	if err == nil {
		for res.Next() {
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.ServerUUID,
				&billInfo.NetworkTraffic,
				&billInfo.TrafficChargePerKB,
				&billInfo.DiscountRate)
		}
	}
	return &billInfo, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int, date, billType string) (*model.VolumeBill, *errors.HccError) {
	var billInfo model.VolumeBill

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err == nil {
		for res.Next() {
			res.Scan(&billInfo.GroupID,
				&billInfo.Date,
				&billInfo.HDDSize,
				&billInfo.SSDSize,
				&billInfo.NVMESize,
				&billInfo.HDDChargePerGB,
				&billInfo.SSDChargePerGB,
				&billInfo.NVMEChargePerGB,
				&billInfo.DiscountRate)
		}
	}
	return &billInfo, err
}

func (bill *Billing) ReadBillingData(groupID *[]int32, dateStart, dateEnd, billType string) (*[][]model.Bill, *errors.HccErrorStack) {
	var billList [][]model.Bill
	errStack := errors.NewHccErrorStack()

	for _, gid := range *groupID {
		res, err := dao.GetBill(int(gid), dateStart, dateEnd, billType)
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

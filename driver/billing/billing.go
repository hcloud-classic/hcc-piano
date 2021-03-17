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

func (bil *Billing) UpdateNetworkBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateNodeBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateServerBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateVolumeBillingInfo(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateBillingData(groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) UpdateAllBillingData() *errors.HccErrorStack {
	return nil
}

func (bill *Billing) ReadNetworkBillingInfo(dateStart, dateEnd time.Time, groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) ReadNodeBillingInfo(dateStart, dateEnd time.Time, groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) ReadServerBillingInfo(dateStart, dateEnd time.Time, groupID *[]int32) *errors.HccErrorStack {
	return nil
}

func (bill *Billing) ReadVolumeBillingInfo(dateStart, dateEnd time.Time, groupID *[]int32) *errors.HccErrorStack {
	return nil
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

func (bill *Billing) ReadAllBillingData(dateStart, dateEnd time.Time, billType string) (*[]*model.Bill, *errors.HccErrorStack) {
	return nil, nil
}

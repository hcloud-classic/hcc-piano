package billing

import (
	"time"

	"hcc/piano/lib/logger"

	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var BillingDriver *Billing = &Billing{
	lastUpdate:  time.Now(),
	updateTimer: nil,
	StopTimer:   nil,
}

func reserveRegisterUpdateTimer(groupIDList *[]int32) {
	defer BillingDriver.UpdateBillingInfo(groupIDList)
	logger.Logger.Println("Register billing info update timer")

	now := time.Now().Add(1 * time.Hour)
	<-time.After(time.Until(time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), 0, 0, 0, now.Location())))

	BillingDriver.RunUpdateTimer()
}

func Init() *errors.HccError {

	groupIDList := []int32{1000, 1001, 1002}
	logger.Logger.Println("Update billing info in boot up time")
	BillingDriver.UpdateBillingInfo(&groupIDList)

	go reserveRegisterUpdateTimer(&groupIDList)

	return nil
}

func End() {
	if BillingDriver.StopTimer != nil {
		BillingDriver.StopTimer()
	}
}

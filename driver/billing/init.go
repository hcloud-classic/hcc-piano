package billing

import (
	"time"

	"hcc/piano/lib/logger"

	errors "innogrid.com/hcloud-classic/hcc_errors"
)

var BillingDriver *Billing = nil

func reserveRegisterUpdateTimer() {
	defer logger.Logger.Println("Update billing info.")

	now := time.Now().Add(1 * time.Hour)
	<-time.After(time.Until(time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), 0, 0, 0, now.Location())))

	logger.Logger.Println("Register billing info update timer")
	BillingDriver.RunUpdateTimer()
}

func Init() *errors.HccError {
	BillingDriver = &Billing{
		lastUpdate:  time.Now(),
		updateTimer: nil,
		StopTimer:   nil,
	}

	logger.Logger.Println("Update billing info in boot up time")

	go reserveRegisterUpdateTimer()

	return nil
}

func End() {
	if BillingDriver.StopTimer != nil {
		BillingDriver.StopTimer()
	}
}

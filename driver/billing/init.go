package billing

import (
	"hcc/piano/lib/config"
	"time"

	"hcc/piano/lib/logger"
)

var DriverBilling = &Billing{
	lastUpdate:  time.Now(),
	updateTimer: nil,
	StopTimer:   nil,
}

func reserveRegisterUpdateTimer() {
	defer DriverBilling.UpdateBillingInfo()
	logger.Logger.Println("Register billing info update timer")

	now := time.Now().Add(time.Duration(config.Billing.UpdateInterval) * time.Second)
	<-time.After(time.Until(time.Date(now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), 0, now.Location())))

	DriverBilling.RunUpdateTimer()
}

func Init() {
	logger.Logger.Println("Update billing info in boot up time")
	DriverBilling.UpdateBillingInfo()

	go reserveRegisterUpdateTimer()
}

func End() {
	if DriverBilling.StopTimer != nil {
		DriverBilling.StopTimer()
	}
}

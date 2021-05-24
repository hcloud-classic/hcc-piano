package billing

import "time"

func getTodayByNumString() string {
	currentTime := time.Now()
	return currentTime.Format("060102")
}

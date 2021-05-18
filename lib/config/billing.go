package config

type billing struct {
	UpdateInterval    int64 `goconf:"billing:billing_update_interval_sec"`     // UpdateInterval : Billing update interval (Seconds)
}

// Billing : Billing config structure
var Billing billing

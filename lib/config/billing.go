package config

type billing struct {
	Debug                     string `goconf:"billing:billing_debug"`                        // Debug : Enable debug logs for billing
	UpdateInterval    int64 `goconf:"billing:billing_update_interval_sec"`     // UpdateInterval : Billing update interval (Seconds)
}

// Billing : Billing config structure
var Billing billing

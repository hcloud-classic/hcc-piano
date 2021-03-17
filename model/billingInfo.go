package model

type Bill struct {
	BillID        uint64  `json:"bill_id"`
	ChargeNode    float32 `json:"charge_node"`
	ChargeServer  float32 `json:"charge_server"`
	ChargeNetwork float32 `json:"charge_network"`
	ChargeVolume  float32 `json:"charge_volume"`
}

type DiscountInfo struct {
	DiscountID    int     `json:"disount_id"`
	GroupID       int     `json:"group_id"`
	Expired       bool    `json:"expired"`
	Target        string  `json:"target"`
	DiscountRate  float32 `json:"discount_rate"`
	DiscountStart string  `json:"discount_start"`
	DiscountEnd   string  `json:"discount_end"`
	DiscountDesc  string  `json:"discount_desc"`
}

type NetworkBill struct {
	GroupID            int     `json:"group_id"`
	Date               string  `json:"date"`
	SubnetCount        int     `json:"subnet_count"`
	AIPCount           int     `json:"adaptiveip_count"`
	SubnetChargePerCnt float32 `json:"subnet_charge_per_cnt"`
	AIPChargePerCnt    float32 `json:"adaptiveip_charge_per_cnt"`
	DiscountRate       float32 `json:"discount_rate"`
}

type NodeBill struct {
	GroupID      int     `json:"group_id"`
	Date         string  `json:"date"`
	NodeUUID     string  `json:"node_uuid"`
	DefChargeCPU float32 `json:"default_charge_cpu"`
	DefChargeMEM float32 `json:"default_charge_memory"`
	DefChargeNIC float32 `json:"default_charge_nic"`
	DiscountRate float32 `json:"discount_rate"`
}

type ServerBill struct {
	GroupID            int     `json:"group_id"`
	Date               string  `json:"date"`
	ServerUUID         string  `json:"server_uuid"`
	NetworkTraffic     uint64  `json:"network_traffic"`
	TrafficChargePerKB float32 `json:"traffic_charge_per_KB"`
	DiscountRate       float32 `json:"discount_rate"`
}

type VolumeBill struct {
	GroupID         int     `json:"group_id"`
	Date            string  `json:"date"`
	HDDSize         uint64  `json:"hdd_size"`
	SSDSize         uint64  `json:"ssd_size"`
	NVMESize        uint64  `json:"nvme_size"`
	HDDChargePerKB  float32 `json:"hdd_charge_per_KB"`
	SSDChargePerKB  float32 `json:"ssd_charge_per_KB"`
	NVMEChargePerKB float32 `json:"nvme_charge_per_KB"`
	DiscountRate    float32 `json:"discount_rate"`
}

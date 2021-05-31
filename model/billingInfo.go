package model

type Bill struct {
	Date          string `json:"date"`
	GroupID       int64  `json:"group_id"`
	ChargeNode    int64  `json:"charge_node"`
	ChargeServer  int64  `json:"charge_server"`
	ChargeNetwork int64  `json:"charge_network"`
	ChargeVolume  int64  `json:"charge_volume"`
}

type NetworkBill struct {
	GroupID          int64 `json:"group_id"`
	ChargeSubnet     int64 `json:"charge_subnet"`
	ChargeAdaptiveIP int64 `json:"charge_adaptive_ip"`
}

type NodeBill struct {
	GroupID   int64  `json:"group_id"`
	NodeUUID  string `json:"node_uuid"`
	ChargeCPU int64  `json:"charge_cpu"`
	ChargeMEM int64  `json:"charge_memory"`
	ChargeNIC int64  `json:"charge_nic"`
}

type ServerBill struct {
	GroupID       int64  `json:"group_id"`
	ServerUUID    string `json:"server_uuid"`
	ChargeTraffic int64  `json:"charge_traffic"`
}

type VolumeBill struct {
	GroupID   int64  `json:"group_id"`
	Date      string `json:"date"`
	ChargeSSD int64  `json:"charge_ssd"`
	ChargeHDD int64  `json:"charge_hdd"`
}

type DailyBill struct {
	Date          string `json:"date"`
	GroupID       int64  `json:"group_id"`
	ChargeNode    int64  `json:"charge_node"`
	ChargeServer  int64  `json:"charge_server"`
	ChargeNetwork int64  `json:"charge_network"`
	ChargeVolume  int64  `json:"charge_volume"`
}

type BillDetail struct {
	BillID        int            `json:"bill_id"`
	DetailNode    *[]NodeBill    `json: "detail_node"`
	DetailServer  *[]ServerBill  `json: "detail_server"`
	DetailNetwork *[]NetworkBill `json: "detail_network"`
	DetailVolume  *[]VolumeBill  `json: "detail_volume"`
}

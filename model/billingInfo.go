package model

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

type NetworkBill struct {
	GroupID          int64 `json:"group_id"`
	ChargeSubnet     int64 `json:"charge_subnet"`
	ChargeAdaptiveIP int64 `json:"charge_adaptive_ip"`
}

type VolumeBill struct {
	GroupID   int64  `json:"group_id"`
	Date      string `json:"date"`
	ChargeSSD int64  `json:"charge_ssd"`
	ChargeHDD int64  `json:"charge_hdd"`
}

type Bill struct {
	Date          string `json:"date"`
	GroupID       int64  `json:"group_id"`
	ChargeNode    int64  `json:"charge_node"`
	ChargeServer  int64  `json:"charge_server"`
	ChargeNetwork int64  `json:"charge_network"`
	ChargeVolume  int64  `json:"charge_volume"`
}

type DailyBill struct {
	Date          string `json:"date"`
	GroupID       int64  `json:"group_id"`
	ChargeNode    int64  `json:"charge_node"`
	ChargeServer  int64  `json:"charge_server"`
	ChargeNetwork int64  `json:"charge_network"`
	ChargeVolume  int64  `json:"charge_volume"`
}

type Node struct {
	UUID     string `json:"uuid"`
	CPUCores int    `json:"cpu_cores"`
	Memory   int    `json:"memory"`
	NICSpeed string `json:"nic_speed"`
}

type DetailNode struct {
	Node     Node     `json:"node"`
	NodeBill NodeBill `json:"node_bill"`
}

type BillDetail struct {
	Date       int         `json:"date"`
	GroupID    int64       `json:"group_id"`
	DetailNode *[]NodeBill `json:"detail_node"`

	DetailServer  *[]ServerBill  `json: "detail_server"`
	DetailNetwork *[]NetworkBill `json: "detail_network"`
	DetailVolume  *[]VolumeBill  `json: "detail_volume"`
}

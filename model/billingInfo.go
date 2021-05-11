package model

type Bill struct {
	BillID        uint64  `json:"bill_id"`
	ChargeNode    int64 `json:"charge_node"`
	ChargeServer  int64 `json:"charge_server"`
	ChargeNetwork int64 `json:"charge_network"`
	ChargeVolume  int64 `json:"charge_volume"`
}

type NetworkBill struct {
	GroupID            int     `json:"group_id"`
	Date               string  `json:"date"`
	SubnetCharge int64 `json:"subnet_charge"`
	AdaptiveIPCharge    int64 `json:"adaptive_ip_charge"`
}

type NodeBill struct {
	GroupID      int     `json:"group_id"`
	Date         string  `json:"date"`
	NodeUUID     string  `json:"node_uuid"`
	ChargeCPU    int64 `json:"charge_cpu"`
	ChargeMEM    int64 `json:"charge_memory"`
	ChargeNIC    int64 `json:"charge_nic"`
}

type ServerBill struct {
	GroupID            int     `json:"group_id"`
	Date               string  `json:"date"`
	ServerUUID         string  `json:"server_uuid"`
	// TODO : RunningTime
	NetworkTraffic     uint64  `json:"network_traffic"`
	TrafficChargePerKB float32 `json:"traffic_charge_per_KB"`
}

type VolumeBill struct {
	GroupID         int     `json:"group_id"`
	Date            string  `json:"date"`
	HDDCharge  float32 `json:"hdd_charge"`
	SSDCharge  float32 `json:"ssd_charge"`
}

type BillDetail struct {
	BillID        int          `json:"bill_id"`
	DetailNode    *NodeBill    `json: "detail_node"`
	DetailServer  *ServerBill  `json: "detail_server"`
	DetailNetwork *NetworkBill `json: "detail_network"`
	DetailVolume  *VolumeBill  `json: "detail_volume"`
}

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
	GroupName     string `json:"group_name"`
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
	Uptime   string `json:"uptime"`
}

type DetailNode struct {
	Node     Node     `json:"node"`
	NodeBill NodeBill `json:"node_bill"`
}

type Server struct {
	Name           string `json:"name"`
	NetworkTraffic string `json:"network_traffic"`
}

type DetailServer struct {
	Server     Server     `json:"server"`
	ServerBill ServerBill `json:"server_bill"`
}

type Volume struct {
	UUID      string `json:"uuid"`
	Pool      string `json:"pool"`
	UsageType string `json:"usage_type"`
	DiskType  string `json:"disk_type"`
	DiskSize  int    `json:"disk_size"`
	Cost      int64  `json:"cost"`
}

type DetailVolume struct {
	Volumes    []Volume   `json:"volumes"`
	VolumeBill VolumeBill `json:"volume_bill"`
}

type Subnet struct {
	SubnetName string `json:"subnet_name"`
	DomainName string `json:"domain_name"`
	NetworkIP  string `json:"network_ip"`
	GatewayIP  string `json:"gateway_ip"`
	Cost       int64  `json:"cost"`
}

type AdaptiveIP struct {
	ServerName     string `json:"server_name"`
	PublicIP       string `json:"public_ip"`
	PrivateIP      string `json:"private_ip"`
	PrivateGateway string `json:"private_gateway"`
}

type DetailNetwork struct {
	Subnets     []Subnet     `json:"subnet"`
	AdaptiveIPs []AdaptiveIP `json:"adaptive_ip"`
	NetworkBill NetworkBill  `json:"network_bill"`
}

type BillDetail struct {
	Date          int             `json:"date"`
	GroupID       int64           `json:"group_id"`
	GroupName     string          `json:"group_name"`
	DetailNode    *[]DetailNode   `json:"detail_node"`
	DetailServer  *[]DetailServer `json:"detail_server"`
	DetailVolume  *DetailVolume   `json:"detail_volume"`
	DetailNetwork *DetailNetwork  `json:"detail_network"`
}

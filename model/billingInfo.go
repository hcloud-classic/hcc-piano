package model

type NodeBill struct {
	GroupID   int64  `json:"group_id"`
	Date      string `json:"date"`
	NodeUUID  string `json:"node_uuid"`
	ChargeCPU int64  `json:"charge_cpu"`
	ChargeMEM int64  `json:"charge_memory"`
	ChargeNIC int64  `json:"charge_nic"`
}

type ServerBill struct {
	GroupID       int64  `json:"group_id"`
	Date          string `json:"date"`
	ServerUUID    string `json:"server_uuid"`
	ChargeTraffic int64  `json:"charge_traffic"`
}

type SubnetBill struct {
	GroupID      int64  `json:"group_id"`
	Date         string `json:"date"`
	SubnetUUID   string `json:"subnet_uuid"`
	ChargeSubnet int64  `json:"charge_subnet"`
}

type AdaptiveIPBill struct {
	GroupID          int64  `json:"group_id"`
	Date             string `json:"date"`
	ServerUUID       string `json:"server_uuid"`
	ChargeAdaptiveIP int64  `json:"charge_adaptiveip"`
}

type VolumeBill struct {
	GroupID    int64  `json:"group_id"`
	Date       string `json:"date"`
	VolumeUUID string `json:"volume_uuid"`
	ChargeSSD  int64  `json:"charge_ssd"`
	ChargeHDD  int64  `json:"charge_hdd"`
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

type Subnet struct {
	SubnetName string `json:"subnet_name"`
	DomainName string `json:"domain_name"`
	NetworkIP  string `json:"network_ip"`
	GatewayIP  string `json:"gateway_ip"`
}

type AdaptiveIP struct {
	ServerName     string `json:"server_name"`
	PublicIP       string `json:"public_ip"`
	PrivateIP      string `json:"private_ip"`
	PrivateGateway string `json:"private_gateway"`
}

type DetailSubnet struct {
	Subnet     Subnet     `json:"subnet"`
	SubnetBill SubnetBill `json:"subnet_bill"`
}

type DetailAdaptiveIP struct {
	AdaptiveIP     AdaptiveIP     `json:"adaptive_ip"`
	AdaptiveIPBill AdaptiveIPBill `json:"adaptive_ip_bill"`
}

type Volume struct {
	UUID      string `json:"uuid"`
	Pool      string `json:"pool"`
	UsageType string `json:"usage_type"`
	DiskType  string `json:"disk_type"`
	DiskSize  int    `json:"disk_size"`
}

type DetailVolume struct {
	Volume     Volume     `json:"volume"`
	VolumeBill VolumeBill `json:"volume_bill"`
}

type BillDetail struct {
	Date             string              `json:"date"`
	GroupID          int64               `json:"group_id"`
	DetailNode       *[]DetailNode       `json:"detail_node"`
	DetailServer     *[]DetailServer     `json:"detail_server"`
	DetailSubnet     *[]DetailSubnet     `json:"detail_subnet"`
	DetailAdaptiveIP *[]DetailAdaptiveIP `json:"detail_adaptive_ip"`
	DetailVolume     *[]DetailVolume     `json:"detail_volume"`
}

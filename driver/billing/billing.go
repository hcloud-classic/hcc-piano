package billing

import (
	"errors"
	"fmt"
	"hcc/piano/lib/config"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
	"time"

	"hcc/piano/action/grpc/client"
	"hcc/piano/dao"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
)

type Billing struct {
	updateTimer *time.Ticker
	StopTimer   func()
	IsRunning   bool
}

func (bill *Billing) RunUpdateTimer() {
	if bill.updateTimer == nil {
		bill.updateTimer = time.NewTicker(time.Duration(config.Billing.UpdateInterval) * time.Second)
	} else {
		// upper go v1.15
		// bill.updateTimer.Reset(duration)
		bill.updateTimer.Stop()
		bill.updateTimer = time.NewTicker(time.Duration(config.Billing.UpdateInterval) * time.Second)

		return
	}
	done := make(chan bool)
	bill.StopTimer = func() {
		done <- true
	}

	go func() {
		for true {
			select {
			case <-done:
				logger.Logger.Println("RunUpdateTimer(): Stopping billing update timer")

				bill.updateTimer.Stop()
				bill.updateTimer = nil

				return
			case <-bill.updateTimer.C:
				if config.Billing.Debug == "on" {
					logger.Logger.Println("RunUpdateTimer(): Updating billing information")
				}
				DriverBilling.UpdateBillingInfo()
				break
			}
		}
	}()
}

func (bill *Billing) UpdateBillingInfo() {
	bill.IsRunning = true

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting group list")
	}
	resGetGroupList, err := client.RC.GetGroupList()
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): GetGroupList(): " + err.Error())
		bill.IsRunning = false
		return
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting node_billing_info")
	}
	nodeBillList, err := getNodeBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getNodeBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting node_billing_info")
		}
		err = dao.InsertNodeBillingInfo(nodeBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertNodeBillingInfo(): " + err.Error())
		}
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting server_billing_info")
	}
	serverBillList, err := getServerBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getServerBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting server_billing_info")
		}
		err = dao.InsertServerBillingInfo(serverBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertServerBillingInfo(): " + err.Error())
		}
	}

	networkBillList, err := getNetworkBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getNetworkBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting network_billing_info")
		}
		err = dao.InsertNetworkBillingInfo(networkBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertNetworkBillingInfo(): " + err.Error())
		}
	}

	volumeBillList, err := getVolumeBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getVolumeBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting volume_billing_info")
		}
		err = dao.InsertVolumeBillingInfo(volumeBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertVolumeBillingInfo(): " + err.Error())
		}
	}

	if err == nil {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Getting daily_info")
		}
		dailyBillList := dao.GetDailyInfo(resGetGroupList.Group, nodeBillList, serverBillList, networkBillList, volumeBillList)
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting daily_info")
		}
		err = dao.InsertDailyInfo(dailyBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertDailyInfo(): " + err.Error())
		}
	}

	bill.IsRunning = false
}

func calcNicSpeed(nicSpeedMbps int32) string {
	var nicSpeed = "error"

	switch nicSpeedMbps {
	case 10:
		nicSpeed = "10Mbps"
	case 100:
		nicSpeed = "100Mbps"
	case 1000:
		nicSpeed = "1Gbps"
	case 2500:
		nicSpeed = "2.5Gbps"
	case 5000:
		nicSpeed = "5Gbps"
	case 10000:
		nicSpeed = "10Gbps"
	case 20000:
		nicSpeed = "20Gbps"
	case 40000:
		nicSpeed = "40Gbps"
	}

	return nicSpeed
}

func calcNodeUptime(uptimeMs int64) string {
	var uptimeSec = 0
	var uptimeMin = 0
	var uptimeHour = 0
	var uptimeDay = 0
	var uptimeStr = ""

	if uptimeMs >= 1000 {
		uptimeSec = int(uptimeMs / int64(1000))
	} else {
		return strconv.Itoa(int(uptimeMs)) + "ms"
	}
	if uptimeSec >= 60 {
		uptimeMin = uptimeSec / 60
		uptimeSec = uptimeSec % 60
	}
	if uptimeMin >= 60 {
		uptimeHour = uptimeMin / 60
		uptimeMin = uptimeMin % 60
	}
	if uptimeHour >= 24 {
		uptimeDay = uptimeHour / 24
		uptimeHour = uptimeHour % 24
	}

	if uptimeDay > 0 {
		uptimeStr = strconv.Itoa(uptimeDay) + "d"
	}
	if uptimeHour > 0 {
		if uptimeDay > 0 {
			uptimeStr += " "
		}
		uptimeStr = uptimeStr + strconv.Itoa(uptimeHour) + "h"
	}
	if uptimeMin > 0 {
		if uptimeDay > 0 || uptimeHour > 0 {
			uptimeStr += " "
		}
		uptimeStr = uptimeStr + strconv.Itoa(uptimeMin) + "m"
	}
	if uptimeSec > 0 {
		if uptimeDay > 0 || uptimeHour > 0 || uptimeMin > 0 {
			uptimeStr += " "
		}
		uptimeStr = uptimeStr + strconv.Itoa(uptimeSec) + "s"
	}

	return uptimeStr
}

func (bill *Billing) readNodeBillingInfo(groupID int64, date, billType string) (*[]model.DetailNode, error) {
	var detailNodes []model.DetailNode
	var uptimeMS int64

	res, err := dao.GetBillInfo(groupID, date, billType, "node")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var detailNode model.DetailNode

		_ = res.Scan(&detailNode.NodeBill.GroupID,
			&detailNode.NodeBill.Date,
			&detailNode.NodeBill.NodeUUID,
			&detailNode.NodeBill.ChargeCPU,
			&detailNode.NodeBill.ChargeMEM,
			&detailNode.NodeBill.ChargeNIC,
			&uptimeMS)

		resGetNode, err := client.RC.GetNode(detailNode.NodeBill.NodeUUID)
		if err != nil {
			return nil, err
		}

		detailNode.Node = model.Node{
			UUID:     resGetNode.Node.UUID,
			CPUCores: int(resGetNode.Node.CPUCores),
			Memory:   int(resGetNode.Node.Memory),
			NICSpeed: calcNicSpeed(resGetNode.Node.NicSpeedMbps),
			Uptime:   calcNodeUptime(uptimeMS),
		}

		detailNodes = append(detailNodes, detailNode)
	}

	return &detailNodes, err
}

func calcTraffic(trafficKB int64) string {
	if trafficKB > 999 {
		dot := trafficKB % 1024
		resultMB := float32(trafficKB) / 1024
		if resultMB > 999 {
			resultGB := resultMB / 1024
			if resultGB > 999 {
				resultTB := resultGB / 1024
				if dot != 0 {
					return fmt.Sprintf("%.2f", resultTB) + "TB"
				}

				return strconv.Itoa(int(resultTB)) + "TB"
			}
			if dot != 0 {
				return fmt.Sprintf("%.2f", resultGB) + "GB"
			}

			return strconv.Itoa(int(resultGB)) + "GB"
		}
		if dot != 0 {
			return fmt.Sprintf("%.2f", resultMB) + "MB"
		}

		return strconv.Itoa(int(resultMB)) + "MB"
	}

	return strconv.Itoa(int(trafficKB)) + "KB"
}

func (bill *Billing) readServerBillingInfo(groupID int64, date, billType string) (*[]model.DetailServer, error) {
	var detailServers []model.DetailServer
	var trafficKB int64

	res, err := dao.GetBillInfo(groupID, date, billType, "server")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var detailServer model.DetailServer

		_ = res.Scan(&detailServer.ServerBill.GroupID,
			&detailServer.ServerBill.Date,
			&detailServer.ServerBill.ServerUUID,
			&detailServer.ServerBill.ChargeTraffic,
			&trafficKB)

		resGetServer, err := client.RC.GetServer(detailServer.ServerBill.ServerUUID)
		if err != nil {
			return nil, err
		}

		detailServer.Server = model.Server{
			Name:           resGetServer.Server.ServerName,
			NetworkTraffic: calcTraffic(trafficKB),
		}

		detailServers = append(detailServers, detailServer)
	}

	return &detailServers, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int64, date, billType string) (*model.DetailVolume, error) {
	var detailVolume model.DetailVolume
	var volumes []model.Volume

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	resGetCharge, err := client.RC.GetCharge(groupID)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		_ = res.Scan(&detailVolume.VolumeBill.GroupID,
			&detailVolume.VolumeBill.Date,
			&detailVolume.VolumeBill.ChargeSSD,
			&detailVolume.VolumeBill.ChargeHDD)

		resGetServerList, err := client.RC.GetServerList(&pb.ReqGetServerList{
			Server: &pb.Server{
				GroupID: groupID,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, server := range resGetServerList.Server {
			resGetVolumeList, err := client.RC.GetVolumeList(&pb.ReqGetVolumeList{
				Volume: &pb.Volume{
					Action:     "read_list",
					ServerUUID: server.UUID,
				},
				// TODO : Should be control row and page later?
				Row:  10,
				Page: 1,
			})
			if err != nil {
				return nil, err
			}

			for _, volume := range resGetVolumeList.Volume {
				useType := strings.ToLower(volume.UseType)
				size, _ := strconv.Atoi(volume.Size)

				// TODO : Should it be change to identify SSD or HDD later?
				var diskType string
				var cost int64
				if useType == "os" {
					diskType = "SSD"
					cost = resGetCharge.Charge.ChargeSSDPerGB * int64(size)
				} else if useType == "data" {
					diskType = "HDD"
					cost = resGetCharge.Charge.ChargeHDDPerGB * int64(size)
				}

				volume := model.Volume{
					UUID:      volume.UUID,
					Pool:      volume.Pool,
					UsageType: volume.UseType,
					DiskType:  diskType,
					DiskSize:  size,
					Cost:      cost,
				}

				volumes = append(volumes, volume)
			}
		}
	}

	detailVolume.Volumes = volumes

	return &detailVolume, err
}

func (bill *Billing) readNetworkBillingInfo(groupID int64, date, billType string) (*model.DetailNetwork, error) {
	var detailNetwork model.DetailNetwork
	var subnets []model.Subnet
	var adaptiveIPs []model.AdaptiveIP

	res, err := dao.GetBillInfo(groupID, date, billType, "network")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	resGetCharge, err := client.RC.GetCharge(groupID)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		_ = res.Scan(&detailNetwork.NetworkBill.GroupID,
			&detailNetwork.NetworkBill.Date,
			&detailNetwork.NetworkBill.ChargeSubnet,
			&detailNetwork.NetworkBill.ChargeAdaptiveIP)

		resGetServerList, err := client.RC.GetServerList(&pb.ReqGetServerList{
			Server: &pb.Server{
				GroupID: groupID,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, server := range resGetServerList.Server {

			resGetSubnetList, err := client.RC.GetSubnetList(&pb.ReqGetSubnetList{
				Subnet: &pb.Subnet{
					ServerUUID: server.UUID,
				},
			})
			if err != nil {
				return nil, err
			}

			for _, subnet := range resGetSubnetList.Subnet {
				subnet := model.Subnet{
					SubnetName: subnet.SubnetName,
					DomainName: subnet.DomainName,
					NetworkIP:  subnet.NetworkIP,
					GatewayIP:  subnet.Gateway,
					Cost:       resGetCharge.Charge.ChargeSubnetPerCnt,
				}

				subnets = append(subnets, subnet)
			}

			resGetAdaptiveIPServerList, err := client.RC.GetAdaptiveIPServerList(&pb.ReqGetAdaptiveIPServerList{
				AdaptiveipServer: &pb.AdaptiveIPServer{
					ServerUUID: server.UUID,
				},
			})
			if err != nil {
				return nil, err
			}

			for _, adaptiveIP := range resGetAdaptiveIPServerList.AdaptiveipServer {
				adaptiveIP := model.AdaptiveIP{
					ServerName:     server.ServerName,
					PublicIP:       adaptiveIP.PublicIP,
					PrivateIP:      adaptiveIP.PrivateIP,
					PrivateGateway: adaptiveIP.PrivateGateway,
					Cost: resGetCharge.Charge.ChargeAdaptiveIPPerCnt,
				}

				adaptiveIPs = append(adaptiveIPs, adaptiveIP)
			}
		}
	}

	detailNetwork.Subnets = subnets
	detailNetwork.AdaptiveIPs = adaptiveIPs

	return &detailNetwork, err
}

func (bill *Billing) ReadBillingData(groupID *[]int64, dateStart, dateEnd, billType string, row, page int) (*[]model.Bill, error) {
	var billList []model.Bill

	res, err := dao.GetBill(groupID, dateStart, dateEnd, billType, row, page)
	if err != nil {
		logger.Logger.Println("ReadBillingData(): dao.GetBill(): " + err.Error())
		return &billList, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var bill model.Bill
		_ = res.Scan(&bill.Date,
			&bill.GroupID,
			&bill.GroupName,
			&bill.ChargeNode,
			&bill.ChargeServer,
			&bill.ChargeNetwork,
			&bill.ChargeVolume)
		billList = append(billList, bill)
	}

	return &billList, nil
}

func (bill *Billing) ReadBillingDetail(groupID int64, date, billType string) (*model.BillDetail, error) {
	var err error
	var returnErr error = nil
	var billingDetail model.BillDetail

	billingDetail.DetailNode, err = bill.readNodeBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNodeBillingInfo(): " + err.Error())
		returnErr = errors.New(err.Error())
	}

	billingDetail.DetailServer, err = bill.readServerBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readServerBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	billingDetail.DetailVolume, err = bill.readVolumeBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readVolumeBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	billingDetail.DetailNetwork, err = bill.readNetworkBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readNetworkBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	return &billingDetail, returnErr
}

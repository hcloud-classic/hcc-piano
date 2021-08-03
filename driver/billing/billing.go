package billing

import (
	"errors"
	"fmt"
	"hcc/piano/lib/config"
	"innogrid.com/hcloud-classic/hcc_errors"
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

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting subnet_billing_info")
	}
	subnetBillList, err := getSubnetBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getSubnetBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting subnet_billing_info")
		}
		err = dao.InsertSubnetBillingInfo(subnetBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertSubnetBillingInfo(): " + err.Error())
		}
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting adaptiveip_billing_info")
	}
	adaptiveIPBillList, err := getAdaptiveIPBillingInfo(resGetGroupList.Group)
	if err != nil {
		logger.Logger.Println("UpdateBillingInfo(): getAdaptiveIPBillingInfo(): " + err.Error())
	} else {
		if config.Billing.Debug == "on" {
			logger.Logger.Println("RunUpdateTimer(): Inserting adaptiveip_billing_info")
		}
		err = dao.InsertAdaptiveIPBillingInfo(adaptiveIPBillList)
		if err != nil {
			logger.Logger.Println("UpdateBillingInfo(): InsertAdaptiveIPBillingInfo(): " + err.Error())
		}
	}

	if config.Billing.Debug == "on" {
		logger.Logger.Println("RunUpdateTimer(): Getting volume_billing_info")
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
		dailyBillList := dao.GetDailyInfo(resGetGroupList.Group)
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

		if resGetNode.HccErrorStack != nil && resGetNode.HccErrorStack.ErrStack != nil {
			if resGetNode.HccErrorStack.ErrStack[0].ErrCode == hcc_errors.FluteSQLNoResult {
				resGetNode.Node.UUID += " (Deleted)"
			} else {
				resGetNode.Node.UUID = "error"
			}
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

		if resGetServer.HccErrorStack != nil && resGetServer.HccErrorStack.ErrStack != nil {
			if resGetServer.HccErrorStack.ErrStack[0].ErrCode == hcc_errors.ViolinSQLNoResult {
				resGetServer.Server.ServerName = detailServer.ServerBill.ServerUUID + " (Deleted)"
			} else {
				resGetServer.Server.UUID = "error"
			}
		}

		detailServer.Server = model.Server{
			Name:           resGetServer.Server.ServerName,
			NetworkTraffic: calcTraffic(trafficKB),
		}

		detailServers = append(detailServers, detailServer)
	}

	return &detailServers, err
}

func (bill *Billing) readSubnetBillingInfo(groupID int64, date, billType string) (*[]model.DetailSubnet, error) {
	var detailSubnets []model.DetailSubnet

	res, err := dao.GetBillInfo(groupID, date, billType, "subnet")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var detailSubnet model.DetailSubnet

		_ = res.Scan(&detailSubnet.SubnetBill.GroupID,
			&detailSubnet.SubnetBill.Date,
			&detailSubnet.SubnetBill.SubnetUUID,
			&detailSubnet.SubnetBill.ChargeSubnet)

		resGetSubnet, err := client.RC.GetSubnet(detailSubnet.SubnetBill.SubnetUUID)
		if err != nil {
			return nil, err
		}

		if resGetSubnet.HccErrorStack != nil && resGetSubnet.HccErrorStack.ErrStack != nil {
			if resGetSubnet.HccErrorStack.ErrStack[0].ErrCode == hcc_errors.HarpSQLNoResult {
				resGetSubnet.Subnet.SubnetName = detailSubnet.SubnetBill.SubnetUUID + " (Deleted)"
			} else {
				resGetSubnet.Subnet.SubnetName = "error"
			}
		}

		detailSubnet.Subnet = model.Subnet{
			SubnetName: resGetSubnet.Subnet.SubnetName,
			DomainName: resGetSubnet.Subnet.DomainName,
			NetworkIP:  resGetSubnet.Subnet.NetworkIP,
			GatewayIP:  resGetSubnet.Subnet.Gateway,
		}

		detailSubnets = append(detailSubnets, detailSubnet)
	}

	return &detailSubnets, err
}

func (bill *Billing) readAdaptiveIPBillingInfo(groupID int64, date, billType string) (*[]model.DetailAdaptiveIP, error) {
	var detailAdaptiveIPs []model.DetailAdaptiveIP

	res, err := dao.GetBillInfo(groupID, date, billType, "adaptiveip")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var detailAdaptiveIP model.DetailAdaptiveIP
		var serverName string

		_ = res.Scan(&detailAdaptiveIP.AdaptiveIPBill.GroupID,
			&detailAdaptiveIP.AdaptiveIPBill.Date,
			&detailAdaptiveIP.AdaptiveIPBill.ServerUUID,
			&detailAdaptiveIP.AdaptiveIPBill.ChargeAdaptiveIP)

		resGetAdaptiveIPServer, err := client.RC.GetAdaptiveIPServer(detailAdaptiveIP.AdaptiveIPBill.ServerUUID)
		if err != nil {
			return nil, err
		}

		if resGetAdaptiveIPServer.HccErrorStack != nil && resGetAdaptiveIPServer.HccErrorStack.ErrStack != nil {
			if resGetAdaptiveIPServer.HccErrorStack.ErrStack[0].ErrCode == hcc_errors.HarpSQLNoResult {
				serverName = detailAdaptiveIP.AdaptiveIPBill.ServerUUID + " (AdaptiveIP Deleted)"
			} else {
				serverName = "error"
			}
		} else {
			resGetServer, err := client.RC.GetServer(resGetAdaptiveIPServer.AdaptiveipServer.ServerUUID)
			if err != nil {
				return nil, err
			}
			serverName = resGetServer.Server.ServerName
		}

		detailAdaptiveIP.AdaptiveIP = model.AdaptiveIP{
			ServerName:     serverName,
			PublicIP:       resGetAdaptiveIPServer.AdaptiveipServer.PublicIP,
			PrivateIP:      resGetAdaptiveIPServer.AdaptiveipServer.PrivateIP,
			PrivateGateway: resGetAdaptiveIPServer.AdaptiveipServer.PrivateGateway,
		}

		detailAdaptiveIPs = append(detailAdaptiveIPs, detailAdaptiveIP)
	}

	return &detailAdaptiveIPs, err
}

func (bill *Billing) readVolumeBillingInfo(groupID int64, date, billType string) (*[]model.DetailVolume, error) {
	var detailVolumes []model.DetailVolume

	res, err := dao.GetBillInfo(groupID, date, billType, "volume")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	var volumes []model.Volume

	resGetServerList, err := client.RC.GetServerList(&pb.ReqGetServerList{
		Server: &pb.Server{},
	})
	if err != nil {
		return nil, err
	}

	for _, server := range resGetServerList.Server {
		resGetVolumeList, err := client.RC.GetVolumeList(&pb.ReqGetVolumeList{
			Volume: &pb.Volume{
				Action:     "single_server_allocated",
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
			if useType == "os" {
				diskType = "SSD"
			} else if useType == "data" {
				diskType = "HDD"
			}

			volumes = append(volumes, model.Volume{
				UUID:      volume.UUID,
				Pool:      volume.Pool,
				UsageType: volume.UseType,
				DiskType:  diskType,
				DiskSize:  size,
			})
		}
	}

	for res.Next() {
		var detailVolume model.DetailVolume

		_ = res.Scan(&detailVolume.VolumeBill.GroupID,
			&detailVolume.VolumeBill.Date,
			&detailVolume.VolumeBill.VolumeUUID,
			&detailVolume.VolumeBill.ChargeSSD,
			&detailVolume.VolumeBill.ChargeHDD)

		var volumeFound bool
		for _, volume := range volumes {
			if detailVolume.VolumeBill.VolumeUUID == volume.UUID {
				detailVolume.Volume = model.Volume{
					UUID:      volume.UUID,
					Pool:      volume.Pool,
					UsageType: volume.UsageType,
					DiskType:  volume.DiskType,
					DiskSize:  volume.DiskSize,
				}
				volumeFound = true

				break
			}
		}

		if !volumeFound {
			detailVolume.Volume = model.Volume{
				UUID: detailVolume.VolumeBill.VolumeUUID + " (Deleted)",
			}
		}

		detailVolumes = append(detailVolumes, detailVolume)
	}

	return &detailVolumes, err
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

	billingDetail.Date = date
	billingDetail.GroupID = groupID

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

	billingDetail.DetailSubnet, err = bill.readSubnetBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readSubnetBillingInfo(): " + err.Error())
		if returnErr == nil {
			returnErr = errors.New("")
		}
		returnErr = errors.New(returnErr.Error() + "\n" + err.Error())
	}

	billingDetail.DetailAdaptiveIP, err = bill.readAdaptiveIPBillingInfo(groupID, date, billType)
	if err != nil {
		logger.Logger.Println("ReadBillingDetail(): bill.readAdaptiveIPBillingInfo(): " + err.Error())
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

	return &billingDetail, returnErr
}

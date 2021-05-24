package billing

import (
	"errors"
	"hcc/piano/action/grpc/client"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
)

var nicSpeedsMbps = []int32{10, 100, 1000, 2500, 5000, 10000, 20000, 40000}

func getChargeNIC(chargeNICList string, nicSpeedMbps int32) (int64, error) {
	chargeNICs := strings.Split(chargeNICList, ",")
	if len(chargeNICs) != len(nicSpeedsMbps) {
		goto OUT
	}

	for i, speed := range nicSpeedsMbps {
		if speed == nicSpeedMbps {
			chargeNic, _ := strconv.Atoi(chargeNICs[i])
			return int64(chargeNic), nil
		}
	}

OUT:
	return 0, errors.New("invalid charge_nic_list")
}

func getNodeBillingInfo(groupList []*pb.Group) (*[]model.NodeBill, error) {
	var billList []model.NodeBill

	for _, group := range groupList {
		resGetCharge, err := client.RC.GetCharge(group.Id)
		if err != nil {
			return nil, err
		}

		resGetNodeList, err := client.RC.GetNodeList(&pb.ReqGetNodeList{
			Node: &pb.Node{
				GroupID: group.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, node := range resGetNodeList.Node {
			var chargeCPU int64 = 0
			var chargeMEM int64 = 0
			var chargeNIC int64 = 0

			if strings.ToLower(node.Status) == "on" {
				resGetNodeUptime, err := client.RC.GetNodeUptime(&pb.ReqGetNodeUptime{
					NodeUUID: node.UUID,
					Day:      getTodayByNumString(),
				})
				if err != nil {
					logger.Logger.Println("getNodeBillingInfo(): Failed to get nodeUptime of nodeUUID=" + node.UUID)
					continue
				}

				nodeUptimeMs := resGetNodeUptime.NodeUptime.UptimeMs

				chargeCPU = int64(float64(resGetCharge.Charge.ChargeCPUPerCore * int64(node.CPUCores)) / float64(24 * 3600 * 1000) * float64(nodeUptimeMs))
				chargeMEM = int64(float64(resGetCharge.Charge.ChargeMemoryPerGB * int64(node.Memory)) / float64(24 * 3600 * 1000) * float64(nodeUptimeMs))
				chargeNIC, err = getChargeNIC(resGetCharge.Charge.ChargeNicList, node.NicSpeedMbps)
				if err != nil {
					logger.Logger.Println("getNodeBillingInfo(): Failed to get chargeNIC of nodeUUID=" + node.UUID)
					continue
				}
			}

			billList = append(billList, model.NodeBill{
				GroupID:   int(group.Id),
				NodeUUID:  node.UUID,
				ChargeCPU: chargeCPU,
				ChargeMEM: chargeMEM,
				ChargeNIC: chargeNIC,
			})
		}
	}

	return &billList, nil
}

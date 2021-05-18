package billing

import (
	"errors"
	"hcc/piano/action/grpc/client"
	"hcc/piano/lib/config"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
	"time"
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
			chargeNIC, err := getChargeNIC(resGetCharge.Charge.ChargeNicList, node.NicSpeedMbps)
			if err != nil {
				return nil, err
			}

			if strings.ToLower(node.Status) == "off" {
				continue
			}

			billList = append(billList, model.NodeBill{
				GroupID:   int(group.Id),
				NodeUUID:  node.UUID,
				ChargeCPU: resGetCharge.Charge.ChargeCPUPerCore * int64(node.CPUCores) / 30 / (24 * 3600) * config.Billing.UpdateInterval,
				ChargeMEM: resGetCharge.Charge.ChargeMemoryPerGB * int64(node.Memory) / 30 / (24 * 3600) * config.Billing.UpdateInterval,
				ChargeNIC: chargeNIC / 30 / (24 * 3600) * config.Billing.UpdateInterval,
			})
		}
	}

	return &billList, nil
}

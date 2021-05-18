package billing

import (
	"errors"
	"hcc/piano/action/grpc/client"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
	"time"
)

func getNodeBillingInfo(groupList []*pb.Group) (*[]model.NodeBill, error) {
	var billList []model.NodeBill

	now := time.Now()

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

			billList = append(billList, model.NodeBill{
				GroupID:   int(group.Id),
				Date:      now.Format("060102"),
				NodeUUID:  node.UUID,
				ChargeCPU: resGetCharge.Charge.ChargeCPUPerCore * int64(node.CPUCores),
				ChargeMEM: resGetCharge.Charge.ChargeMemoryPerGB * int64(node.Memory),
				ChargeNIC: chargeNIC,
			})
		}
	}

	return &billList, nil
}

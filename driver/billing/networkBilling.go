package billing

import (
	"hcc/piano/action/grpc/client"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
)

func getSubnetBillingInfo(groupList []*pb.Group) (*[]model.SubnetBill, error) {
	var billList []model.SubnetBill

	for _, group := range groupList {
		if group.Id == 1 {
			continue
		}

		resGetCharge, err := client.RC.GetCharge(group.Id)
		if err != nil {
			return nil, err
		}

		resGetSubnetList, err := client.RC.GetSubnetList(&pb.ReqGetSubnetList{
			Subnet: &pb.Subnet{
				GroupID: group.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, subnet := range resGetSubnetList.Subnet {
			billList = append(billList, model.SubnetBill{
				GroupID:      group.Id,
				SubnetUUID:   subnet.UUID,
				ChargeSubnet: resGetCharge.Charge.ChargeSubnetPerCnt,
			})
		}
	}

	return &billList, nil
}

func getAdaptiveIPBillingInfo(groupList []*pb.Group) (*[]model.AdaptiveIPBill, error) {
	var billList []model.AdaptiveIPBill

	for _, group := range groupList {
		if group.Id == 1 {
			continue
		}

		resGetCharge, err := client.RC.GetCharge(group.Id)
		if err != nil {
			return nil, err
		}

		resGetAdaptiveIPServerList, err := client.RC.GetAdaptiveIPServerList(&pb.ReqGetAdaptiveIPServerList{
			AdaptiveipServer: &pb.AdaptiveIPServer{
				GroupID: group.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, adaptiveIP := range resGetAdaptiveIPServerList.AdaptiveipServer {
			billList = append(billList, model.AdaptiveIPBill{
				GroupID:          group.Id,
				ServerUUID:       adaptiveIP.ServerUUID,
				ChargeAdaptiveIP: resGetCharge.Charge.ChargeAdaptiveIPPerCnt,
			})
		}
	}

	return &billList, nil
}

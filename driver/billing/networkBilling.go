package billing

import (
	"hcc/piano/action/grpc/client"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
)

func getNetworkBillingInfo(groupList []*pb.Group) (*[]model.NetworkBill, error) {
	var billList []model.NetworkBill

	for _, group := range groupList {
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

		resGetAdaptiveIPServerList, err := client.RC.GetAdaptiveIPServerList(&pb.ReqGetAdaptiveIPServerList{
			AdaptiveipServer: &pb.AdaptiveIPServer{
				GroupID: group.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		billList = append(billList, model.NetworkBill{
			GroupID:          int(group.Id),
			ChargeSubnet:     resGetCharge.Charge.ChargeSubnetPerCnt * int64(len(resGetSubnetList.Subnet)),
			ChargeAdaptiveIP: resGetCharge.Charge.ChargeAdaptiveIPPerCnt * int64(len(resGetAdaptiveIPServerList.AdaptiveipServer)),
		})
	}

	return &billList, nil
}

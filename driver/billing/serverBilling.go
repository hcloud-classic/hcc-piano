package billing

import (
	"hcc/piano/action/grpc/client"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
)

func getServerBillingInfo(groupList []*pb.Group) (*[]model.ServerBill, error) {
	var billList []model.ServerBill

	for _, group := range groupList {
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

		for _, adaptiveipServer := range resGetAdaptiveIPServerList.AdaptiveipServer {
			resGetTraffic, err := client.RC.GetTraffic(adaptiveipServer.ServerUUID, getTodayByNumString())
			if err != nil {
				logger.Logger.Println("getServerBillingInfo(): Failed to get traffic info for serverUUID=" + adaptiveipServer.ServerUUID)
				continue
			}

			trafficTotalKB := resGetTraffic.Traffic.TxKB + resGetTraffic.Traffic.RxKB

			billList = append(billList, model.ServerBill{
				GroupID:   int(group.Id),
				ServerUUID: adaptiveipServer.ServerUUID,
				ChargeTraffic: int64(resGetCharge.Charge.ChargeTrafficPerKB * float32(trafficTotalKB)),
			})
		}
	}

	return &billList, nil
}

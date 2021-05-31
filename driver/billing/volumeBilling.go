package billing

import (
	"hcc/piano/action/grpc/client"
	"hcc/piano/model"
	"innogrid.com/hcloud-classic/pb"
	"strconv"
	"strings"
)

func getVolumeBillingInfo(groupList []*pb.Group) (*[]model.VolumeBill, error) {
	var billList []model.VolumeBill

	for _, group := range groupList {
		if group.Id == 1 {
			continue
		}

		resGetCharge, err := client.RC.GetCharge(group.Id)
		if err != nil {
			return nil, err
		}

		resGetServerList, err := client.RC.GetServerList(&pb.ReqGetServerList{
			Server: &pb.Server{
				GroupID: group.Id,
			},
		})
		if err != nil {
			return nil, err
		}

		// TODO : Should it be change to identify SSD or HDD later?
		var ssdGBTotal = 0
		var hddGBTotal = 0

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

				if useType == "os" {
					ssdGBTotal += size
				} else if useType == "data" {
					hddGBTotal += size
				}
			}
		}

		billList = append(billList, model.VolumeBill{
			GroupID:   group.Id,
			ChargeSSD: resGetCharge.Charge.ChargeSSDPerGB * int64(ssdGBTotal),
			ChargeHDD: resGetCharge.Charge.ChargeHDDPerGB * int64(hddGBTotal),
		})
	}

	return &billList, nil
}

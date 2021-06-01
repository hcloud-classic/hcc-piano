package server

import (
	"context"
	"encoding/json"
	"hcc/piano/action/grpc/client"
	"hcc/piano/dao"
	"hcc/piano/model"
	"strconv"
	"strings"

	"hcc/piano/action/grpc/errconv"
	"hcc/piano/driver/billing"
	"hcc/piano/driver/influxdb"

	errors "innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

type pianoServer struct {
	pb.UnimplementedPianoServer
}

func (s *pianoServer) Telegraph(_ context.Context, in *pb.ReqMetricInfo) (*pb.ResMonitoringData, error) {
	series := influxdb.GetInfluxData(in)

	return series, nil
}

func (s *pianoServer) GetBillingData(_ context.Context, in *pb.ReqBillingData) (*pb.ResBillingData, error) {
	var data *[]model.Bill
	var count int
	var err error

	var groupID = &in.GroupID
	var groupIDAll []int64

	var resBillingData = pb.ResBillingData{
		BillingType:   "UNDEFINED",
		GroupID:       []int64{},
		Result:        nil,
		HccErrorStack: nil,
	}

	switch strings.ToLower(in.BillingType) {
	case "yearly":
		if len(in.DateStart) != 2 {
			goto WrongStartDate
		}
		if len(in.DateEnd) != 2 {
			goto WrongEndDate
		}
	case "monthly":
		if len(in.DateStart) != 4 {
			goto WrongStartDate
		}
		if len(in.DateEnd) != 4 {
			goto WrongEndDate
		}
	case "daily":
		if len(in.DateStart) != 6 {
			goto WrongStartDate
		}
		if len(in.DateEnd) != 6 {
			goto WrongEndDate
		}
	default:
		resBillingData.HccErrorStack = errconv.HccStackToGrpc(
			errors.NewHccErrorStack(
				errors.NewHccError(errors.PianoGrpcArgumentError, "-> Unsupport BillingType")))
		goto OUT
	}

	_, err = strconv.Atoi(in.DateStart)
	if err != nil {
		goto WrongStartDate
	}
	_, err = strconv.Atoi(in.DateEnd)
	if err != nil {
		goto WrongEndDate
	}

	resBillingData.BillingType = in.BillingType
	resBillingData.GroupID = in.GroupID

	if len(*groupID) == 0 {
		resGetGroupList, err := client.RC.GetGroupList()
		if err != nil {
			goto ERROR
		}

		for _, group := range resGetGroupList.Group {
			if group.Id == 1 {
				continue
			}
			groupIDAll = append(groupIDAll, group.Id)
		}

		groupID = &groupIDAll
	}

	data, err = billing.DriverBilling.ReadBillingData(
		groupID, in.DateStart, in.DateEnd,
		in.BillingType, int(in.Row), int(in.Page))
	if data != nil {
		resBillingData.Result, _ = json.Marshal(*data)
	} else {
		resBillingData.Result = []byte{}
	}
	if err != nil {
		goto ERROR
	}

	count, err = dao.GetBillCount(groupID, in.DateStart, in.DateEnd,
		in.BillingType)
	if err != nil {
		goto ERROR
	}
	resBillingData.TotalNum = int32(count)

	goto OUT
ERROR:
	resBillingData.HccErrorStack = errconv.HccStackToGrpc(
		errors.NewHccErrorStack(
			errors.NewHccError(errors.PianoInternalOperationFail, err.Error())))

	goto OUT
WrongStartDate:
	resBillingData.HccErrorStack = errconv.HccStackToGrpc(
		errors.NewHccErrorStack(
			errors.NewHccError(errors.PianoGrpcArgumentError, "-> Wrong start date")))

	goto OUT
WrongEndDate:
	resBillingData.HccErrorStack = errconv.HccStackToGrpc(
		errors.NewHccErrorStack(
			errors.NewHccError(errors.PianoGrpcArgumentError, "-> Wrong end date")))

	goto OUT
OUT:
	return &resBillingData, nil
}

func (s *pianoServer) GetBillingDetail(_ context.Context, in *pb.ReqBillingData) (*pb.ResBillingData, error) {
	var resBillingDetail = pb.ResBillingData{
		BillingType:   "UNDEFINED",
		GroupID:       []int64{},
		Result:        nil,
		HccErrorStack: nil,
	}

	if len(in.GroupID) > 1 {
		resBillingDetail.HccErrorStack = errconv.HccStackToGrpc(
			errors.NewHccErrorStack(
				errors.NewHccError(errors.PianoGrpcArgumentError, "-> Too many Group ID")))

	} else {
		switch in.BillingType {
		case "YEARLY":
			fallthrough
		case "MONTHLY":
			fallthrough
		case "DAILY":
			resBillingDetail.BillingType = in.BillingType
			resBillingDetail.GroupID = in.GroupID

			data, err := billing.DriverBilling.ReadBillingDetail(in.GroupID[0], in.DateStart, in.BillingType)
			if data != nil {
				resBillingDetail.Result, _ = json.Marshal(*data)
			} else {
				resBillingDetail.Result = []byte{}
			}
			if err != nil {
				resBillingDetail.HccErrorStack = errconv.HccStackToGrpc(
					errors.NewHccErrorStack(
						errors.NewHccError(errors.PianoInternalOperationFail, err.Error())))
			}
		default:
			resBillingDetail.HccErrorStack = errconv.HccStackToGrpc(
				errors.NewHccErrorStack(
					errors.NewHccError(errors.PianoGrpcArgumentError, "-> Unsupport BillingType")))
		}

	}

	return &resBillingDetail, nil
}

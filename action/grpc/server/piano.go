package server

import (
	"context"
	"encoding/json"
	"strconv"

	"hcc/piano/action/grpc/errconv"
	"hcc/piano/driver/billing"
	"hcc/piano/driver/influxdb"

	errors "innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

type pianoServer struct {
	pb.UnimplementedPianoServer
}

func (s *pianoServer) Telegraph(ctx context.Context, in *pb.ReqMetricInfo) (*pb.ResMonitoringData, error) {
	series := influxdb.GetInfluxData(in)

	return series, nil
}

func (s *pianoServer) GetBillingData(ctx context.Context, in *pb.ReqBillingData) (*pb.ResBillingData, error) {
	var resBillingData = pb.ResBillingData{
		BillingType:   "UNDEFINED",
		GroupID:       []int32{},
		Result:        nil,
		HccErrorStack: nil,
	}

	switch in.BillingType {
	case "YEARLY":
		in.DateStart = in.DateStart / 10000 * 10000
		in.DateEnd = in.DateEnd / 10000 * 10000
		fallthrough
	case "MONTHLY":
		in.DateStart = in.DateStart / 100 * 100
		in.DateEnd = in.DateEnd / 100 * 100
		fallthrough
	case "DAILY":
		resBillingData.BillingType = in.BillingType
		resBillingData.GroupID = in.GroupID

		data, err := billing.DriverBilling.ReadBillingData(
			&(in.GroupID), strconv.Itoa(int(in.DateStart)), strconv.Itoa(int(in.DateEnd)),
			in.BillingType, int(in.Row), int(in.Page))
		resBillingData.Result, _ = json.Marshal(*data)
		if err != nil {
			resBillingData.HccErrorStack = errconv.HccStackToGrpc(
				errors.NewHccErrorStack(
					errors.NewHccError(errors.PianoInternalOperationFail, err.Error())))
		}
	default:
		resBillingData.HccErrorStack = errconv.HccStackToGrpc(
			errors.NewHccErrorStack(
				errors.NewHccError(errors.PianoGrpcArgumentError, "-> Unsupport BillingType")))
	}

	return &resBillingData, nil
}

func (s *pianoServer) GetBillingDetail(ctx context.Context, in *pb.ReqBillingData) (*pb.ResBillingData, error) {
	var resBillingDetail = pb.ResBillingData{
		BillingType:   "UNDEFINED",
		GroupID:       []int32{},
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
			in.DateStart = in.DateStart / 10000 * 10000
			in.DateEnd = in.DateEnd / 10000 * 10000
			fallthrough
		case "MONTHLY":
			in.DateStart = in.DateStart / 100 * 100
			in.DateEnd = in.DateEnd / 100 * 100
			fallthrough
		case "DAILY":
			resBillingDetail.BillingType = in.BillingType
			resBillingDetail.GroupID = in.GroupID

			data, err := billing.DriverBilling.ReadBillingDetail(in.GroupID[0], strconv.Itoa(int(in.DateStart)), in.BillingType)
			resBillingDetail.Result, _ = json.Marshal(*data)
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

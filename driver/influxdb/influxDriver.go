package influxdb

import (
	"encoding/json"

	"hcc/piano/action/grpc/errconv"
	"hcc/piano/model"

	"innogrid.com/hcloud-classic/hcc_errors"
	"innogrid.com/hcloud-classic/pb"
)

// GetInfluxData - cgs
func GetInfluxData(in *pb.ReqMetricInfo) *pb.ResMonitoringData {
	var resMonitoringData pb.ResMonitoringData
	var monitoringData pb.MonitoringData
	var metricInfo model.MetricInfo

	if in.GetMetricInfo() == nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.PianoGrpcArgumentError, "metricInfo is nil"))
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(errStack)

		resMonitoringData.MonitoringData = &monitoringData

		return &resMonitoringData
	}

	if jsonInfo, e := json.Marshal(in.GetMetricInfo()); e != nil {

	} else {
		if e = json.Unmarshal(jsonInfo, &metricInfo); e != nil {

		}
	}
	resData, err := Influx.ReadMetric(metricInfo)
	if err != nil {
		errStack := hcc_errors.NewHccErrorStack(hcc_errors.NewHccError(hcc_errors.PianoInfluxDBReadMetricError, err.Error()))
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(errStack)

		monitoringData.Uuid = ""
		resMonitoringData.MonitoringData = &monitoringData

		return &resMonitoringData
	}

	monitoringData.Uuid = metricInfo.UUID
	monitoringData.Result, err = json.Marshal(resData)

	resMonitoringData.MonitoringData = &monitoringData

	return &resMonitoringData
}

package influxdb

import (
	"encoding/json"

	"hcc/piano/action/grpc/errconv"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/lib/errors"
	"hcc/piano/model"
)

// GetInfluxData - cgs
func GetInfluxData(in *rpcpiano.ReqMetricInfo) *rpcpiano.ResMonitoringData {
	var resMonitoringData rpcpiano.ResMonitoringData
	var monitoringData rpcpiano.MonitoringData
	var metricInfo model.MetricInfo

	if in.GetMetricInfo() == nil {
		errStack := errors.ReturnHccError(errors.PianoGrpcArgumentError, "metricInfo is nil")
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)

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
		errStack := errors.ReturnHccError(errors.PianoInfluxDBReadMetricError, err.Error())
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)

		monitoringData.Uuid = ""
		resMonitoringData.MonitoringData = &monitoringData

		return &resMonitoringData
	}

	monitoringData.Uuid = metricInfo.UUID
	monitoringData.Result, err = json.Marshal(resData)

	resMonitoringData.MonitoringData = &monitoringData

	return &resMonitoringData
}

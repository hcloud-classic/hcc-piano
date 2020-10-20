package influxdb

import (
	"encoding/json"
	"fmt"
	"strings"

	"hcc/piano/action/grpc/errconv"
	msgType "hcc/piano/action/grpc/pb/rpcmsgType"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/lib/errors"
	"hcc/piano/lib/logger"

	"github.com/influxdata/influxdb1-client/models"
)

// GetInfluxData - cgs
func GetInfluxData(in *rpcpiano.ReqMetricInfo) *rpcpiano.ResMonitoringData {
	var resMonitoringData rpcpiano.ResMonitoringData
	var monitoringData rpcpiano.MonitoringData
	var resSeriesList []*rpcpiano.Series

	if in.GetMetricInfo() == nil {
		errStack := errors.ReturnHccError(errors.PianoGrpcArgumentError, "metricInfo is nil")
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)

		monitoringData.Metric = ""
		monitoringData.SubMetric = ""
		resSeriesList = append(resSeriesList, &rpcpiano.Series{})
		monitoringData.Series = resSeriesList
		resMonitoringData.MonitoringData = &monitoringData

		return &resMonitoringData
	}

	metricInfo := in.GetMetricInfo()
	metric := metricInfo.Metric
	subMetric := strings.ReplaceAll(metricInfo.SubMetric, " ", "")
	period := metricInfo.Period
	aggregateType := metricInfo.AggregateType
	duration := metricInfo.Duration
	uuid := metricInfo.Uuid

	resData, err := Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		errStack := errors.ReturnHccError(errors.PianoInfluxDBReadMetricError, err.Error())
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)

		monitoringData.Metric = ""
		monitoringData.SubMetric = ""
		resSeriesList = append(resSeriesList, &rpcpiano.Series{})
		monitoringData.Series = resSeriesList
		resMonitoringData.MonitoringData = &monitoringData

		return &resMonitoringData
	}
	queryResult := resData.(models.Row)

	var seriesList []*msgType.Series
	for _, list := range queryResult.Values {
		var series msgType.Series
		for _, value := range list {
			val, _ := value.(json.Number).Float64()
			series.Values = append(series.Values, val)
		}
		seriesList = append(seriesList, &series)
	}

	logger.Logger.Println("queryResult [" + metric + "] (time," + subMetric + ")" + ": " + fmt.Sprintf("%v", queryResult.Values))

	monitoringData.Metric = metric
	monitoringData.SubMetric = subMetric
	monitoringData.UUID = uuid
	monitoringData.Series = seriesList

	resMonitoringData.MonitoringData = &monitoringData

	return &resMonitoringData
}

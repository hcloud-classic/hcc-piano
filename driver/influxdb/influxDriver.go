package influxdb

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	"hcc/piano/action/grpc/errconv"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/lib/errors"
	"hcc/piano/lib/logger"
	"strconv"
)

// GetInfluxData - cgs
func GetInfluxData(in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {
	var resMonitoringData rpcpiano.ResMonitoringData

	if in.GetMetricInfo() == nil {
		errStack := errors.ReturnHccError(errors.PianoGrpcArgumentError, "metricInfo is nil")
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)
		return &resMonitoringData, nil
	}

	var monitoringData rpcpiano.MonitoringData
	var seriesList []rpcpiano.Series
	var resSeriesList []*rpcpiano.Series

	metricInfo := in.GetMetricInfo()
	metric := metricInfo.Metric
	subMetric := metricInfo.SubMetric
	period := metricInfo.Period
	aggregateType := metricInfo.AggregateType
	duration := metricInfo.Duration
	if metric == "net" {
		durationInt, _ := strconv.Atoi(duration[:len(duration)-1])
		duration = strconv.Itoa(durationInt + 1) + duration[len(duration)-1:]
	}
	uuid := metricInfo.Uuid

	queryResult, err := Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		errStack := errors.ReturnHccError(errors.PianoInfluxDBReadMetricError, err.Error())
		resMonitoringData.HccErrorStack = errconv.HccStackToGrpc(&errStack)
		return &resMonitoringData, nil
	}
	logger.Logger.Println("queryResult (" +metric + ", " + subMetric + ")" +  ": " + fmt.Sprintf("%v", queryResult))

	dataLength := len(queryResult.(models.Row).Values)
	if metric == "net" {
		dataLength--
	}
	logger.Logger.Println("data : ", queryResult.(models.Row).Values)

	for i := 0; i < dataLength; i++ {
		var time int64
		var value int64

		time = int64(i)

		switch metric {
		case "cpu":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			value = int64(valueFloat * 100)
		case "mem", "disk":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			value = int64(valueFloat)
		case "diskio":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			value = int64(valueFloat / 1024 / 1024) // MB
		case "net":
			valueBeforeStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueBeforeFloat, _ := strconv.ParseFloat(valueBeforeStr, 64)
			valueBefore := int64(valueBeforeFloat / 1024 / 1024) // MB

			valueAfterStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i + 1][1])
			valueAfterFloat, _ := strconv.ParseFloat(valueAfterStr, 64)
			valueAfter := int64(valueAfterFloat / 1024 / 1024) // MB

			value = valueAfter - valueBefore
		default:
			continue
		}

		seriesList = append(seriesList, rpcpiano.Series{
			Time:  time,
			Value: value,
		})
	}

	for i := range seriesList {
		resSeriesList = append(resSeriesList, &seriesList[i])
	}

	monitoringData.Metric = metric
	monitoringData.SubMetric = subMetric
	monitoringData.Series = resSeriesList
	resMonitoringData.MonitoringData = &monitoringData

	return &resMonitoringData, nil
}

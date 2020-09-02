package influxdb

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
	"strconv"
)

// GetInfluxData - cgs
func GetInfluxData(in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {

	if in.GetMetricInfo() == nil {
		return nil, errors.New("metricInfo is nil")
	}

	var resMonitoringData rpcpiano.ResMonitoringData
	var monitoringData rpcpiano.MonitoringData
	var seriesEntry rpcpiano.Series
	var seriesList []rpcpiano.Series
	var seriesResponse []*rpcpiano.Series

	metricInfo := in.GetMetricInfo()
	metric := metricInfo.Metric
	subMetric := metricInfo.SubMetric
	period := metricInfo.Period
	aggregateType := metricInfo.AggregateType
	duration := metricInfo.Duration
	uuid := metricInfo.Uuid

	//if !metricOk || !subMetricOk || !periodOk || !aggregateTypeOk || !durationOk || !uuidOk {
	//	return nil, nil
	//}

	var telegraf model.Telegraf
	var series []model.Series
	var s model.Series

	var err error

	queryResult, err := Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		return nil, nil
	}
	logger.Logger.Println("queryResult : " + fmt.Sprintf("%v", queryResult))
	telegraf.UUID = fmt.Sprintf("%v", queryResult.(models.Row).Tags["host"])

	dataLength := len(queryResult.(models.Row).Values)
	logger.Logger.Println("data : ", queryResult.(models.Row).Values)

	if queryResult == nil {
		for j := 0; j < 11; j++ {
			s.Time = j
			s.Value = 0
			series = append(series, s)
		}
	} else if dataLength < 11 {
		for j := 0; j < 11-dataLength; j++ {
			s.Time = j
			s.Value = 0
			series = append(series, s)
		}
	}

	for i := 0; i < dataLength; i++ {
		s.Time = 11 - dataLength + i

		seriesEntry.Time = int64(11 - dataLength + i)

		switch metric {
		case "cpu":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 100)
			seriesEntry.Value = int64(valueFloat * 100)
			break
		case "mem":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			seriesEntry.Value = int64(valueFloat * 1)
			break
		case "net":
		case "disk":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			seriesEntry.Value = int64(valueFloat * 1)
			break
			//case "net":
			//	valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			//	valueInt, _ := strconv.ParseInt(valueStr, 10, 64)
			//	s.Value = int(valueInt * 1)
			//	break

			//case "process":
			//	valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			//	valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			//	s.Value = int(valueFloat * 1)
			//	seriesEntry.Value = int64(valueFloat * 1)
			//	break

		}

		series = append(series, s)
		seriesList = append(seriesList, seriesEntry)
	}

	telegraf.Series = series
	telegraf.Metric = metric
	telegraf.SubMetric = subMetric

	for j := 0; j < dataLength; j++ {
		seriesResponse = append(seriesResponse, &seriesList[j])
	}
	monitoringData.Metric = metric
	monitoringData.SubMetric = subMetric
	monitoringData.Series = seriesResponse
	resMonitoringData.MonitoringData = &monitoringData

	return &resMonitoringData, nil
}

package driver

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	"hcc/piano/action/grpc/rpcpiano"
	"hcc/piano/lib/influxdb"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
	"strconv"
)

// GetInfluxData - cgs
//func GetInfluxData(args map[string]interface{}) (interface{}, error) {
func GetInfluxData(in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {

	if in.GetMetricInfo() == nil {
		return nil, errors.New("metricInfo is nil")
	}

	var monitoringData *rpcpiano.ResMonitoringData
	metricInfo := in.GetMetricInfo()

	metric := metricInfo.Metric
	subMetric := metricInfo.SubMetric
	period := metricInfo.Period
	aggregateType := metricInfo.AggregateType
	duration := metricInfo.Duration
	uuid := metricInfo.Uuid

	//metric, metricOk := args["metric"].(string)
	//subMetric, subMetricOk := args["subMetric"].(string)
	//period, periodOk := args["period"].(string)
	//aggregateType, aggregateTypeOk := args["aggregateType"].(string)
	//duration, durationOk := args["duration"].(string)
	//uuid, uuidOk := args["uuid"].(string)

	//if !metricOk || !subMetricOk || !periodOk || !aggregateTypeOk || !durationOk || !uuidOk {
	//	return nil, nil
	//}

	var telegraf model.Telegraf
	var series []model.Series
	var s model.Series
	var err error

	queryResult, err := influxdb.Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
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

		switch metric {
		case "cpu":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 100)
			break
		case "mem":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			break
		case "net":
		case "disk":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			break
			//case "net":
			//	valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			//	valueInt, _ := strconv.ParseInt(valueStr, 10, 64)
			//	s.Value = int(valueInt * 1)
			//	break
		}

		series = append(series, s)
	}

	telegraf.Series = series
	telegraf.Metric = metric
	telegraf.SubMetric = subMetric

	return monitoringData, nil
}

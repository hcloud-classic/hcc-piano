package driver

import (
	"fmt"
	"github.com/influxdata/influxdb1-client/models"
	"strconv"

	//	"time"
	//	"hcc/piano/lib/logger"
	"hcc/piano/lib/influxdb"
	"hcc/piano/model"
)

// GetInfluxData - cgs
func GetInfluxData(args map[string]interface{}) (interface{}, error) {

	metric, metricOk := args["metric"].(string)
	subMetric, subMetricOk := args["subMetric"].(string)
	period, periodOk := args["period"].(string)
	aggregateType, aggregateTypeOk := args["aggregateType"].(string)
	duration, durationOk := args["duration"].(string)
	uuid, uuidOk := args["uuid"].(string)

	if !metricOk || !subMetricOk || !periodOk || !aggregateTypeOk || !durationOk || !uuidOk {
		return nil, nil
	}

	var telegraf model.Telegraf
	var series []model.Series
	var s model.Series

	//queryResult, err := influxdb.Influx.ReadMetric("cpu", "s", "avg", "1m", "hcc-ubuntu")
	queryResult, err := influxdb.Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		return nil, nil
	}

	//value := fmt.Sprintf("%v", queryResult.(models.Row).Values)
	//value1 := queryResult.(models.Row).Values[0][0]
	//logger.Logger.Println("InfluxDB queryResult : " + value)
	//logger.Logger.Println("value1 : " + fmt.Sprintf("%v", value1))
	//logger.Logger.Println("value length : " + fmt.Sprintf("%v", len(queryResult.(models.Row).Values)))

	telegraf.UUID = fmt.Sprintf("%v", queryResult.(models.Row).Tags["host"])

	for i := 0; i < len(queryResult.(models.Row).Values); i++ {

		//s.Time = fmt.Sprintf("%v", queryResult.(models.Row).Values[i][0])
		s.Time = i

		valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
		valueFloat, _ := strconv.ParseFloat(valueStr, 64)
		s.Value = int(valueFloat * 100)
		series = append(series, s)
	}
	telegraf.Series = series
	telegraf.SubMetric = subMetric

	return telegraf, nil
}

package influxdb

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	influxBuilder "github.com/Scalingo/go-utils/influx"
	influxdbClient "github.com/influxdata/influxdb1-client/v2"
	"innogrid.com/hcloud-classic/hcc_errors"

	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
)

// HostInfo : Contain infos of InfluxDB's host
type HostInfo struct {
	URL      string
	Username string
	Password string
}

// InfluxInfo : Contain infos of InfluxDB
type InfluxInfo struct {
	HostInfo HostInfo
	Database string
	Clients  []influxdbClient.Client
}

// Influx : Exported variable to get infos of InfluxDB
var Influx InfluxInfo

// Init : Initialize InfluxDB connection
func Init() *hcc_errors.HccError {
	hostInfo := HostInfo{
		URL:      "http://" + config.Influxdb.Address + ":" + strconv.FormatInt(config.Influxdb.Port, 10),
		Username: config.Influxdb.ID,
		Password: config.Influxdb.Password,
	}
	Influx = InfluxInfo{HostInfo: hostInfo, Database: config.Influxdb.Db}
	err := Influx.InitInfluxDB()
	if err != nil {
		return hcc_errors.NewHccError(
			hcc_errors.PianoInternalInitFail, "influxdb.Init(): "+err.Error())
	}
	return nil
}

// InitInfluxDB : Check if InfluxDB is available
func (s *InfluxInfo) InitInfluxDB() error {
	client, err := influxdbClient.NewHTTPClient(influxdbClient.HTTPConfig{
		Addr:     s.HostInfo.URL,
		Username: s.HostInfo.Username,
		Password: s.HostInfo.Password,
	})
	if err != nil {
		logger.Logger.Println("NewHTTPClient error")
		return err
	}
	if _, _, err := client.Ping(time.Millisecond * 100); err != nil {
		logger.Logger.Println("Ping error")
		return err
	}

	s.Clients = append(s.Clients, client)

	return nil
}

// ReadMetric : Read metrics from InfluxDB
func (s *InfluxInfo) ReadMetric(metricInfo model.MetricInfo) (interface{}, error) {
	influx := s.Clients[0]

	queryString, err := s.GenerateQuery(metricInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println("ReadMetric query : " + queryString)

	query := influxdbClient.NewQuery(queryString, s.Database, metricInfo.Period)
	res, _ := influx.Query(query)

	if res.Err != "" {
		logger.Logger.Println("ReadMetric(): res.Err")
		return nil, errors.New(res.Err)
	}

	return res.Results, nil
}

// GenerateQuery : Generate the query for InfluxDB
func (s *InfluxInfo) GenerateQuery(metricInfo model.MetricInfo) (string, error) {

	// InfluxDB 쿼리 생성
	var subMetricList = strings.Split(metricInfo.SubMetric, ",")
	var aggregateFnList = strings.Split(metricInfo.AggregateFn, ",")

	fmt.Println(metricInfo.AggregateFn)

	query := influxBuilder.NewQuery().On(metricInfo.Metric)

	for index, sub := range subMetricList {
		query = query.Field(sub, aggregateFnList[index])
	}

	query = query.Where("host", influxBuilder.Equal, influxBuilder.String(metricInfo.UUID))
	if metricInfo.Metric == "cpu" {
		query = query.And("cpu", influxBuilder.Equal, influxBuilder.String("cpu-total"))
	}
	if metricInfo.Time != "" {
		query = query.And("time", influxBuilder.MoreThan, metricInfo.Time)
	}
	if metricInfo.GroupBy != "" {
		query = query.GroupByTag(strings.Split(metricInfo.GroupBy, ",")...)
	}

	limit, _ := strconv.Atoi(metricInfo.Limit)
	query = query.OrderByTime("DESC").Limit(limit)
	queryString := query.Build()

	return queryString, nil
}

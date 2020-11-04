package influxdb

import (
	"errors"
	"fmt"
	influxBuilder "github.com/Scalingo/go-utils/influx"
	influxdbClient "github.com/influxdata/influxdb1-client/v2"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"strconv"
	"strings"
	"time"
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
func Init() error {
	hostInfo := HostInfo{
		URL:      "http://" + config.Influxdb.Address + ":" + strconv.FormatInt(config.Influxdb.Port, 10),
		Username: config.Influxdb.ID,
		Password: config.Influxdb.Password,
	}
	Influx = InfluxInfo{HostInfo: hostInfo, Database: config.Influxdb.Db}
	err := Influx.InitInfluxDB()
	if err != nil {
		return err
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
func (s *InfluxInfo) ReadMetric(metric string, subMetric string, period string, aggregateType string, duration string, uuid string) (interface{}, error) {
	logger.Logger.Println("ReadMetric")
	influx := s.Clients[0]

	queryString, err := s.GenerateQuery(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		return nil, err
	}
	logger.Logger.Println("ReadMetric query : " + queryString)
	fmt.Println("ReadMetric query : " + queryString)

	query := influxdbClient.NewQuery(queryString, s.Database, period)
	res, _ := influx.Query(query)
	if res.Err != "" {
		return nil, errors.New(res.Err)
	}

	if len(res.Results) != 0 {
		if len(res.Results[0].Series) != 0 {
			logger.Logger.Println("ReadMetric - series")
			return res.Results[0].Series[0], nil
		}
	}

	logger.Logger.Println("ReadMetric(): failed to get metric")
	return nil, errors.New("failed to get metric")
}

// GenerateQuery : Generate the query for InfluxDB
func (s *InfluxInfo) GenerateQuery(metric string, subMetric string, period string, aggregateType string, duration string, uuid string) (string, error) {

	// InfluxDB 쿼리 생성
	var subMetricList = strings.Split(subMetric, ",")

	query := influxBuilder.NewQuery().On(metric)

	for _, sub := range subMetricList {
		if metric == "net" {
			query = query.Field(sub, "difference")
		} else {
			query = query.Field(sub, "")
		}
	}

	query = query.Where("host", influxBuilder.Equal, influxBuilder.String(uuid))
	if metric == "cpu" {
		query = query.And("cpu", influxBuilder.Equal, influxBuilder.String("cpu-total"))
	}
	if metric == "net" {
		query = query.And("interface", influxBuilder.Equal, influxBuilder.String("eth0"))
	}
	if metric == "disk" {
		query = query.And("path", influxBuilder.Equal, influxBuilder.String("/"))
	}
	if aggregateType != "" {
		query = query.And("time", influxBuilder.MoreThan, aggregateType)
	}

	limits, _ := strconv.Atoi(duration)
	query = query.OrderByTime("DESC").Limit(limits)
	queryString := query.Build()

	return queryString, nil
}

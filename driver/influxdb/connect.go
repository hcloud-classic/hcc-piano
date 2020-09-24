package influxdb

import (
	"errors"
	"fmt"
	influxBuilder "github.com/Scalingo/go-utils/influx"
	influxdbClient "github.com/influxdata/influxdb1-client/v2"
	"hcc/piano/lib/config"
	"hcc/piano/lib/logger"
	"strconv"
	"time"
)

// HostInfo - cgs
type HostInfo struct {
	URL      string
	Username string
	Password string
}

// InfluxInfo - cgs
type InfluxInfo struct {
	HostInfo HostInfo
	Database string
	Clients  []influxdbClient.Client
}

// Influx - cgs
var Influx InfluxInfo

// Prepare - cgs
func Init() error {
	hostInfo := HostInfo{
		URL:      "http://" + config.Influxdb.Address + ":" + strconv.FormatInt(config.Influxdb.Port, 10),
		Username: config.Influxdb.Id,
		Password: config.Influxdb.Password,
	}
	Influx = InfluxInfo{HostInfo: hostInfo, Database: config.Influxdb.Db}
	err := Influx.InitInfluxDB()
	if err != nil {
		return err
	}
	return nil
}

// InitInfluxDB - cgs
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

// ReadMetric - cgs
func (s *InfluxInfo) ReadMetric(metric string, subMetric string, period string, aggregateType string, duration string, uuid string) (interface{}, error) {
	logger.Logger.Println("ReadMetric")
	influx := s.Clients[0]

	queryString, err := s.GenerateQuery(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		return nil, err
	}
	logger.Logger.Println("ReadMetric query : " + queryString)
	fmt.Println("ReadMetric query : " + queryString)

	query := influxdbClient.NewQuery(queryString, s.Database, "")
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

	return nil, errors.New("failed to get metric")
}

// GenerateQuery - cgs
func (s *InfluxInfo) GenerateQuery(metric string, subMetric string, period string, aggregateType string, duration string, uuid string) (string, error) {

	// 통계 기준 설정
	switch aggregateType {
	case "median":
	case "middle":
		aggregateType = "median"
	case "mode":
	case "frequency":
		aggregateType = "mode"
	case "average":
		aggregateType = "mean"
	default:
		aggregateType = "mean"
	}

	// 시간 범위 설정
	timeDuration := fmt.Sprintf("now() - %s", duration)

	// 시간 단위 설정
	var timeCriteria time.Duration
	switch period {
	case "s":
		timeCriteria = time.Second
	case "m":
		timeCriteria = time.Minute
	case "h":
		timeCriteria = time.Hour
	case "d":
		timeCriteria = time.Hour * 24
	}

	// InfluxDB 쿼리 생성
	var query influxBuilder.Query

	switch metric {
	case "cpu":
		switch subMetric {
		case "usage_system":
			query = influxBuilder.NewQuery().On(metric).
				Field("usage_system", aggregateType)
			break
		case "usage_user":
			query = influxBuilder.NewQuery().On(metric).
				Field("usage_system", aggregateType)
			break
		}
		break
	case "mem":
		switch subMetric {
		case "used_percent":
			query = influxBuilder.NewQuery().On(metric).
				Field("used_percent", aggregateType)
			break
		case "swap_total":
			query = influxBuilder.NewQuery().On(metric).
				Field("swap_total", aggregateType)
			break
		}
		break
	case "disk":
		switch subMetric {
		case "used_percent":
			query = influxBuilder.NewQuery().On(metric).
				Field("used_percent", aggregateType)
			break
		}
		break
	case "diskio":
		switch subMetric {
		case "read_bytes":
			query = influxBuilder.NewQuery().On(metric).
				Field("read_bytes", aggregateType)
			break
		}
		break
	case "net":
		switch subMetric {
		case "bytes_recv":
			query = influxBuilder.NewQuery().On(metric).
				Field("bytes_recv", aggregateType)
			break
		case "bytes_sent":
			query = influxBuilder.NewQuery().On(metric).
				Field("bytes_sent", aggregateType)
			break
		}
		break
	default:
		return "", errors.New("not found metric")
	}

	hostname := uuid
	query = query.Where("time", influxBuilder.MoreThan, timeDuration).
		And("\"host\"", influxBuilder.Equal, "'"+hostname+"'").
		GroupByTime(timeCriteria).
		GroupByTag("\"host\"").
		Fill(influxBuilder.None).
		OrderByTime("ASC")

	queryString := query.Build()

	return queryString, nil
}

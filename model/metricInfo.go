package model

type MetricInfo struct {
	UUID        string `json:"uuid"`
	Metric      string `json:"metric"`
	SubMetric   string `json:"subMetric"`
	Period      string `json:"period"`
	AggregateFn string `json:"aggregateFn"`
	Duration    string `json:"duration"`
	Time        string `json:"time"`
	GroupBy     string `json:"groupBy"`
	OrderBy     string `json:"orderBy"`
	Limit       string `json:"limit"`
}

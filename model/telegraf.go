package model

type Telegraf struct {
	Metric    string   `json:"metric"`
	SubMetric string   `json:"subMetric"`
	UUID      string   `json:"id"`
	Series    []Series `json:"data"`
}

type Series struct {
	Time  int `json:"x"`
	Value int `json:"y"`
}

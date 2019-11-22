package model

// Telegraf - cgs
type Telegraf struct {
	SubMetric string   `json:"subMetric"`
	UUID      string   `json:"id"`
	Series    []Series `json:"data"`
}

// Series - cgs
type Series struct {
	//Time  string `json:"x"`
	Time  int `json:"x"`
	Value int `json:"y"`
}

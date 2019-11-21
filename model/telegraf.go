package model

// Telegraf - cgs
type Telegraf struct {
	UUID   string   `json:"uuid"`
	Series []Series `json:"series"`
}

// Series - cgs
type Series struct {
	Time  string `json:"time"`
	Value string `json:"value"`
}

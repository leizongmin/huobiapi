package types

type KLineData struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Data   []struct {
		ID     int64   `json:"id"`
		Amount float64 `json:"amount"`
		Count  int64   `json:"count"`
		Open   float64 `json:"open"`
		Close  float64 `json:"close"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
		Vol    float64 `json:"vol"`
	} `json:"data"`
}

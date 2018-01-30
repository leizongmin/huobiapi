package data_type

import "encoding/json"

type Kline struct {
	Ch   string    `json:"ch"`
	Ts   uint      `json:"ts"`
	Tick KlineTick `json:"tick"`
}

type KlineTick struct {
	ID     uint    `json:"id"`
	Amount float64 `json:"amount"`
	Count  uint    `json:"count"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Vol    float64 `json:"vol"`
}

func DecodeKline(raw []byte) (*Kline, error) {
	var ret = &Kline{}
	if err := json.Unmarshal(raw, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

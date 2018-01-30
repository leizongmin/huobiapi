package data_type

import "encoding/json"

type Trade struct {
	Ch   string `json:"ch"`
	Ts   uint   `json:"ts"`
	Tick struct {
		ID   uint        `json:"id"`
		Ts   uint        `json:"ts"`
		Data []TradeItem `json:"data"`
	} `json:"tick"`
}

type TradeItem struct {
	Ts        uint    `json:"ts"`
	ID        uint    `json:"id"`
	Direction string  `json:"direction"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
}

func DecodeTrade(raw []byte) (*Trade, error) {
	var ret = &Trade{}
	if err := json.Unmarshal(raw, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

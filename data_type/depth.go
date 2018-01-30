package data_type

import "encoding/json"

type Depth struct {
	Ch   string `json:"ch"`
	Ts   uint   `json:"ts"`
	Tick struct {
		Bids [][]float64 `json:"bids"`
		Asks [][]float64 `json:"asks"`
	} `json:"tick"`
}

func DecodeDepth(raw []byte) (*Depth, error) {
	var ret = &Depth{}
	if err := json.Unmarshal(raw, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

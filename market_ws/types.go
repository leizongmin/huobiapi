package market_ws

type pongData struct {
	Pong int `json:"pong"`
}

type pingData struct {
	Ping int `json:"ping"`
}

type subData struct {
	Sub string `json:"sub"`
	ID string `json:"id"`
}

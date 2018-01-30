package data_type

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeKLine(t *testing.T) {
	str := `{"ch":"market.eosusdt.kline.1min","tick":{"amount":280.483600000000000000,"close":14.290000000000000000,"count":9,"high":14.300000000000000000,"id":1516870800,"low":14.290000000000000000,"open":14.290000000000000000,"vol":4010.253191000000000000000000000000000000},"ts":1516870810953}`
	data, err := DecodeKline([]byte(str))
	assert.NoError(t, err)
	fmt.Println(data)
	assert.Equal(t, "market.eosusdt.kline.1min", data.Ch)
	assert.Equal(t, uint(1516870810953), data.Ts)
	assert.Equal(t, 280.4836, data.Tick.Amount)
	assert.Equal(t, 14.29, data.Tick.Close)
	assert.Equal(t, 14.29, data.Tick.Open)
	assert.Equal(t, uint(9), data.Tick.Count)
	assert.Equal(t, 14.3, data.Tick.High)
	assert.Equal(t, 14.29, data.Tick.Low)
	assert.Equal(t, uint(1516870800), data.Tick.ID)
	assert.Equal(t, 4010.253191, data.Tick.Vol)
}

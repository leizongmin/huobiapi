package data_type

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeTrade(t *testing.T) {
	str := `{"ch":"market.eosusdt.trade.detail","tick":{"data":[{"amount":8.333400000000000000,"direction":"sell","id":17592232489498,"price":14.290000000000000000,"ts":1516870810708}],"id":1557850942,"ts":1516870810708},"ts":1516870811033}`
	data, err := DecodeTrade([]byte(str))
	assert.NoError(t, err)
	fmt.Println(data)
	assert.Equal(t, "market.eosusdt.trade.detail", data.Ch)
	assert.Equal(t, uint(1516870811033), data.Ts)
	assert.Equal(t, uint(1516870810708), data.Tick.Ts)
	assert.Equal(t, uint(1557850942), data.Tick.ID)
	assert.Equal(t, uint(17592232489498), data.Tick.Data[0].ID)
	assert.Equal(t, uint(1516870810708), data.Tick.Data[0].Ts)
	assert.Equal(t, "sell", data.Tick.Data[0].Direction)
	assert.Equal(t, 8.3334, data.Tick.Data[0].Amount)
	assert.Equal(t, 14.29, data.Tick.Data[0].Price)
}

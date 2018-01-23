package market_ws

import (
	"fmt"
	"testing"
)

func TestMarket(t *testing.T) {
	m, err := NewMarket()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(m)
	}
}

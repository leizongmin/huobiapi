package market_ws

import (
	"testing"
	"fmt"
)

func TestMarket(t *testing.T) {
	m, err := NewMarket()
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(m)
	}
}

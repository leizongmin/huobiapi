package main

import (
	"fmt"

	"github.com/leizongmin/huobiapi"
)

func main() {
	client, err := huobiapi.NewClient("key id", "key secret")
	if err != nil {
		panic(err)
	}
	ret, err := client.Request("GET", "/market/history/trade", huobiapi.ParamsData{
		"symbol": "eosusdt",
		"size":   "10",
	})
	data, err := ret.Get("data").Array()
	if err != nil {
		panic(err)
	}
	for _, v := range data {
		fmt.Println(v)
	}
}

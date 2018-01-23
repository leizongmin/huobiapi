package market

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leizongmin/huobiapi/client"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetKLine(t *testing.T) {
	options := client.ClientOptions{
		AccessKeyId:      "",
		AccessKeySecret:  "",
		SignatureMethod:  "",
		SignatureVersion: "",
	}
	c := NewClient(options)
	data, err := c.GetKLine()
	assert.NoError(t, err)
	fmt.Println(data)
	str, err := json.Marshal(data)
	assert.NoError(t, err)
	fmt.Println(string(str))
}

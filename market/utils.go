package market

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

/// 返回随机字符串
func getRandomString(n uint) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getMapKeys(data map[string]bool) []string {
	var keys []string
	for k, _ := range data {
		keys = append(keys, k)
	}
	return keys
}

func getUinxMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

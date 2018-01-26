package market

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
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

/// 取毫秒时间戳
func getUinxMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

/// 解压gzip的数据
func unGzipData(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

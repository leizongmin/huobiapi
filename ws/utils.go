package ws

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// getRandomString 返回随机字符串
func getRandomString(n uint) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// getUinxMillisecond 取毫秒时间戳
func getUinxMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// unGzipData 解压gzip的数据
func unGzipData(buf []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func getMapKeys(m map[string]string) (keys []string) {
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func sortKeys(keys []string) []string {
	sort.Strings(keys)
	return keys
}

func computeHmac256(data string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/// 拼接query字符串
func encodeQueryString(query map[string]string) string {
	var keys = sortKeys(getMapKeys(query))
	var len = len(keys)
	var lines = make([]string, len)
	for i := 0; i < len; i++ {
		var k = keys[i]
		lines[i] = url.QueryEscape(k) + "=" + url.QueryEscape(query[k])
	}
	return strings.Join(lines, "&")
}

func GenSignature(query map[string]string, accessKeySecret string) string {
	var pre = "GET" + "\n" + "api-aws.huobi.pro" + "\n" + "/ws/v1" + "\n"
	eqs := encodeQueryString(query)
	fmt.Println(eqs)
	eqs = pre + eqs
	return computeHmac256(eqs, accessKeySecret)
}
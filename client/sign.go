package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

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

type Sign struct {
	AccessKeyId      string
	AccessKeySecret  string
	SignatureMethod  string
	SignatureVersion string
}

func NewSign(accessKeyId, accessKeySecret string) *Sign {
	return &Sign{
		AccessKeyId:      accessKeyId,
		AccessKeySecret:  accessKeySecret,
		SignatureMethod:  "HmacSHA256",
		SignatureVersion: "2",
	}
}

func (s *Sign) Get(method, host, path, timestamp string, params map[string]string) (string, error) {
	var str = method + "\n" + host + "\n" + path + "\n"
	params["AccessKeyId"] = s.AccessKeyId
	params["SignatureMethod"] = s.SignatureMethod
	params["SignatureVersion"] = s.SignatureVersion
	params["Timestamp"] = timestamp
	str += encodeQueryString(params)
	return computeHmac256(str, s.AccessKeySecret), nil
}

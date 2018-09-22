package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

func isGetMethod(method string) bool {
	if method == "GET" || method == "HEAD" {
		return true
	}
	return false
}

/// 请求参数
type ParamData = map[string]string

/// 发送原始请求
func SendRequest(sign *Sign, method, scheme, host, path string, data ParamData) (*simplejson.Json, error) {
	var body *bytes.Buffer
	method = strings.ToUpper(method)
	if data == nil {
		data = ParamData{}
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
	// 参与计算签名的参数
	signData := make(map[string]string)
	if isGetMethod(method) {
		// GET 请求所有参数都参与签名计算，POST 请求业务参数不参与签名计算
		for k, v := range data {
			signData[k] = v
		}
	}
	if s, err := sign.Get(method, host, path, timestamp, signData); err != nil {
		return nil, err
	} else {
		signData["Signature"] = s
	}
	path += "?" + encodeQueryString(signData)
	if isGetMethod(method) == false {
		// POST 请求 JSON
		if b, err := json.Marshal(data); err != nil {
			return nil, err
		} else {
			body = bytes.NewBuffer(b)
		}
	}

	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, scheme+"://"+host+path, body)
	} else {
		req, err = http.NewRequest(method, scheme+"://"+host+path, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-cn")
	if isGetMethod(method) {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	json, err := simplejson.NewJson(resBody)
	if err != nil {
		return nil, err
	}
	var status = json.Get("status").MustString()
	if status == "error" {
		return json, fmt.Errorf(json.Get("err-msg").MustString())
	}
	return json, nil
}

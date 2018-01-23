package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func isGetMethod(method string) bool {
	if method == "GET" || method == "HEAD" {
		return true
	}
	return false
}

type ParamData = map[string]string

func SendRequest(sign *Sign, method, host, path string, data ParamData) ([]byte, error) {
	var body *bytes.Buffer
	method = strings.ToUpper(method)
	if data == nil {
		data = ParamData{}
	}

	if s, err := sign.Get(method, host, path, time.Now().Format("2017-05-11T15:19:30"), data); err != nil {
		return nil, err
	} else {
		data["Signature"] = s
	}
	if isGetMethod(method) {
		path += "?" + encodeQueryString(data)
	} else {
		if b, err := json.Marshal(data); err != nil {
			return nil, err
		} else {
			body = bytes.NewBuffer(b)
		}
	}
	fmt.Println(body)
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, "https://"+host+path, body)
	} else {
		req, err = http.NewRequest(method, "https://"+host+path, nil)
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
	return resBody, nil
}

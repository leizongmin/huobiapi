# huobiapi

火币网 API Go 客户端

## 安装

本模块使用 [godep](https://github.com/golang/dep) 作为包管理工具

```bash
dep ensure -v -add github.com/leizongmin/huobiapi
```

## Websocket 版行情数据订阅查询

说明：返回的结果数据使用 [go-simplejson](https://github.com/bitly/go-simplejson) 存储

详细代码参考 [examples](https://github.com/leizongmin/huobiapi/tree/master/examples) 目录

```go
package main

import (
    "fmt"
    "github.com/leizongmin/huobiapi"
)

func main() {
    // 创建客户端实例
    market, err := huobiapi.NewMarket()
    if err != nil {
        panic(err)
    }
    // 订阅主题
    market.Subscribe("market.eosusdt.trade.detail", func(topic string, json *simplejson.Json, raw []byte) {
        // 收到数据更新时回调
        fmt.Println(topic, json, raw)
    })
    // 请求数据
    json, err := market.Request("market.eosusdt.detail")
    if err != nil {
        panic(err)
    } else {
        fmt.Println(json)
    }
    // 进入阻塞等待，这样不会导致进程退出
    market.Loop()
}
```

## RESTful 版行情和交易查询

```go
package main

import (
    "fmt"
    "github.com/leizongmin/huobiapi"
)

func main() {
    // 创建客户端实例
    client, err := huobiapi.NewClient("key id", "key secret")
    if err != nil {
        panic(err)
    }
    // 发送请求
    ret, err := client.Request("GET", "/market/history/trade", huobiapi.ParamsData{
        "symbol": "eosusdt",
        "size": "10",
    })
    data, err := ret.Get("data").Array()
    if err != nil {
        panic(err)
    }
    for _, v := range data {
        fmt.Println(v)
    }
}
```

## License

```text
MIT License

Copyright (c) 2018 Zongmin Lei <leizongmin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

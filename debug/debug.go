package debug

import (
	"log"
	"os"
	"strings"
)

var IsOutputDebug bool = false

func init() {
	if v, ok := os.LookupEnv("HUOBI_DEBUG"); ok {
		v = strings.ToLower(v)
		if v == "1" || v == "true" || v == "yes" || v == "ok" {
			IsOutputDebug = true
		}
	}
}

func Println(a ...interface{}) {
	if IsOutputDebug {
		log.Println(a...)
	}
}

package main

import (
	"fmt"
	"net/http"
	"simple_gin/gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL是：%q\n", r.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q", k, v)
		}
	})
	r.Run(":8081")
}

// git 代理：
//git config --global http.proxy http://127.0.0.1:1080
//git config --global https.proxy http://127.0.0.1:1080
// 取消代理：
//git config --global --unset http.proxy
//git config --global --unset https.proxy

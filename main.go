package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

// 这里实现了Handler接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL是%q\n", req.URL.Path)
	case "/hello":
		fmt.Fprintf(w, "URL是%q\n", req.URL.Path)
	default:
		fmt.Fprintf(w, "404 NOT FOUND : %s", req.URL)
	}
}
func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":8081", engine))
}

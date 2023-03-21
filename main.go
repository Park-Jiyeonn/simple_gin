package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("使用了")
		_, err := fmt.Fprintf(w, "URL = %q", req.URL)
		if err != nil {
			fmt.Println(err.Error())
		}
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "URL = %q", req.URL.Path)
	if err != nil {
		fmt.Println(err.Error())
	}
}

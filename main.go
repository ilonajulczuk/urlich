package main

import (
	"github.com/ilonajulczuk/urlich/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/add", handlers.AddHandler)
	http.HandleFunc("/", handlers.ViewHandler)
	http.ListenAndServe(":9081", nil)
}

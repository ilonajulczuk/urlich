package main

import (
	"github.com/ilonajulczuk/urlich/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.ViewHandler)
	http.HandleFunc("/add/", handlers.AddHandler)
	http.ListenAndServe(":9081", nil)
}

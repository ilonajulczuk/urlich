package main

import (
	"github.com/ilonajulczuk/urlich/handlers"
	"github.com/ilonajulczuk/urlich/pages"
	"net/http"
)

func main() {
	db := pages.NewRedisClient(&pages.RedisOptions{})
	var pc = &handlers.PageController{db}

	http.HandleFunc("/add", pc.AddHandler)
	http.HandleFunc("/", pc.ViewHandler)
	http.ListenAndServe(":9081", nil)
}

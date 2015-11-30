package main

import (
	"github.com/ilonajulczuk/urlich/handlers"
	"github.com/ilonajulczuk/urlich/pages"
	"net/http"
)

func main() {
	redisOptions := &pages.RedisOptions{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	db := pages.NewRedisClient(redisOptions)
	var pc = &handlers.PageController{db}

	http.HandleFunc("/add", pc.AddHandler)
	http.HandleFunc("/", pc.ViewHandler)
	http.ListenAndServe(":9081", nil)
}

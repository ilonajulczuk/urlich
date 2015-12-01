package main

import (
	"github.com/ilonajulczuk/urlich/handlers"
	"github.com/ilonajulczuk/urlich/pages"
	"net/http"
	"os"
	"strconv"
)

func envOrDefault(envName string, defaultVal string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultVal
	}
	return val
}

func main() {
	dbNumber, err := strconv.Atoi(envOrDefault("REDIS_DB", "0"))
	if err != nil {
		panic(err)
	}
	redisOptions := &pages.RedisOptions{
		Addr:     envOrDefault("REDIS_HOST", "localhost:6379"),
		Password: envOrDefault("REDIS_PASS", ""),
		DB:       int64(dbNumber),
	}

	db := pages.NewRedisClient(redisOptions)
	var pc = &handlers.PageController{db}

	http.HandleFunc("/add", pc.AddHandler)
	http.HandleFunc("/", pc.ViewHandler)
	http.ListenAndServe(envOrDefault("PORT", ":9081"), nil)
}

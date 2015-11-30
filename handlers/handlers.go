package handlers

import (
	"encoding/json"
	"errors"
	"gopkg.in/redis.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	maxLimit   = 10 * 1024 * 1024
	PagePrefix = "urlich:page:"
)

var (
	validPath = regexp.MustCompile("^/[a-zA-Z0-9-]{1,32}$")
	client    = connectToRedis()
)

type Page struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	key := findKey(r.URL.Path)
	if key == "" {
		http.NotFound(w, r)
		return
	}
	p, err := loadPage(key)
	if err == nil {
		http.Redirect(w, r, p.URL, http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	var page *Page
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, maxLimit))
	if err != nil {
		log.Println("Can't read request body")
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Println("Can't close request body")
		return
	}
	if err := json.Unmarshal(body, &page); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(442) // Unprocessable entity.
		return
	}

	if err := storePage(page); err != nil {
		log.Println("failure at storing page")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF=8")
	w.WriteHeader(http.StatusCreated)
	log.Println("Created a page", page)
}

func findKey(path string) string {
	key := validPath.FindString(path)
	if key != "" {
		// Remove '/' from beginning.
		return key[1:]
	}
	return key
}

func connectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Unable to connect to redis")
		panic(err)
	}
	return client
}

func loadPage(key string) (*Page, error) {
	rawData, err := client.Get(PagePrefix + key).Result()
	if err != nil {
		return nil, err
	}
	var page *Page
	if rawData == "" {
		return nil, errors.New("not found")
	}
	json.Unmarshal([]byte(rawData), &page)
	return page, nil
}

func storePage(page *Page) error {
	serializedPage, err := json.Marshal(page)
	if err != nil {
		return err
	}

	if err := client.Set(PagePrefix+page.Key, serializedPage, 0).Err(); err != nil {
		return err
	}
	return nil
}

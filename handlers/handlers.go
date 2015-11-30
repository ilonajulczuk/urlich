package handlers

import (
	"encoding/json"
	"github.com/ilonajulczuk/urlich/pages"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	maxLimit = 10 * 1024 * 1024
)

var (
	validPath = regexp.MustCompile("^/[a-zA-Z0-9-]{1,32}$")
)

type PageController struct {
	DB pages.PageClient
}

func (c *PageController) ViewHandler(w http.ResponseWriter, r *http.Request) {
	key := findKey(r.URL.Path)
	if key == "" {
		http.NotFound(w, r)
		return
	}
	p, err := c.DB.LoadPage(key)
	if err == nil {
		http.Redirect(w, r, p.URL, http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func (c *PageController) AddHandler(w http.ResponseWriter, r *http.Request) {
	var page *pages.Page
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

	if err := c.DB.StorePage(page); err != nil {
		log.Println("failure at storing page")
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF=8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(page)
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

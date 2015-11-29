package handlers

import (
	"errors"
	"net/http"
	"regexp"
)

type Page struct {
	Key string
	URL string
}

var pages = map[string]*Page{
	"this": &Page{"this", "http://atte.ro"},
}

var validPath = regexp.MustCompile("^/[a-zA-Z0-9-]{1,32}$")

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
}

func loadPage(key string) (*Page, error) {
	page, ok := pages[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return page, nil
}

func findKey(path string) string {
	key := validPath.FindString(path)
	if key != "" {
		// Remove '/' from beginning.
		return key[1:]
	}
	return key
}

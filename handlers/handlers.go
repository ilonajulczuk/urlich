package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ilonajulczuk/urlich/pages"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
)

const (
	maxLimit    = 10 * 1024 * 1024
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// DefaultKeyLen is default length of generated key.
	DefaultKeyLen = 8
)

var (
	validPath     = regexp.MustCompile("^/[a-zA-Z0-9-]{1,32}$")
	unprocessable = errors.New("unprocessable entity")
)

// ApiError represents client error in api.
type ApiError struct {
	Error string `json:"error"`
}

// PageController handles creating and redirecting to pages.
// Uses DB as an access to data storage.
type PageController struct {
	DB pages.PageClient
}

func createRandomKey() string {
	return randStringBytes(DefaultKeyLen)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// ViewHanler handles redirection to page with given key.
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

// AddHandler creates new page at given key.
func (c *PageController) AddHandler(w http.ResponseWriter, r *http.Request) {
	page, err := parsePage(r.Body)
	if err != nil {
		w.WriteHeader(422)
		return
	}

	var generated bool
	if page.Key == "" {
		page.Key = createRandomKey()
		generated = true
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF=8")
	if err = c.DB.StorePage(page); err != nil {
		if generated {
			for err == pages.PageKeyAlreadyTakenError {
				page.Key = createRandomKey()
				err = c.DB.StorePage(page)
			}
		}
		if err == pages.PageKeyAlreadyTakenError {
			w.WriteHeader(400)
			apiError := &ApiError{"key already taken"}
			json.NewEncoder(w).Encode(apiError)
			return
		}
		if err != nil {
			log.Println("failure at storing page")
			w.WriteHeader(500)
			return
		}
	}

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

func parsePage(rc io.ReadCloser) (*pages.Page, error) {
	body, err := ioutil.ReadAll(io.LimitReader(rc, maxLimit))
	if err != nil {
		return nil, unprocessable
	}
	if err := rc.Close(); err != nil {
		return nil, unprocessable
	}

	var page *pages.Page
	err = json.Unmarshal(body, &page)
	return page, err
}

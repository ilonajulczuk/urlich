package pages

import (
	"encoding/json"
	"errors"
	"gopkg.in/redis.v3"
	"log"
)

const PagePrefix = "urlich:pages:"

type Page struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type PageClient interface {
	LoadPage(string) (*Page, error)
	StorePage(*Page) error
}

type RedisOptions struct {
	Addr     string
	Password string
	DB       int64
}

type RedisPageClient struct {
	Client *redis.Client
}

func NewRedisClient(options *RedisOptions) *RedisPageClient {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Println("Unable to connect to redis")
		panic(err)
	}
	return &RedisPageClient{client}
}

func (r *RedisPageClient) LoadPage(key string) (*Page, error) {
	rawData, err := r.Client.Get(PagePrefix + key).Result()
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

func (r *RedisPageClient) StorePage(page *Page) error {
	serializedPage, err := json.Marshal(page)
	if err != nil {
		return err
	}

	if err := r.Client.Set(PagePrefix+page.Key, serializedPage, 0).Err(); err != nil {
		return err
	}
	return nil
}

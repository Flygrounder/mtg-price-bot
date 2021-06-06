package caching

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type CacheClient struct {
	Storage    *redis.Client
	Expiration time.Duration
}

func NewClient(addr string, passwd string, expiration time.Duration, db int) *CacheClient {
	return &CacheClient{
		Storage: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: passwd,
			DB:       db,
		}),
		Expiration: expiration,
	}
}

func (client *CacheClient) Set(key string, prices []cardsinfo.ScgCardPrice) {
	value, _ := json.Marshal(prices)
	client.Storage.Set(key, value, client.Expiration)
}

func (client *CacheClient) Get(key string) ([]cardsinfo.ScgCardPrice, error) {
	c, err := client.Storage.Get(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "No such key in cache")
	}
	var prices []cardsinfo.ScgCardPrice
	json.Unmarshal([]byte(c), &prices)
	return prices, nil
}

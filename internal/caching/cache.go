package caching

import (
	"time"

	"github.com/go-redis/redis"
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

func (client *CacheClient) Set(key string, value string) {
	client.Storage.Set(key, value, client.Expiration)
}

func (client *CacheClient) Get(key string) (string, error) {
	return client.Storage.Get(key).Result()
}

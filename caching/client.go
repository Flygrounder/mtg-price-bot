package caching

import (
	"github.com/go-redis/redis"
	"time"
)

type CacheClient struct {
	storage    *redis.Client
	Expiration time.Duration
}

var client *CacheClient

func GetClient() *CacheClient {
	if client != nil {
		return client
	}
	client = new(CacheClient)
	client.Init()
	return client
}

func (client *CacheClient) Init() {
	client.storage = redis.NewClient(&redis.Options{
		Addr:     HOST_NAME,
		Password: PASSWORD,
		DB:       0,
	})
	client.Expiration = CACHE_EXPIRATION
}

func (client *CacheClient) Set(key string, value string) {
	client.storage.Set(key, value, client.Expiration)
}

func (client *CacheClient) Get(key string) (string, error) {
	return client.storage.Get(key).Result()
}

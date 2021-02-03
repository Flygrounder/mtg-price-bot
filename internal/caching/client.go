package caching

import (
	"github.com/go-redis/redis"
	"time"
)

type CacheClient struct {
	Storage    *redis.Client
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
	client.Storage = redis.NewClient(&redis.Options{
		Addr:     HostName,
		Password: Password,
		DB:       0,
	})
	client.Expiration = CacheExpiration
}

func (client *CacheClient) Set(key string, value string) {
	client.Storage.Set(key, value, client.Expiration)
}

func (client *CacheClient) Get(key string) (string, error) {
	return client.Storage.Get(key).Result()
}

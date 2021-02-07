package caching

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	c := NewClient("addr", "123", time.Hour, 1)
	assert.Equal(t, time.Hour, c.Expiration)
	assert.Equal(t, 1, c.Storage.Options().DB)
	assert.Equal(t, "addr", c.Storage.Options().Addr)
	assert.Equal(t, "123", c.Storage.Options().Password)
}

func TestGetSet(t *testing.T) {
	client, s := getTestClient()
	defer s.Close()

	keyName := "test_key"
	value := "test_value"
	client.Set(keyName, value)
	val, err := client.Get(keyName)
	assert.Nil(t, err)
	assert.Equal(t, value, val)
}

func TestExpiration(t *testing.T) {
	client, s := getTestClient()
	defer s.Close()

	client.Expiration = time.Millisecond
	keyName := "test_key"
	value := "test_value"
	client.Set(keyName, value)
	s.FastForward(time.Millisecond * 2)
	val, err := client.Get(keyName)
	assert.Zero(t, val)
	assert.NotNil(t, err)
}

func getTestClient() (*CacheClient, *miniredis.Miniredis) {
	s, _ := miniredis.Run()
	fmt.Println(s.Addr())
	c := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	return &CacheClient{
		Storage:    c,
		Expiration: 0,
	}, s
}

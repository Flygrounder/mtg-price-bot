package tests

import (
	"github.com/flygrounder/mtg-price-vk/caching"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	client := getTestClient()
	keyName := "test_key"
	value := "test_value"
	client.Set(keyName, value)
	val, err := client.Get(keyName)
	assert.Nil(t, err)
	assert.Equal(t, value, val)
}

func TestExpiration(t *testing.T) {
	client := getTestClient()
	client.Expiration = time.Millisecond
	keyName := "test_key"
	value := "test_value"
	client.Set(keyName, value)
	time.Sleep(time.Millisecond * 2)
	val, err := client.Get(keyName)
	assert.Zero(t, val)
	assert.NotNil(t, err)
}

func getTestClient() *caching.CacheClient {
	client := new(caching.CacheClient)
	client.Init()
	return client
}

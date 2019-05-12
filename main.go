package main

import (
	"github.com/flygrounder/mtg-price-vk/vk"
	"github.com/gin-gonic/gin"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	go (func() {
		for {
			time.Sleep(5 * time.Second)
			println(runtime.NumGoroutine())
		}
	})()
	rand.Seed(time.Now().UTC().UnixNano())
	r := gin.Default()
	r.POST("callback/message", vk.HandleMessage)
	r.Run(":80")
}

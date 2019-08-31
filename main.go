package main

import (
	"github.com/flygrounder/go-mtg-vk/vk"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	r := gin.Default()
	r.POST("callback/message", vk.HandleMessage)
	r.Run(":80")
}

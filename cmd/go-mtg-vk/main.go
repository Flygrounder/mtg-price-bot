package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/vk"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	logFile, _ := os.OpenFile("logs/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	log.SetOutput(logFile)
	r := gin.Default()
	r.POST("callback/message", vk.HandleMessage)
	_ = r.Run(":8000")
}

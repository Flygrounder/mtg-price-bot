package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/vk"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := gin.Default()

	groupId, _ := strconv.ParseInt(os.Getenv("VK_GROUP_ID"), 10, 64)
	handler := vk.Handler{
		Sender: &vk.ApiSender{
			Token: os.Getenv("VK_TOKEN"),
		},
		Logger:             log.New(os.Stdout, "", 0),
		SecretKey:          os.Getenv("VK_SECRET_KEY"),
		GroupId:            groupId,
		ConfirmationString: os.Getenv("VK_CONFIRMATION_STRING"),
		DictPath:           "./assets/additional_cards.json",
		Cache:              caching.GetClient(),
	}

	r.POST("callback/message", handler.HandleMessage)
	_ = r.Run(":8000")
}

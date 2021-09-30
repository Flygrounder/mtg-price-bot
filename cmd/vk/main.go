package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"

	"github.com/gin-gonic/gin"
	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/vk"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := gin.Default()

	groupId, _ := strconv.ParseInt(os.Getenv("VK_GROUP_ID"), 10, 64)
	dict, _ := os.Open("./assets/additional_cards.json")
	dictBytes, _ := ioutil.ReadAll(dict)
	var dictMap map[string]string
	_ = json.Unmarshal(dictBytes, &dictMap)
	logger := log.New(os.Stdout, "", 0)
	handler := vk.Handler{
		Scenario: &scenario.Scenario{
			Sender: &vk.ApiSender{
				Token:  os.Getenv("VK_TOKEN"),
				Logger: logger,
			},
			Logger: logger,
			Cache:  caching.NewClient("redis:6379", "", time.Hour*24, 0),
			InfoFetcher: &cardsinfo.Fetcher{
				Dict: dictMap,
			},
		},
		SecretKey:          os.Getenv("VK_SECRET_KEY"),
		GroupId:            groupId,
		ConfirmationString: os.Getenv("VK_CONFIRMATION_STRING"),
	}

	r.POST("callback/message", handler.HandleMessage)
	_ = r.Run(":8000")
}

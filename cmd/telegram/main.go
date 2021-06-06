package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"

	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/telegram"
)

func main() {
	dict, _ := os.Open("./assets/additional_cards.json")
	dictBytes, _ := ioutil.ReadAll(dict)
	var dictMap map[string]string
	_ = json.Unmarshal(dictBytes, &dictMap)
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	sender := &telegram.Sender{
	    API: bot,
	}
	sc := &scenario.Scenario{
	    Sender: sender,
	    Logger: log.New(os.Stdout, "", 0),
	    Cache:  caching.NewClient("redis:6379", "", time.Hour*24, 0),
	    InfoFetcher: &cardsinfo.Fetcher{
		Dict: dictMap,
	    },
	}

	u := tgbotapi.NewUpdate(0)
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
	    if update.Message == nil {
		continue
	    }

	    go sc.HandleSearch(&scenario.UserMessage{
		Body: update.Message.Text,
		UserId: update.Message.Chat.ID,
	    })
	} 
}

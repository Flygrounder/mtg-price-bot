package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"
	"gitlab.com/flygrounder/go-mtg-vk/internal/telegram"
)

const welcomeMessage = "Здравствуйте, вас приветствует бот для поиска цен на карты MTG, введите название карты, которая вас интересует."

func main() {
	token, exists := os.LookupEnv("TG_TOKEN")
	if !exists {
		panic("TG_TOKEN environment variable not defined")
	}
	bot, _ := tgbotapi.NewBotAPI(token)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dsn, exists := os.LookupEnv("YDB_CONNECTION_STRING")
	if !exists {
		panic("YDB_CONNECTION_STRING environment variable not defined")
	}

	db, err := ydb.Open(ctx,
		dsn,
		environ.WithEnvironCredentials(ctx),
	)
	if err != nil {
		panic(fmt.Errorf("connect error: %w", err))
	}
	defer func() { _ = db.Close(ctx) }()

	sender := &telegram.Sender{
		API: bot,
	}
	cache := &caching.CacheClient{
		Storage:    db.Table(),
		Expiration: 12 * time.Hour,
		Prefix:     db.Name(),
	}
	err = cache.Init(context.Background())
	if err != nil {
		panic(fmt.Errorf("init error: %w", err))
	}
	sc := &scenario.Scenario{
		Sender:      sender,
		Logger:      log.New(os.Stdout, "", 0),
		Cache:       cache,
		InfoFetcher: &cardsinfo.Fetcher{},
	}

	u := tgbotapi.NewUpdate(0)
	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			sender.Send(update.Message.Chat.ID, welcomeMessage)
			continue
		}

		go sc.HandleSearch(context.Background(), &scenario.UserMessage{
			Body:   update.Message.Text,
			UserId: update.Message.Chat.ID,
		})
	}
}

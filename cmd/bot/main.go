package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gitlab.com/flygrounder/go-mtg-vk/internal/caching"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
	"gitlab.com/flygrounder/go-mtg-vk/internal/scenario"
	"gitlab.com/flygrounder/go-mtg-vk/internal/telegram"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"gitlab.com/flygrounder/go-mtg-vk/internal/vk"
)

type config struct {
	tgToken              string
	vkGroupId            int64
	vkSecretKey          string
	vkToken              string
	vkConfirmationString string
	ydbConnectionString  string
}

func getConfig() *config {
	var cfg config
	var exists bool
	var err error

	cfg.tgToken, exists = os.LookupEnv("TG_TOKEN")
	if !exists {
		panic("TG_TOKEN environment variable not defined")
	}

	vkGroupId, exists := os.LookupEnv("VK_GROUP_ID")
	if !exists {
		panic("VK_GROUP_ID environment variable not defined")
	}
	cfg.vkGroupId, err = strconv.ParseInt(vkGroupId, 10, 64)
	if err != nil {
		panic("VK_GROUP_ID is not a number")
	}

	cfg.vkSecretKey, exists = os.LookupEnv("VK_SECRET_KEY")
	if !exists {
		panic("VK_SECRET_KEY environment variable not defined")
	}

	cfg.vkToken, exists = os.LookupEnv("VK_TOKEN")
	if !exists {
		panic("VK_TOKEN environment variable not defined")
	}

	cfg.vkConfirmationString, exists = os.LookupEnv("VK_CONFIRMATION_STRING")
	if !exists {
		panic("VK_CONFIRMATION_STRING environment variable not defined")
	}

	cfg.ydbConnectionString, exists = os.LookupEnv("YDB_CONNECTION_STRING")
	if !exists {
		panic("YDB_CONNECTION_STRING environment variable not defined")
	}

	return &cfg
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := getConfig()

	bot, _ := tgbotapi.NewBotAPI(cfg.tgToken)
	r := gin.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := ydb.Open(ctx,
		cfg.ydbConnectionString,
		environ.WithEnvironCredentials(ctx),
	)
	if err != nil {
		panic(fmt.Errorf("connect error: %w", err))
	}
	defer func() { _ = db.Close(ctx) }()

	cache := &caching.CacheClient{
		Storage:    db.Table(),
		Expiration: 12 * time.Hour,
		Prefix:     db.Name(),
	}
	err = cache.Init(context.Background())
	if err != nil {
		panic(fmt.Errorf("init error: %w", err))
	}

	logger := log.New(os.Stdout, "", 0)
	fetcher := &cardsinfo.Fetcher{}
	handler := vk.Handler{
		Scenario: &scenario.Scenario{
			Sender: &vk.ApiSender{
				Token:  cfg.vkToken,
				Logger: logger,
			},
			Logger:      logger,
			InfoFetcher: fetcher,
			Cache:       cache,
		},
		SecretKey:          cfg.vkSecretKey,
		GroupId:            cfg.vkGroupId,
		ConfirmationString: cfg.vkConfirmationString,
	}

	tgHandler := telegram.Handler{
		Scenario: &scenario.Scenario{
			Sender: &telegram.Sender{
				API: bot,
			},
			Logger:      logger,
			InfoFetcher: fetcher,
			Cache:       cache,
		},
	}

	r.POST("vk", handler.HandleMessage)
	r.POST("tg", tgHandler.HandleMessage)
	_ = r.Run(":8000")
}

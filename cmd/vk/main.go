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

	"github.com/gin-gonic/gin"
	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"gitlab.com/flygrounder/go-mtg-vk/internal/vk"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := gin.Default()

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

	cache := &caching.CacheClient{
		Storage:    db.Table(),
		Expiration: 12 * time.Hour,
		Prefix:     db.Name(),
	}
	err = cache.Init(context.Background())
	if err != nil {
		panic(fmt.Errorf("init error: %w", err))
	}

	groupId, _ := strconv.ParseInt(os.Getenv("VK_GROUP_ID"), 10, 64)
	logger := log.New(os.Stdout, "", 0)
	handler := vk.Handler{
		Scenario: &scenario.Scenario{
			Sender: &vk.ApiSender{
				Token:  os.Getenv("VK_TOKEN"),
				Logger: logger,
			},
			Logger:      logger,
			InfoFetcher: &cardsinfo.Fetcher{},
			Cache:       cache,
		},
		SecretKey:          os.Getenv("VK_SECRET_KEY"),
		GroupId:            groupId,
		ConfirmationString: os.Getenv("VK_CONFIRMATION_STRING"),
	}

	r.POST("callback/message", handler.HandleMessage)
	_ = r.Run(":8000")
}

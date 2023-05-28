package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"gitlab.com/flygrounder/go-mtg-vk/internal/cardsinfo"
)

type CacheClient struct {
	Storage    table.Client
	Expiration time.Duration
	Prefix     string
}

func (client *CacheClient) Init(ctx context.Context) error {
	return client.Storage.Do(ctx, func(ctx context.Context, s table.Session) error {
		return s.CreateTable(
			ctx,
			path.Join(client.Prefix, "cache"),
			options.WithColumn("card", types.TypeString),
			options.WithColumn("prices", types.Optional(types.TypeJSON)),
			options.WithColumn("created_at", types.Optional(types.TypeTimestamp)),
			options.WithTimeToLiveSettings(
				options.NewTTLSettings().ColumnDateType("created_at").ExpireAfter(client.Expiration),
			),
			options.WithPrimaryKeyColumn("card"),
		)
	})
}

func (client *CacheClient) Set(ctx context.Context, key string, prices []cardsinfo.ScgCardPrice) error {
	const query = `
	DECLARE $cacheData AS List<Struct<
		card: String,
		prices: Json,
		created_at: Timestamp>>;

	INSERT INTO cache SELECT cd.card AS card, cd.prices AS prices, cd.created_at AS created_at FROM AS_TABLE($cacheData) cd LEFT OUTER JOIN cache c ON cd.card = c.card WHERE c.card IS NULL`
	value, _ := json.Marshal(prices)
	return client.Storage.Do(ctx, func(ctx context.Context, s table.Session) (err error) {
		_, _, err = s.Execute(ctx, writeTx(), query, table.NewQueryParameters(
			table.ValueParam("$cacheData", types.ListValue(
				types.StructValue(
					types.StructFieldValue("card", types.StringValueFromString(key)),
					types.StructFieldValue("prices", types.JSONValueFromBytes(value)),
					types.StructFieldValue("created_at", types.TimestampValueFromTime(time.Now())),
				))),
		))
		return
	})
}

func (client *CacheClient) Get(ctx context.Context, key string) ([]cardsinfo.ScgCardPrice, error) {
	const query = `
	DECLARE $card AS String;

	SELECT UNWRAP(prices) AS prices FROM cache WHERE card = $card`
	var pricesStr string
	var res result.Result
	err := client.Storage.Do(ctx, func(ctx context.Context, s table.Session) (err error) {
		_, res, err = s.Execute(ctx, readTx(), query, table.NewQueryParameters(
			table.ValueParam("$card", types.StringValueFromString(key)),
		))
		return
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get key")
	}
	ok := res.NextResultSet(ctx)
	if !ok {
		return nil, errors.New("no key")
	}
	ok = res.NextRow()
	if !ok {
		return nil, errors.New("no key")
	}
	err = res.ScanNamed(
		named.Required("prices", &pricesStr),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan prices: %w", err)
	}
	var prices []cardsinfo.ScgCardPrice
	_ = json.Unmarshal([]byte(pricesStr), &prices)
	return prices, nil
}

func writeTx() *table.TransactionControl {
	return table.TxControl(table.BeginTx(
		table.WithSerializableReadWrite(),
	), table.CommitTx())
}

func readTx() *table.TransactionControl {
	return table.TxControl(table.BeginTx(
		table.WithOnlineReadOnly(),
	), table.CommitTx())
}

package internal

import (
	"context"
	"time"

	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v3"
)

type CoinGecko interface {
	FetchUsdPriceMap(denoms ...string) (map[string]float64, error)
	GetPriceInUsd(denoms ...string) (map[string]map[string]float64, error)
	GetAPIKey() string
}

type Bot interface {
	Delete(msg telebot.Editable) error
	Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error)
	Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc)
	ProcessUpdate(u telebot.Update)
}

type BotClient interface {
	GetPrice(ticker string) (float64, bool)
	SetPrice(ticker string, price float64)
	GetBalances(address string, denoms []string) (map[string]float64, error)
	TransferToken(to string, amount float64, denom string) (string, error)
	GetAddress() string
	GetDecimal(denom string) int32
	PlaceSpotOrder(denomIn, denomOut string, amount, price float64) (string, error)
	GetActiveOrders(marketId string) ([]types.LimitOrderInfo, error)
	CancelOrder(marketID, orderHash string) (string, error)
	ToMessage(order types.LimitOrderInfo, showDetail bool) string
	GetRedisInstance() RedisClient
}

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
}

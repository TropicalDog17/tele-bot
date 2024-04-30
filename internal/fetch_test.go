package internal_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/go-redis/redismock/v9"
)

func MockFetchUsdPriceMap(redisClient internal.RedisClient, coinGeckoClient internal.CoinGecko, tokens ...string) error {
	ctx := context.Background()
	for _, token := range tokens {
		// Fetch data from redis
		tokenKey := fmt.Sprintf("price:%s", token)
		coinGeckoID := internal.ConvertToCoinGeckoTicker(token)
		// If data is not found in redis, fetch from CoinGecko
		// If data is found in CoinGecko, store it in redis
		price, err := redisClient.Get(ctx, tokenKey).Result()
		if err != nil || price == "" {
			// Fetch from CoinGecko
			fetchedPrice, err := internal.MockGetPriceInUsd(coinGeckoID, coinGeckoClient)
			if err != nil {
				return err
			}
			// Store in redis
			err = redisClient.Set(ctx, tokenKey, fmt.Sprintf("%f", fetchedPrice), 0).Err()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func TestFetchUsdPriceMapRedis(t *testing.T) {
	db, mock := redismock.NewClientMock()

	key := "price:atom"

	// mock ignoring `call api()`

	mock.ExpectGet(key).RedisNil()
	mock.Regexp().ExpectSet(key, "1.2", 0).SetErr(errors.New("FAIL"))

	err := MockFetchUsdPriceMap(db, &internal.CoinGeckoClient{}, "atom")
	if err == nil || err.Error() != "FAIL" {
		t.Error("wrong error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

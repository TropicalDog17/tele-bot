package internal_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/TropicalDog17/tele-bot/internal"
	mock_internal "github.com/TropicalDog17/tele-bot/tests/mocks"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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

func TestFetchMarkets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db, mock := redismock.NewClientMock()
	mockExchangeClient := mock_internal.NewMockExchangeClient(ctrl)

	mockPair1Response := &spotExchangePB.SpotMarketInfo{
		MarketId:     "pair1",
		MarketStatus: "active",
		Ticker:       "ATOM/INJ",
		BaseDenom:    "ATOM",
		QuoteDenom:   "INJ",
	}
	mockPair2Response := &spotExchangePB.SpotMarketInfo{
		MarketId:     "pair2",
		MarketStatus: "active",
		Ticker:       "BTC/USDT",
		BaseDenom:    "BTC",
		QuoteDenom:   "USDT",
	}
	key := "markets"
	pair1 := "ATOM/INJ"
	pair2 := "BTC/USDT"

	mockExchangeClient.EXPECT().GetActiveMarkets(gomock.Any(), gomock.Any()).Return([]*spotExchangePB.SpotMarketInfo{mockPair1Response, mockPair2Response}, nil)
	mock.ExpectHSet(key, pair1, mockPair1Response.MarketId).SetVal(1)
	mock.ExpectHSet(key, pair2, mockPair2Response.MarketId).SetVal(1)
	err := internal.FetchMarkets(db, mockExchangeClient)
	require.NoError(t, err)
}

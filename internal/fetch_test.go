package internal_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
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

// func TestFetchUsdPriceMapRedis(t *testing.T) {
// 	db, mock := redismock.NewClientMock()

// 	key := "price:atom"

// 	// mock ignoring `call api()`

// 	mock.ExpectGet(key).RedisNil()
// 	mock.Regexp().ExpectSet(key, "1.2", 0).SetErr(errors.New("FAIL"))

// 	err := MockFetchUsdPriceMap(db, internal.CoinGecko, "atom")
// 	if err == nil || err.Error() != "FAIL" {
// 		t.Error("wrong error")
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Error(err)
// 	}
// }

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

func TestSyncOrderHappyCase(t *testing.T) {
	mockOrder1 := types.LimitOrderInfo{
		OrderHash: "order1",
		MarketID:  "pair1",
	}
	mockOrder2 := types.LimitOrderInfo{
		OrderHash: "order2",
		MarketID:  "pair1",
	}
	mockOrder3 := types.LimitOrderInfo{
		OrderHash: "order3",
		MarketID:  "pair2",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, mock := redismock.NewClientMock()

	// GetActiveMarkets returns a map of market pairs
	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1", "BTC/USDT": "pair2"}, nil)
	botClient.EXPECT().GetActiveOrders("pair1").Return([]types.LimitOrderInfo{mockOrder1, mockOrder2}, nil)
	botClient.EXPECT().GetActiveOrders("pair2").Return([]types.LimitOrderInfo{mockOrder3}, nil)
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()

	mock.MatchExpectationsInOrder(false) // Disable strict order matching
	mock.ExpectHKeys("address").SetVal([]string{})

	// Expect HGet, HSet, and HSet for each order in any order
	for _, order := range []types.LimitOrderInfo{mockOrder1, mockOrder2, mockOrder3} {
		mock.ExpectHGet("address", order.OrderHash).SetVal("")
		mock.ExpectHSet("address", order.OrderHash, LimitOrderInfoToJson(order)).SetVal(1)
		mock.ExpectHSet("orders", order.OrderHash, order.MarketID).SetVal(1)
	}

	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.NoError(t, err)
}

// TestSyncOrderError1 tests the case where GetActiveMarkets returns an error
func TestSyncOrderError1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, _ := redismock.NewClientMock()

	botClient.EXPECT().GetActiveMarkets().Return(nil, errors.New("error"))
	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.Error(t, err)
}

// TestSyncOrderError2 tests the case where GetActiveOrders returns an error
func TestSyncOrderError2(t *testing.T) {
	mockOrder1 := types.LimitOrderInfo{
		OrderHash: "order1",
		MarketID:  "pair1",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, mock := redismock.NewClientMock()

	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1"}, nil)
	mock.ExpectHKeys("address").SetVal([]string{})

	botClient.EXPECT().GetActiveOrders("pair1").Return(nil, errors.New("error"))
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()

	mock.ExpectHSet("address", mockOrder1.OrderHash, mockOrder1).SetVal(1)
	mock.ExpectHSet("orders", "order1", "pair1").SetVal(1)

	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.Error(t, err)
}

// TestSyncOrderError3 tests the case where SAdd returns an error
func TestSyncOrderError3(t *testing.T) {
	mockOrder1 := types.LimitOrderInfo{
		OrderHash: "order1",
		MarketID:  "pair1",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, mock := redismock.NewClientMock()

	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1"}, nil)
	mock.ExpectHKeys("address").SetVal([]string{})
	botClient.EXPECT().GetActiveOrders("pair1").Return([]types.LimitOrderInfo{mockOrder1}, nil)
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()

	mock.ExpectHSet("address", mockOrder1.OrderHash, mockOrder1).SetErr(errors.New("error"))

	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.Error(t, err)
}

// TestSyncOrderError4 tests the case where HSet returns an error
func TestSyncOrderError4(t *testing.T) {
	mockOrder1 := types.LimitOrderInfo{
		OrderHash: "order1",
		MarketID:  "pair1",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, mock := redismock.NewClientMock()

	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1"}, nil)
	mock.ExpectHKeys("address").SetVal([]string{})

	botClient.EXPECT().GetActiveOrders("pair1").Return([]types.LimitOrderInfo{mockOrder1}, nil)
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()

	mock.ExpectHSet("address", mockOrder1.OrderHash, mockOrder1).SetVal(1)
	mock.ExpectHSet("orders", "order1", "pair1").SetErr(errors.New("error"))

	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.Error(t, err)
}

func TestSyncAndPrune(t *testing.T) {
	mockOrder1 := types.LimitOrderInfo{
		OrderHash: "order1",
		MarketID:  "pair1",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	botClient := mock_internal.NewMockBotClient(ctrl)
	rdb, mock := redismock.NewClientMock()

	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1"}, nil)
	mock.ExpectHKeys("address").SetVal([]string{"order1"})
	botClient.EXPECT().GetActiveOrders("pair1").Return([]types.LimitOrderInfo{mockOrder1}, nil)
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()
	mock.ExpectHGet("address", "order1").SetVal(LimitOrderInfoToJson(mockOrder1))
	err := internal.SyncOrdersToRedis(botClient, rdb)
	require.NoError(t, err)

	// order in redis, but onchain data order is already completed(not active)
	botClient.EXPECT().GetActiveMarkets().Return(map[string]string{"ATOM/INJ": "pair1"}, nil)
	mock.ExpectHKeys("address").SetVal([]string{"order1"})
	botClient.EXPECT().GetActiveOrders("pair1").Return([]types.LimitOrderInfo{}, nil)
	botClient.EXPECT().GetAddress().Return("address").AnyTimes()

	mock.ExpectHDel("orders", "order1").SetVal(1)
	mock.ExpectHDel("address", "order1").SetVal(1)

	err = internal.SyncOrdersToRedis(botClient, rdb)
	require.NoError(t, err)

}

func LimitOrderInfoToJson(order types.LimitOrderInfo) string {
	jsonBytes, _ := json.Marshal(order)
	return string(jsonBytes)
}

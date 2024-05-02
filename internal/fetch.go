package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	"github.com/redis/go-redis/v9"
)

// Fetch market data, market id, and market summary from the exchange client
func FetchDataWithTimeout(redisClient *redis.Client, coinGeckoClient CoinGecko, exchangeClient ExchangeClient) {
	// goroutine fetch usd price map
	go func() {
		err := FetchUsdPriceMap(redisClient, coinGeckoClient, "inj", "atom")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(30 * time.Minute)
	}()

	go func() {
		err := FetchMarkets(redisClient, exchangeClient)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(30 * time.Minute)
	}()

}

func FetchUsdPriceMap(redisClient RedisClient, coinGeckoClient CoinGecko, tokens ...string) error {
	ctx := context.Background()
	for _, token := range tokens {
		// Fetch data from redis
		tokenKey := fmt.Sprintf("price:%s", token)
		token = ConvertToCoinGeckoTicker(token)
		// If data is not found in redis, fetch from CoinGecko
		// If data is found in CoinGecko, store it in redis
		price, err := redisClient.Get(ctx, tokenKey).Result()
		if err != nil || price == "" {
			// Fetch from CoinGecko
			fetchedPrice, err := GetPriceInUsd(token, coinGeckoClient)
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

// TODO: handle convert to coin gecko ticker
func ConvertToCoinGeckoTicker(denom string) string {
	if denom == "inj" {
		return "injective-protocol"
	} else if denom == "atom" {
		return "cosmos"
	}
	return denom
}

func GetPriceInUsd(denom string, coinGeckoClient CoinGecko) (float64, error) {

	// Request by http
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids="+denom+"&vs_currencies=usd", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-cg-demo-api-key", coinGeckoClient.GetAPIKey())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	// json result is like this:
	// {
	//   "cosmos": {
	//     "usd": 8.66
	//   }
	// }
	var result map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	if result[denom] == nil {
		return 0, errors.New("denom not found")
	}
	return result[denom]["usd"], nil
}

func MockGetPriceInUsd(_ string, _ CoinGecko) (float64, error) {
	return 1.2, nil
}

func FetchMarkets(redisClient RedisClient, ExchangeClient ExchangeClient) error {
	ctx := context.Background()
	req := &spotExchangePB.MarketsRequest{
		MarketStatus: "active",
	}
	res, err := ExchangeClient.GetActiveMarkets(ctx, req)
	if err != nil {
		return err
	}
	for _, market := range res {
		ticker := market.Ticker
		marketId := market.MarketId
		// Store in redis
		err = redisClient.HSet(ctx, "markets", ticker, marketId).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncOrdersToRedis(client BotClient, rdb RedisClient) error {
	markets, err := client.GetActiveMarkets()
	if err != nil {
		return fmt.Errorf("error fetching active markets: %v", err)
	}

	for _, marketID := range markets {
		marketOrders, err := client.GetActiveOrders(marketID)
		if err != nil {
			return fmt.Errorf("error fetching active orders for market %s: %v", marketID, err)
		}
		ctx := context.Background()
		for _, order := range marketOrders {
			orderID := order.OrderHash
			_, err = rdb.SAdd(ctx, client.GetAddress(), order).Result()
			if err != nil {
				return fmt.Errorf("error syncing order %s to Redis: %v", client.GetAddress(), err)
			}
			_, err = rdb.HSet(ctx, client.GetAddress(), orderID, marketID).Result()
			if err != nil {
				return fmt.Errorf("error mapping order %s to marketid to redis: %v", orderID, err)
			}
		}
	}

	return nil
}

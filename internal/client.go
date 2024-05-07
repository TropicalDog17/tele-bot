package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/utils"
	"github.com/TropicalDog17/tele-bot/internal/database"
	"github.com/TropicalDog17/tele-bot/internal/types"
)

type Client struct {
	client          *exchange.MbClient
	coinGeckoClient CoinGecko
	redisClient     RedisClient
}

type CoinGeckoClient struct {
	apiKey string
}

func NewCoinGeckoClient() *CoinGeckoClient {
	return &CoinGeckoClient{
		apiKey: os.Getenv("COINGECKO_API_KEY"),
	}
}
func InitExchangeClient() *exchange.MbClient {
	exchangeClient := exchange.NewMbClient("local", configtypes.DefaultConfig())
	return exchangeClient
}

func NewClient(privateKey string) *Client {
	client := exchange.NewMbClient("local", privateKey, configtypes.DefaultConfig())
	client.ChainClient.AdjustKeyring("user4")
	redisClient := database.NewRedisInstance()

	cgClient := NewCoinGeckoClient()
	go FetchDataWithTimeout(redisClient, cgClient, client)
	c := &Client{
		client:          client,
		coinGeckoClient: cgClient,
		redisClient:     redisClient,
	}
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Sync orders to redis")
			err := SyncOrdersToRedis(c, c.redisClient)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	return c
}

func (c *Client) GetPrice(ticker string) (float64, bool) {
	ctx := context.Background()
	price, found := c.redisClient.Get(ctx, fmt.Sprintf("price:%s", ticker)).Result()
	if found != nil {
		return 0, false
	}
	floatPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, false
	}
	return floatPrice, true
}

func (c *CoinGeckoClient) FetchUsdPriceMap(denoms ...string) (map[string]float64, error) {
	// TODO: fix the hardcode
	priceMap, err := c.GetPriceInUsd("inj", "atom")
	if err != nil {
		return nil, err
	}
	result := make(map[string]float64)
	for key, value := range priceMap {
		result[key] = value["usd"]
	}
	return result, nil
}
func (c *Client) SetPrice(ticker string, price float64) {
	ctx := context.Background()
	c.redisClient.Set(ctx, fmt.Sprintf("price:%s", ticker), fmt.Sprintf("%f", price), 0)
}

func (c *Client) GetBalances(address string, denoms []string) (map[string]float64, error) {
	balances := make(map[string]float64)
	for _, denom := range denoms {
		balance, err := c.client.ChainClient.GetBalance(address, denom)
		if err != nil {
			return nil, err
		}
		balances[denom] = balance
	}
	return balances, nil
}

func (c *Client) TransferToken(to string, amount float64, denom string) (string, error) {
	return c.client.ChainClient.TransferToken(to, amount, denom)
}

func (c *Client) GetAddress() string {
	return c.client.ChainClient.SenderAddress.String()
}

// This works for most of the tokens
func (c *Client) GetDecimal(denom string) int32 {
	if denom == "inj" {
		return 18
	}
	return 6
}

func (c *Client) PlaceSpotOrder(denomIn, denomOut string, amount, price float64) (string, error) {
	spotMarket, err := c.FetchSpotMarket(denomIn, denomOut)
	if err != nil {
		return "", err
	}

	orderType := DetermineOrderType(spotMarket, denomIn, denomOut)
	order := c.client.NewSpotOrder(orderType, spotMarket.MarketId, price, amount)
	txHash, err := c.client.PlaceSpotOrder(order)
	if err != nil {
		return "", err
	}
	return txHash, nil

}
func (c *Client) FetchSpotMarket(denomIn, denomOut string) (*exchangetypes.SpotMarket, error) {
	ticker1 := denomIn + denomOut
	ticker2 := denomOut + denomIn
	spot1, err1 := c.client.GetSpotMarketFromTicker(ticker1)
	spot2, err2 := c.client.GetSpotMarketFromTicker(ticker2)
	if err1 != nil && err2 != nil {
		return nil, err1
	}
	if err1 != nil {
		return spot2, nil
	}
	return spot1, nil
}

func DetermineOrderType(spotMarket *exchangetypes.SpotMarket, denomIn, denomOut string) exchangetypes.OrderType {

	var orderType exchangetypes.OrderType
	if spotMarket.BaseDenom == denomIn {
		orderType = exchangetypes.OrderType_BUY
	} else {
		orderType = exchangetypes.OrderType_SELL
	}
	return orderType
}

func (c *CoinGeckoClient) GetPriceInUsd(denoms ...string) (map[string]map[string]float64, error) {
	var ticker []string
	for _, denom := range denoms {
		if denom == "inj" {
			ticker = append(ticker, "injective-protocol")
		} else if denom == "atom" {
			ticker = append(ticker, "cosmos")
		} else {
			ticker = append(ticker, denom)
		}
	}

	tickerString := strings.Join(ticker, ",")
	// Request by http
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids="+tickerString+"&vs_currencies=usd", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-cg-demo-api-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// json result is like this:
	// {
	//   "cosmos": {
	//     "usd": 8.66
	//   }
	// }

	// Parse the json result
	var result map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	// Replace ticker to denom
	for _, denom := range denoms {
		if denom == "inj" {
			result["inj"] = result["injective-protocol"]
			delete(result, "injective-protocol")
		} else if denom == "atom" {
			result["atom"] = result["cosmos"]
			delete(result, "cosmos")
		}
	}
	return result, nil

}

func (c *CoinGeckoClient) GetAPIKey() string {
	return c.apiKey
}

func (c *Client) ToMessage(order types.LimitOrderInfo, showDetail bool) string {
	priceOut, _ := c.GetPrice(order.DenomOut)
	priceIn, _ := c.GetPrice(order.DenomIn)
	if !showDetail {

		return fmt.Sprintf(`📊 Limit Order - %s
	⬩ Mode: %s
	⬩ Amount: %.3f %s
	⬩ Limit Price: $%.3f (0.00%%)
	⬩ Pay Token: %s
	⬩ OrderID: %s
	⬩ MarketID: %s
	`, order.Direction, order.Direction, order.Amount, order.DenomIn, order.Price, order.DenomOut, order.OrderHash, order.MarketID)
	} else {
		return fmt.Sprintf(`📊 Limit Order - %s
	⬩ Mode: %s
	⬩ Amount: %.3f %s
	⬩ Limit Price: $%.3f (0.00%%)
	⬩ Pay Token: %s
	⬩ OrderID: %s
	⬩ MarketID: %s
	After the order is filled: 
	You will receive: %.3f %s ($%.3f)
	You will pay: %.3f %s ($%.3f)`, order.Direction, order.Direction, order.Amount, order.DenomIn, order.Price, order.DenomOut, order.OrderHash, order.MarketID, order.Amount, order.DenomIn, order.Amount*priceIn, order.Amount*order.Price, order.DenomOut, order.Amount*order.Price*priceOut)
	}

}

func (c *Client) GetActiveOrders(marketId string) ([]types.LimitOrderInfo, error) {
	ctx := context.Background()
	orders, err := c.client.ChainClient.GetInjectiveChainClient().FetchChainAccountAddressSpotOrders(ctx, marketId, c.GetAddress())
	if err != nil {
		return nil, err
	}
	if len(orders.Orders) == 0 {
		return nil, nil
	}
	out := make([]types.LimitOrderInfo, 0)
	marketInfo, err := c.client.GetSpotMarket(marketId)
	if err != nil {
		return nil, err
	}

	// TODO: handle pagination
	if len(orders.Orders) > 10 {
		// get the last 10 orders
		orders.Orders = orders.Orders[len(orders.Orders)-10:]
	}
	for _, order := range orders.Orders {
		parsedOrder := types.LimitOrderInfo{}
		if order.IsBuy {
			parsedOrder.DenomIn = marketInfo.BaseDenom
			parsedOrder.DenomOut = marketInfo.QuoteDenom
			parsedOrder.Direction = "buy"
		} else {
			parsedOrder.DenomIn = marketInfo.QuoteDenom
			parsedOrder.DenomOut = marketInfo.BaseDenom
			parsedOrder.Direction = "sell"
		}
		parsedOrder.Price = utils.PriceFromChainFormat(order.Price.String(), c.GetDecimal(parsedOrder.DenomIn), c.GetDecimal(parsedOrder.DenomOut)).InexactFloat64()
		parsedOrder.Amount = utils.QuantityFromChainFormat(order.Quantity.String(), c.GetDecimal(parsedOrder.DenomIn)).InexactFloat64()
		parsedOrder.OrderHash = order.OrderHash
		parsedOrder.MarketID = marketId
		out = append(out, parsedOrder)
	}
	return out, nil
}

func (c *Client) CancelOrder(marketID, orderHash string) (string, error) {
	ctx := context.Background()
	txhash, err := c.client.CancelOrder(ctx, marketID, orderHash)
	if err != nil {
		return "", err
	}
	return txhash, nil
}

func (c *Client) GetRedisInstance() RedisClient {
	return c.redisClient
}

func (c *Client) GetActiveMarkets() (map[string]string, error) {
	ctx := context.Background()
	markets, err := c.GetRedisInstance().HGetAll(ctx, "markets").Result()
	if err != nil {
		return nil, err
	}
	return markets, nil
}

func (c *Client) GetExchangeClient() *exchange.MbClient {
	return c.client
}

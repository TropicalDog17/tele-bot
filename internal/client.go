package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/utils"
	"github.com/TropicalDog17/tele-bot/internal/types"
)

type Client struct {
	client          *exchange.MbClient
	coinGeckoClient *CoinGeckoClient
	priceMap        map[string]float64
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

func NewClient() *Client {
	client := exchange.NewMbClient("local", configtypes.DefaultConfig())
	client.ChainClient.AdjustKeyring("user3")

	cgClient := NewCoinGeckoClient()
	priceMap, err := cgClient.FetchUsdPriceMap("inj", "atom")
	if err != nil {
		panic(err)
	}
	c := &Client{
		client:          client,
		coinGeckoClient: cgClient,
		priceMap:        priceMap,
	}

	return c
}

func (c *Client) GetPrice(ticker string) (float64, bool) {
	price, found := c.priceMap[ticker]
	return price, found
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
	c.priceMap[ticker] = price
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
	// curl --request GET \
	//  --url 'https://api.coingecko.com/api/v3/simple/price?ids=cosmos&vs_currencies=usd' \
	//  --header 'accept: application/json' \
	//  --header 'x-cg-demo-api-key: CG-NahJ6HtSXWRtV1ASvP6EZMCS'
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

// 📊 Limit Order - Buy
// ⬩ Mode: Buy
// ⬩ Token: SOL
// ⬩ Amount: 1.000000 SOL
// ⬩ Limit Price: $0.006382091 (0.00%)
// IN:   1.000000 SOL ($156.69)
// OUT: 156.688456 USDC ($156.69)
func (c *Client) ToMessage(order types.LimitOrderInfo) string {
	priceOut := c.priceMap[order.DenomOut]
	priceIn := c.priceMap[order.DenomIn]
	return fmt.Sprintf(`📊 Limit Order - %s
	⬩ Mode: %s
	⬩ Amount: %.3f %s
	⬩ Limit Price: $%.3f (0.00%%)
	⬩ Pay Token: %s
	You will receive:   %.3f %s ($%.3f)
	You will pay: %.3f %s ($%.3f)`, order.Direction, order.Direction, order.Amount, order.DenomIn, order.Price, order.DenomOut, order.Amount, order.DenomIn, order.Amount*priceIn, order.Amount*order.Price, order.DenomOut, order.Amount*order.Price*priceOut)
}

func (c *Client) GetActiveOrders(marketId string) ([]types.LimitOrderInfo, error) {
	ctx := context.Background()
	orders, err := c.client.ChainClient.GetInjectiveChainClient().FetchChainAccountAddressSpotOrders(ctx, marketId, c.GetAddress())
	if err != nil {
		fmt.Println(err)
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
		out = append(out, parsedOrder)
	}
	return out, nil
}

package internal

import (
	"encoding/json"
	"net/http"
	"os"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
)

type Client struct {
	client          *exchange.MbClient
	coinGeckoClient *CoinGeckoClient
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
	return &Client{
		client:          client,
		coinGeckoClient: cgClient,
	}
}

func (c *Client) GetPrice(ticker string) (float64, error) {
	return c.client.GetPrice(ticker)
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

func (c *CoinGeckoClient) GetPriceInUsd(denom string) (float64, error) {
	var ticker string
	if denom == "atom" {
		ticker = "cosmos"
	}
	if denom == "inj" {
		ticker = "injective-protocol"
	}
	// curl --request GET \
	//  --url 'https://api.coingecko.com/api/v3/simple/price?ids=cosmos&vs_currencies=usd' \
	//  --header 'accept: application/json' \
	//  --header 'x-cg-demo-api-key: CG-NahJ6HtSXWRtV1ASvP6EZMCS'

	// Request by http
	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?ids="+ticker+"&vs_currencies=usd", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-cg-demo-api-key", c.apiKey)

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

	// Parse the json result
	var result map[string]map[string]float64
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}
	return result[ticker]["usd"], nil
}

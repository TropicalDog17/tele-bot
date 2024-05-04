package exchange

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"
	spotExchangePB "github.com/InjectiveLabs/sdk-go/exchange/spot_exchange_rpc/pb"
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/chain"
	customtypes "github.com/TropicalDog17/orderbook-go-sdk/pkg/types"
)

var _ CronosClient = (*MbClient)(nil)
var _ ExchangeFetcher = (*MbClient)(nil)

type ExchangeFetcher interface {
	GetPrice(ticker string) (float64, error)
}

type WalletFetcher interface {
	GetBalance() (float64, error)
}

type CronosClient interface {
	GetMarketSummary(marketId string) (customtypes.MarketSummary, error)
}
type MbClient struct {
	ExchangeClient exchangeclient.ExchangeClient
	ChainClient    *chain.ChainClient
	Config         *configtypes.Config
}

func NewMbClient(networkType string, config *configtypes.Config) *MbClient {
	if networkType != "local" {
		panic("Only local network type is supported")
	}

	network := configtypes.DefaultNetwork()
	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}
	chainClient := chain.NewChainClient("genesis") // TODO: refactor hard code
	return &MbClient{
		ExchangeClient: exchangeClient,
		ChainClient:    &chainClient,
		Config:         config,
	}
}

func (c *MbClient) GetPrice(ticker string) (float64, error) {
	ticker = strings.Replace(ticker, "-", "", -1)
	ticker = strings.Replace(ticker, "/", "", -1)
	ticker = strings.ToUpper(ticker)
	marketId := os.Getenv(ticker)
	if marketId == "" {
		return 0, fmt.Errorf("marketId not found for ticker %s", ticker)
	}
	marketSummary, err := c.GetMarketSummary(marketId)
	if err != nil {
		return 0, err
	}
	return marketSummary.Price, nil
}

func (c *MbClient) GetMarketSummary(marketId string) (customtypes.MarketSummary, error) {
	// TODO: fix hard code

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	endpoint := fmt.Sprintf("%s/api/chronos/v1/spot/market_summary?marketId=%s&resolution=24h", c.Config.ChronosEndpoint, marketId)
	var marketSummary customtypes.MarketSummary
	resp, err := client.Get(endpoint)

	if err != nil {
		return marketSummary, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return marketSummary, err
	}
	if err := json.Unmarshal(bodyBytes, &marketSummary); err != nil {
		return marketSummary, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return marketSummary, nil
}

func (c *MbClient) GetSpotMarket(marketId string) (*exchangetypes.SpotMarket, error) {
	ctx := context.Background()
	market, err := c.ChainClient.GetInjectiveChainClient().FetchChainSpotMarket(ctx, marketId)
	if err != nil {
		panic(err)
	}
	return market.Market, nil
}

// func (c *MbClient) GetMarketsAssistant() chainclient.MarketsAssistant {
// 	ctx := context.Background()

// 	marketsAssistant, err := chainclient.NewMarketsAssistantInitializedFromChain(ctx, *c.exchangeClient)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return marketsAssistant
// }

func (c *MbClient) GetChainClient() *chain.ChainClient {
	return c.ChainClient
}

func (c *MbClient) GetDecimals(ctx context.Context, marketId string) (baseDecimal, quoteDecimal int32) {
	market, err := c.ExchangeClient.GetSpotMarket(ctx, marketId)
	if err != nil {
		panic(err)
	}
	baseDecimal = market.Market.BaseTokenMeta.Decimals
	quoteDecimal = market.Market.QuoteTokenMeta.Decimals
	return baseDecimal, quoteDecimal
}

func (c *MbClient) GetMarketSummaryFromTicker(ticker string) (customtypes.MarketSummary, error) {
	ticker = strings.Replace(ticker, "-", "", -1)
	ticker = strings.Replace(ticker, "/", "", -1)
	ticker = strings.ToUpper(ticker)
	marketId := os.Getenv(ticker)
	if marketId == "" {
		return customtypes.MarketSummary{}, fmt.Errorf("marketId not found for ticker %s", ticker)
	}
	marketSummary, err := c.GetMarketSummary(marketId)
	if err != nil {
		return customtypes.MarketSummary{}, err
	}
	return marketSummary, err
}
func (c *MbClient) GetSpotMarketFromTicker(ticker string) (*exchangetypes.SpotMarket, error) {
	ticker = strings.Replace(ticker, "-", "", -1)
	ticker = strings.Replace(ticker, "/", "", -1)
	ticker = strings.ToUpper(ticker)
	marketId := os.Getenv(ticker)
	if marketId == "" {
		return &exchangetypes.SpotMarket{}, fmt.Errorf("marketId not found for ticker %s", ticker)
	}
	spotMarket, err := c.GetSpotMarket(marketId)
	if err != nil {
		return &exchangetypes.SpotMarket{}, err
	}
	return spotMarket, err
}

func (c *MbClient) GetActiveMarkets(ctx context.Context, req *spotExchangePB.MarketsRequest) ([]*spotExchangePB.SpotMarketInfo, error) {
	res, err := c.ExchangeClient.GetSpotMarkets(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Markets, nil
}

package internal

import (
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/chain"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
)

type Client struct {
	ExchangeClient *exchange.MbClient
	ChainClient    chain.ChainClient
}

func InitExchangeClient() *exchange.MbClient {
	exchangeClient := exchange.NewMbClient("local", configtypes.DefaultConfig())
	return exchangeClient
}

func NewClient() *Client {
	return &Client{
		ExchangeClient: InitExchangeClient(),
		ChainClient:    chain.NewChainClient("genesis"),
	}
}

func (c *Client) GetPrice(ticker string) (float64, error) {
	return c.ExchangeClient.GetPrice(ticker)
}

func (c *Client) GetBalances(address string, denoms []string) (map[string]float64, error) {
	balances := make(map[string]float64)
	for _, denom := range denoms {
		balance, err := c.ChainClient.GetBalance(address, denom)
		if err != nil {
			return nil, err
		}
		balances[denom] = balance
	}
	return balances, nil
}

func (c *Client) TransferToken(to string, amount float64, denom string) (string, error) {
	return c.ChainClient.TransferToken(to, amount, denom)
}

func (c *Client) GetAddress() string {
	return c.ChainClient.SenderAddress.String()
}

// This works for most of the tokens
func (c *Client) GetDecimal(denom string) int32 {
	if denom == "inj" {
		return 6
	}
	return 6
}

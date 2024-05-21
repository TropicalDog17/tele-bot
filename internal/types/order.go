package types

import (
	"fmt"

	"github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
)

type LimitOrderInfo struct {
	DenomIn   string  "json:\"denom_in\""
	DenomOut  string  "json:\"denom_out\""
	Amount    float64 "json:\"amount\""
	Price     float64 "json:\"price\""
	Direction string  "json:\"direction\""
	OrderHash string  "json:\"order_hash\""
	MarketID  string  "json:\"market_id\""
}

// NewLimitOrderInfo returns a new LimitOrderInfo with default values
func NewLimitOrderInfo() *LimitOrderInfo {
	// TODO: adjust default from settings
	defaultDenomIn := "atom"
	defaultDenomOut := "inj"
	ticker := defaultDenomIn + "/" + defaultDenomOut

	// mock private key for testing
	exchangeClient := exchange.NewMbClient("local", "", config.DefaultConfig())
	marketSummary, err := exchangeClient.GetMarketSummaryFromTicker(ticker)
	var defaultPrice float64
	if err != nil {
		defaultPrice = 0
	} else {
		defaultPrice = marketSummary.Price
	}
	fmt.Println("Default price: ", defaultPrice)
	return &LimitOrderInfo{
		DenomIn:   "atom",
		DenomOut:  "inj",
		Amount:    1,
		Price:     defaultPrice,
		Direction: "buy",
	}
}

package main

import (
	"fmt"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	configtypes "github.com/TropicalDog17/orderbook-go-sdk/config"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/exchange"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/types"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

func main() {
	exchangeClient := exchange.NewMbClient("local", "F7356863262652E2E9640F9641DBA58FB741998E5B8D838F14612473CC1DEFCC", configtypes.DefaultConfig())
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	marketSummary, err := exchangeClient.GetMarketSummary("0xfbd55f13641acbb6e69d7b59eb335dabe2ecbfea136082ce2eedaba8a0c917a3")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Market Summary: %+v\n", marketSummary)

	price, err := exchangeClient.GetPrice("ATOM/INJ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Price: %f\n", price)

	// Create a spot order
	// chainClient := chain.NewChainClientFromPrivateKey("6C212553111B370A8FFDC682954495B7B90A73CEDAB7106323646A4F2C4E668F")
	spotOrder := types.SpotOrder{
		OrderType: exchangetypes.OrderType_BUY,
		MarketId:  "0xfbd55f13641acbb6e69d7b59eb335dabe2ecbfea136082ce2eedaba8a0c917a3",
		Price:     decimal.NewFromFloat(0.70),
		Quantity:  decimal.NewFromFloat(0.01),
	}
	tx, err := exchangeClient.PlaceSpotOrder(spotOrder)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Transaction: %+v\n", tx)

}

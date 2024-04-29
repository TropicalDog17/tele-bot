package internal_test

import (
	"testing"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestDetermineOrderType(t *testing.T) {
	spotMarket := &exchangetypes.SpotMarket{
		MarketId:   "test",
		BaseDenom:  "atom",
		QuoteDenom: "inj",
	}
	// When want to buy atom in pair atom/inj, should be a buy order
	orderType := internal.DetermineOrderType(spotMarket, "atom", "inj")
	require.Equal(t, exchangetypes.OrderType_BUY, orderType)

	// When want to sell atom in pair atom/inj, should be a sell order
	orderType = internal.DetermineOrderType(spotMarket, "inj", "atom")
	require.Equal(t, exchangetypes.OrderType_SELL, orderType)

}

func TestGetPriceInUsd(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)
	client := internal.NewCoinGeckoClient()
	price, err := client.GetPriceInUsd("atom", "inj")
	require.NoError(t, err)
	require.Greater(t, price["atom"]["usd"], 0.0)
	require.Greater(t, price["inj"]["usd"], 0.0)
}

func TestFetchUsdPriceMap(t *testing.T) {
	err := godotenv.Load("../.env")
	require.NoError(t, err)
	client := internal.NewCoinGeckoClient()
	priceMap, err := client.FetchUsdPriceMap()
	require.NoError(t, err)
	require.Greater(t, priceMap["atom"], 0.0)
	require.Greater(t, priceMap["inj"], 0.0)

}

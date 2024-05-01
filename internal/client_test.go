package internal_test

import (
	"testing"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	"github.com/TropicalDog17/tele-bot/internal"
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

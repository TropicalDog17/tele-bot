package client

import (
	"testing"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	"github.com/stretchr/testify/require"
)

func TestDetermineOrderType(t *testing.T) {
	spotMarket := &exchangetypes.SpotMarket{
		MarketId:   "test",
		BaseDenom:  "atom",
		QuoteDenom: "inj",
	}
	// When want to buy atom in pair atom/inj, should be a buy order
	orderType := DetermineOrderType(spotMarket, "atom", "inj")
	require.Equal(t, exchangetypes.OrderType_BUY, orderType)

	// When want to sell atom in pair atom/inj, should be a sell order
	orderType = DetermineOrderType(spotMarket, "inj", "atom")
	require.Equal(t, exchangetypes.OrderType_SELL, orderType)
}

// func TestTransferToken(t *testing.T) {
// 	client := internal.NewClient("tropicaldog17")
// 	_, err := client.TransferToken("inj1l6377277h4nxap383hzphrf8ckecajdqzfx5uq", 0.0001, "inj")
// 	require.Error(t, err)
// }

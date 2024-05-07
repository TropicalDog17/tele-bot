package utils

import (
	"testing"

	"cosmossdk.io/math"
	decimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestQuantityToChainFormat(t *testing.T) {
	quantity := 0.0001
	baseDecimal := int32(18)
	quantityStr := QuantityToChainFormat(decimal.NewFromFloat(quantity), baseDecimal).String()
	require.Equal(t, "100000000000000.000000000000000000", quantityStr)
}

// Denom: "inj", Amount: math.Int(math.LegacyMustNewDecFromStr(amountStr))},

func TestConstructAmount(t *testing.T) {
	amount := 0.0001
	decimals := int32(18)
	amountDec := QuantityToChainFormat(decimal.NewFromFloat(amount), int32(decimals))
	amountInt := math.NewInt(amountDec.RoundInt64())
	require.Equal(t, int64(100000000000000), amountInt.Int64())

}

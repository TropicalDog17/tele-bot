package chain

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"
)

func (c *ChainClientStruct) GetBalance(address string, denom string) (float64, error) {
	denom = strings.ToLower(denom)
	var decimals int32
	switch denom {
	case "inj":
		decimals = 18
	case "eth":
		decimals = 18
	default:
		decimals = 6
	}
	ctx := context.Background()

	res, err := c.chainClient.GetBankBalance(ctx, address, denom)
	if err != nil {
		return 0, err
	}
	price := decimal.RequireFromString(res.Balance.Amount.String()).Mul(decimal.New(1, -decimals))
	priceFloat, _ := price.Float64()
	return priceFloat, nil

}

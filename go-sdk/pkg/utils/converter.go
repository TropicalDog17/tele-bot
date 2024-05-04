package utils

import (
	"cosmossdk.io/math"
	decimal "github.com/shopspring/decimal"
)

func PriceToChainFormat(price decimal.Decimal, baseDecimal int32, quoteDecimal int32) math.LegacyDec {
	price = price.Shift(quoteDecimal - baseDecimal)
	return math.LegacyMustNewDecFromStr(price.String())
}

func PriceFromChainFormat(price string, baseDecimal int32, quoteDecimal int32) decimal.Decimal {
	priceDecimal := decimal.RequireFromString(price)
	priceDecimal = priceDecimal.Shift(baseDecimal - quoteDecimal)
	return priceDecimal
}

func QuantityToChainFormat(quantity decimal.Decimal, baseDecimal int32) math.LegacyDec {
	quantity = quantity.Shift(baseDecimal)
	return math.LegacyMustNewDecFromStr(quantity.String())
}

func QuantityFromChainFormat(quantity string, baseDecimal int32) decimal.Decimal {
	quantityDecimal := decimal.RequireFromString(quantity)
	quantityDecimal = quantityDecimal.Shift(-baseDecimal)
	return quantityDecimal
}

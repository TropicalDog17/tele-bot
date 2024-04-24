package types

import "fmt"

type LimitOrderInfo struct {
	DenomIn   string  "json:\"denom_in\""
	DenomOut  string  "json:\"denom_out\""
	Amount    float64 "json:\"amount\""
	Price     float64 "json:\"price\""
	Direction string  "json:\"direction\""
}

// NewLimitOrderInfo returns a new LimitOrderInfo with default values
func NewLimitOrderInfo() *LimitOrderInfo {
	return &LimitOrderInfo{
		DenomIn:   "inj",
		DenomOut:  "atom",
		Amount:    1,
		Price:     3.6, // TODO: Should be able to fetch from the market
		Direction: "buy",
	}
}

// 📊 Limit Order - Buy
// ⬩ Mode: Buy
// ⬩ Token: SOL
// ⬩ Amount: 1.000000 SOL
// ⬩ Limit Price: $0.006382091 (0.00%)
// IN:   1.000000 SOL ($156.69)
// OUT: 156.688456 USDC ($156.69)
func (order *LimitOrderInfo) ToMessage() string {
	priceOut := 8.7
	priceIn := 28.5

	return fmt.Sprintf(`📊 Limit Order - %s
⬩ Mode: %s
⬩ TokenToPay: %s
⬩ Amount: %f %s
⬩ Limit Price: $%f (0.00%%)
IN:   %f %s ($%f)
OUT: %f %s ($%f)`, order.Direction, order.Direction, order.DenomOut, order.Amount, order.DenomIn, order.Price, order.Amount, order.DenomIn, order.Amount*priceIn, order.Amount*order.Price, order.DenomOut, order.Amount*order.Price*priceOut)
}

package types

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
		DenomIn:   "atom",
		DenomOut:  "inj",
		Amount:    1,
		Price:     0.3, // TODO: Should be able to fetch from the market
		Direction: "buy",
	}
}

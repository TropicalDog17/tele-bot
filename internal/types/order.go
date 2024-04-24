package types

type LimitOrderInfo struct {
	Denom     string  "json:\"denom\""
	Amount    float64 "json:\"amount\""
	Price     float64 "json:\"price\""
	Direction string  "json:\"direction\""
}

// NewLimitOrderInfo returns a new LimitOrderInfo with default values
func NewLimitOrderInfo() *LimitOrderInfo {
	return &LimitOrderInfo{
		Denom:     "atom",
		Amount:    1,
		Price:     0.23, // TODO: Should be able to fetch from the market
		Direction: "buy",
	}
}

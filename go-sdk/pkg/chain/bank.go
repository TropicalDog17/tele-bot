package chain

import (
	"fmt"

	"cosmossdk.io/math"
	utils "github.com/TropicalDog17/orderbook-go-sdk/pkg/utils"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/shopspring/decimal"
)

func (c *ChainClientStruct) TransferToken(toAddress string, amount float64, denom string) (string, error) {
	var decimals int32
	switch denom {
	case "inj":
		decimals = 18
	case "eth":
		decimals = 18
	default:
		decimals = 6
	}
	amountDec := utils.QuantityToChainFormat(decimal.NewFromFloat(amount), int32(decimals))
	// Trim trailing zeros

	// prepare tx msg
	msg := &banktypes.MsgSend{
		FromAddress: c.SenderAddress.String(),
		ToAddress:   toAddress,
		Amount: []sdktypes.Coin{{
			Denom: "inj", Amount: math.NewInt(amountDec.RoundInt64())},
		},
	}
	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg
	txResp, err := c.chainClient.SyncBroadcastMsg(msg)

	if err != nil {
		fmt.Println(err)
	}

	return txResp.TxResponse.TxHash, nil
}

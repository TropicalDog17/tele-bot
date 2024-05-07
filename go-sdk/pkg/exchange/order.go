package exchange

import (
	"context"
	"fmt"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	"github.com/TropicalDog17/orderbook-go-sdk/pkg/types"
	utils "github.com/TropicalDog17/orderbook-go-sdk/pkg/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderMaker interface {
	PlaceSpotOrder(order types.SpotOrder) error
	PlaceMarketOrder() error
}

func (c *MbClient) PlaceSpotOrder(order types.SpotOrder) (string, error) {
	chainClient := c.ChainClient.GetInjectiveChainClient()
	senderAddress := c.ChainClient.GetSenderAddress()
	ctx := context.Background()

	defaultSubaccountID := chainClient.DefaultSubaccount(senderAddress)
	baseDecimal, quoteDecimal := c.GetDecimals(ctx, order.MarketId)
	spotOrder := exchangetypes.SpotOrder{
		OrderType: exchangetypes.OrderType_BUY,
		MarketId:  order.MarketId,
		OrderInfo: exchangetypes.OrderInfo{
			SubaccountId: defaultSubaccountID.String(),
			Price:        utils.PriceToChainFormat(order.Price, baseDecimal, quoteDecimal),
			Quantity:     utils.QuantityToChainFormat(order.Quantity, baseDecimal),
			Cid:          uuid.NewString(),
		},
	}
	fmt.Println("spot order: ", spotOrder)
	msg := new(exchangetypes.MsgCreateSpotLimitOrder)
	msg.Sender = senderAddress.String()
	msg.Order = spotOrder
	simRes, err := chainClient.SimulateMsg(chainClient.ClientContext(), msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	msgCreateSpotLimitOrderResponse := exchangetypes.MsgCreateSpotLimitOrderResponse{}
	err = msgCreateSpotLimitOrderResponse.Unmarshal(simRes.Result.MsgResponses[0].Value)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg

	tx, err := chainClient.SyncBroadcastMsg(msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	txHash := tx.TxResponse.TxHash

	return txHash, nil
}

func (c *MbClient) NewSpotOrder(orderType exchangetypes.OrderType, marketId string, price float64, quantity float64) types.SpotOrder {
	return types.SpotOrder{
		OrderType: orderType,
		MarketId:  marketId,
		Price:     decimal.NewFromFloat32(float32(price)),
		Quantity:  decimal.NewFromFloat32(float32(quantity)),
	}
}

func (c *MbClient) CancelOrder(ctx context.Context, marketID, orderID string) (string, error) {
	chainClient := c.ChainClient.GetInjectiveChainClient()
	defaultSubaccountID := chainClient.DefaultSubaccount(c.ChainClient.GetSenderAddress())
	msg := &exchangetypes.MsgCancelSpotOrder{
		Sender:       c.ChainClient.GetSenderAddress().String(),
		MarketId:     marketID,
		SubaccountId: defaultSubaccountID.String(),
		OrderHash:    orderID,
	}
	simRes, err := chainClient.SimulateMsg(chainClient.ClientContext(), msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	msgCreateSpotLimitOrderResponse := exchangetypes.MsgCancelSpotOrderResponse{}
	err = msgCreateSpotLimitOrderResponse.Unmarshal(simRes.Result.MsgResponses[0].Value)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg

	tx, err := chainClient.SyncBroadcastMsg(msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	txHash := tx.TxResponse.TxHash

	return txHash, nil
}

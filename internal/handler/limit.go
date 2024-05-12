package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	types "github.com/TropicalDog17/tele-bot/internal/types"
	tele "gopkg.in/telebot.v3"
)

func HandleLimitOrder(b internal.Bot, authRoute *tele.Group, clients map[string]internal.BotClient, limitOrderMenu, createOrderMenu *tele.StoredMessage, globalLimitOrder *types.LimitOrderInfo, currentStep *string, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders *tele.ReplyMarkup) {
	authRoute.Handle(&types.BtnLimitOrder, func(c tele.Context) error {
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		rdb := client.GetRedisInstance()

		text := "📊 Limit Orders\n\nBuy or Sell tokens automatically at your desired price.\n1. Choose to Buy or Sell.\n2. Choose the Token to Buy or Sell.\n3. Select the amount to Buy or Sell.\n4. Set your target buy or sell price.\n5. Pick an expiry time for the order.\n6. Click Create Order and Review, Confirm.\n\nTo manage or view unfilled orders, click Active Orders."

		msg, err := b.Send(c.Chat(), text, types.MenuLimitOrder)
		if err != nil {
			return err
		}

		limitOrderMenu.ChatID = msg.Chat.ID
		limitOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)
		ctx := context.Background()
		err = rdb.HSet(ctx, "limitOrderMenu", "chatID", fmt.Sprintf("%d", limitOrderMenu.ChatID)).Err()
		if err != nil {
			return err
		}

		err = rdb.HSet(ctx, "limitOrderMenu", "messageID", limitOrderMenu.MessageID).Err()
		if err != nil {
			return err
		}

		return nil
	})
	authRoute.Handle(&types.BtnBuyLimitOrder, func(ctx tele.Context) error {
		text := "Place a buy limit order"
		client := clients[ctx.Callback().Sender.Username]
		rdb := client.GetRedisInstance()

		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(types.MenuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		msg, err := b.Send(ctx.Chat(), text, menuCreateLimitOrder)
		if err != nil {
			return err
		}

		createOrderMenu.ChatID = msg.Chat.ID
		createOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in Redis using HSET
		redisCtx := context.Background()
		err = rdb.HSet(redisCtx, "createOrderMenu", "chatID", fmt.Sprintf("%d", createOrderMenu.ChatID)).Err()
		if err != nil {
			return err
		}

		err = rdb.HSet(redisCtx, "createOrderMenu", "messageID", createOrderMenu.MessageID).Err()
		if err != nil {
			return err
		}

		// Adjust global limit order
		globalLimitOrder.Direction = "buy"

		return nil
	})
	authRoute.Handle(&types.BtnSellLimitOrder, func(ctx tele.Context) error {
		text := "Place a buy limit order"
		client := clients[ctx.Callback().Sender.Username]
		rdb := client.GetRedisInstance()
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(types.MenuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		msg, err := b.Send(ctx.Chat(), text, menuCreateLimitOrder)
		if err != nil {
			return err
		}

		createOrderMenu.ChatID = msg.Chat.ID
		createOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in Redis using HSET
		redisCtx := context.Background()
		err = rdb.HSet(redisCtx, "createOrderMenu", "chatID", fmt.Sprintf("%d", createOrderMenu.ChatID)).Err()
		if err != nil {
			return err
		}

		err = rdb.HSet(redisCtx, "createOrderMenu", "messageID", createOrderMenu.MessageID).Err()
		if err != nil {
			return err
		}

		// Adjust global limit order
		globalLimitOrder.Direction = "sell"

		return nil
	})
	authRoute.Handle(&types.BtnAmount, func(c tele.Context) error {
		*currentStep = "limitAmount"
		return c.Send("Enter the amount to buy", tele.ForceReply)
	})
	authRoute.Handle(&types.BtnPrice, func(c tele.Context) error {
		*currentStep = "limitPrice"
		return c.Send("Enter the price to buy", tele.ForceReply)
	})
	authRoute.Handle(&types.BtnToken, func(c tele.Context) error {
		*currentStep = "limitToken"
		return c.Send("Enter the token to buy", tele.ForceReply)
	})
	authRoute.Handle(&types.BtnConfirmOrder, func(c tele.Context) error {
		*currentStep = "confirmOrder"
		client := clients[c.Callback().Sender.Username]

		orderOverview := client.ToMessage(*globalLimitOrder, true)
		return c.Send(orderOverview, menuConfirmOrder)
	})
	authRoute.Handle(&types.BtnConfirmLimitOrder, func(c tele.Context) error {
		client := clients[c.Callback().Sender.Username]
		rdb := client.GetRedisInstance()
		txHash, err := client.PlaceSpotOrder(globalLimitOrder.DenomIn, globalLimitOrder.DenomOut, globalLimitOrder.Amount, globalLimitOrder.Price)
		if err != nil {
			return c.Send("Error placing limit order"+err.Error(), menu)
		}
		// Store the order in Redis
		orderJSON, err := json.Marshal(globalLimitOrder)
		if err != nil {
			return c.Send("Error placing limit order"+err.Error(), menu)
		}
		ctx := context.Background()
		err = rdb.HSet(ctx, client.GetAddress(), globalLimitOrder.OrderHash, string(orderJSON)).Err()
		if err != nil {
			return c.Send("Error placing limit order"+err.Error(), menu)
		}
		explorer := fmt.Sprintf("https://0267-2402-800-61c5-784e-4ddd-972f-73a6-f0a2.ngrok-free.app/injective/tx/%s", txHash)
		button := tele.InlineButton{
			Text: "View Transaction",
			URL:  explorer,
		}
		menu := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{button},
			},
		}
		text := "Transaction sent successfully."
		return c.Send(text, menu)
	})
	// Handle active orders
	authRoute.Handle(&types.BtnActiveOrders, func(c tele.Context) error {
		client, ok := clients[c.Callback().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		// markets, err := client.
		ctx := context.Background()
		msgs := []string{}
		orders := []types.LimitOrderInfo{}
		marketOrders, err := client.GetRedisInstance().HGetAll(ctx, client.GetAddress()).Result()
		if err != nil {
			return c.Send("Error fetching active orders, please try again later")
		}
		for _, orderJSON := range marketOrders {
			order := types.LimitOrderInfo{}
			err = json.Unmarshal([]byte(orderJSON), &order)
			if err != nil {
				return c.Send("Error fetching active orders")
			}
			orders = append(orders, order)
		}

		if len(orders) == 0 {
			return c.Send("No active orders", menu)
		}
		for _, order := range orders {
			msgs = append(msgs, client.ToMessage(order, false))
		}

		return c.Send(strings.Join(msgs, "\n\n"), menuActiveOrders)
	})
	b.Handle(&types.BtnCancelOrder, func(c tele.Context) error {
		*currentStep = "cancelOrder"
		return c.Send("Enter the order id to cancel", tele.ForceReply)
	})
}

func HandleLimitStep(b *tele.Bot, c tele.Context, createOrderMenu *tele.StoredMessage, menuLimitOrder, menuCreateLimitOrder *tele.ReplyMarkup, globalLimitOrder *types.LimitOrderInfo, step *string) error {
	switch *step {
	case "limitAmount":
		amount, err := strconv.ParseFloat(c.Text(), 64)
		if err != nil {
			return c.Send("Invalid amount")
		}
		globalLimitOrder.Amount = amount
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		return internal.DeleteInputMessage(b, c)
	case "limitPrice":
		price, err := strconv.ParseFloat(c.Text(), 64)
		if err != nil {
			return c.Send("Invalid price")
		}
		globalLimitOrder.Price = price
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		return internal.DeleteInputMessage(b, c)
	case "limitToken":
		globalLimitOrder.DenomOut = c.Text()
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err := b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		return internal.DeleteInputMessage(b, c)
	}
	return c.Send("Invalid input")
}

func HandleCancelLimitOrderStep(b *tele.Bot, c tele.Context, client internal.BotClient, globalLimitOrder *types.LimitOrderInfo) error {
	orderId := c.Text()
	marketId, err := client.GetRedisInstance().HGet(context.Background(), "orders", orderId).Result()
	if err != nil {
		return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
	}
	txhash, err := client.CancelOrder(marketId, orderId)

	if err != nil {
		return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
	}
	// Delete the order from the database
	err = client.GetRedisInstance().HDel(context.Background(), client.GetAddress(), orderId).Err()
	if err != nil {
		return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
	}
	return c.Send(fmt.Sprintf("Order cancelled with tx hash: %s", txhash), types.MenuLimitOrder, types.Menu)
}

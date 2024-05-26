package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TropicalDog17/tele-bot/internal"
	types "github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func HandleLimitOrder(b internal.Bot, localizer *i18n.Localizer, authRoute *tele.Group, clients map[string]internal.BotClient, limitOrderMenu, createOrderMenu *tele.StoredMessage, globalLimitOrder *types.LimitOrderInfo, btnLimitOrder, btnBuyLimitOrder, btnSellLimitOrder, btnAmount, btnPrice, btnToken, btnPayWith, btnConfirmOrder, btnConfirmLimitOrder, btnActiveOrder, btnCancelOrder tele.Btn, currentStep *string, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders *tele.ReplyMarkup) {
	authRoute.Handle(&btnLimitOrder, func(c tele.Context) error {
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		rdb := client.GetRedisInstance()

		// text := "📊 Limit Orders\n\nBuy or Sell tokens automatically at your desired price.\n1. Choose to Buy or Sell.\n2. Choose the Token to Buy or Sell.\n3. Select the amount to Buy or Sell.\n4. Set your target buy or sell price.\n5. Click Create Order and Review, Confirm.\n\nTo manage or view unfilled orders, click Active Orders."
		text := localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "LimitOrderText", Other: "📊 Limit Orders\n\nBuy or Sell tokens automatically at your desired price.\n1. Choose to Buy or Sell.\n2. Choose the Token to Buy or Sell.\n3. Select the amount to Buy or Sell.\n4. Set your target buy or sell price.\n5. Click Create Order and Review, Confirm.\n\nTo manage or view unfilled orders, click Active Orders."}})
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
	authRoute.Handle(&btnBuyLimitOrder, func(ctx tele.Context) error {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		text := "Place a buy limit order\n Updated at: " + currentTime
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
	authRoute.Handle(&btnSellLimitOrder, func(ctx tele.Context) error {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		text := "Place a sell limit order\nDefault price updated at: " + currentTime
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
	authRoute.Handle(&btnAmount, func(c tele.Context) error {
		*currentStep = "limitAmount"
		return c.Send("Enter the amount to buy", tele.ForceReply)
	})
	authRoute.Handle(&btnPrice, func(c tele.Context) error {
		*currentStep = "limitPrice"
		return c.Send("Enter the price to buy", tele.ForceReply)
	})
	authRoute.Handle(&btnToken, func(c tele.Context) error {
		*currentStep = "limitToken"
		return c.Send("Enter the token to buy", tele.ForceReply)
	})
	authRoute.Handle(&btnPayWith, func(c tele.Context) error {
		*currentStep = "payWithToken"
		return c.Send("Enter the token to pay with", tele.ForceReply)
	})
	authRoute.Handle(&btnConfirmOrder, func(c tele.Context) error {
		*currentStep = "confirmOrder"
		client := clients[c.Callback().Sender.Username]

		orderOverview := client.ToMessage(*globalLimitOrder, true)
		return c.Send(orderOverview, menuConfirmOrder)
	})
	authRoute.Handle(&btnConfirmLimitOrder, func(c tele.Context) error {
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
		explorerUrl := os.Getenv("EXPLORER_URL")
		explorer := fmt.Sprintf("%s/injective/tx/%s", explorerUrl, txHash)
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
	authRoute.Handle(&btnActiveOrder, func(c tele.Context) error {
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
			return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "NoActiveOrders", Other: "No active orders"}}), types.Menu)
		}
		for _, order := range orders {
			msgs = append(msgs, client.ToMessage(order, false))
		}

		return c.Send(strings.Join(msgs, "\n\n"), menuActiveOrders)
	})
	b.Handle(&btnCancelOrder, func(c tele.Context) error {
		*currentStep = "cancelOrder"
		return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "EnterCancelOrderId", Other: "Enter the order ID to cancel"}}), types.Menu)
	})
}

func HandleLimitStep(b *tele.Bot, c tele.Context, client internal.BotClient, createOrderMenu *tele.StoredMessage, menuLimitOrder, menuCreateLimitOrder *tele.ReplyMarkup, globalLimitOrder *types.LimitOrderInfo, step *string) error {
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
		fmt.Printf("%+v", globalLimitOrder)
		globalLimitOrder.Price = price
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		return internal.DeleteInputMessage(b, c)
	case "limitToken":
		if globalLimitOrder.Direction == "buy" {
			globalLimitOrder.DenomIn = strings.ToLower(c.Text())
		} else {
			globalLimitOrder.DenomOut = strings.ToLower(c.Text())
		}
		market, err := client.GetExchangeClient().GetMarketSummaryFromTicker(globalLimitOrder.DenomIn + "/" + globalLimitOrder.DenomOut)
		if err == nil {
			globalLimitOrder.Price = market.Price
		}

		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		return internal.DeleteInputMessage(b, c)
	case "payWithToken":
		if globalLimitOrder.Direction == "buy" {
			globalLimitOrder.DenomOut = strings.ToLower(c.Text())
		} else {
			globalLimitOrder.DenomIn = strings.ToLower(c.Text())
		}
		fmt.Printf("%+v", globalLimitOrder)

		market, err := client.GetExchangeClient().GetMarketSummaryFromTicker(globalLimitOrder.DenomIn + "/" + globalLimitOrder.DenomOut)
		if err == nil {
			fmt.Println("fetched price is: ", market.Price)
			globalLimitOrder.Price = market.Price
		} else {
			fmt.Println("Error getting market summary: ", err)
		}

		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
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

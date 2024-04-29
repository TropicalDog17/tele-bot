package handler

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	types "github.com/TropicalDog17/tele-bot/internal/types"
	tele "gopkg.in/telebot.v3"
)

func HandleLimitOrder(b *tele.Bot, client *internal.Client, limitOrderMenu, createOrderMenu *tele.StoredMessage, globalLimitOrder *types.LimitOrderInfo, currentStep *string, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders *tele.ReplyMarkup) {
	b.Handle(&types.BtnLimitOrder, func(c tele.Context) error {
		text := "📊 Limit Orders\n\nBuy or Sell tokens automatically at your desired price.\n1. Choose to Buy or Sell.\n2. Choose the Token to Buy or Sell.\n3. Select the amount to Buy or Sell.\n4. Set your target buy or sell price.\n5. Pick an expiry time for the order.\n6. Click Create Order and Review, Confirm.\n\nTo manage or view unfilled orders, click Active Orders."
		msg, err := b.Send(c.Chat(), text, types.MenuLimitOrder)
		if err != nil {
			return err
		}
		limitOrderMenu.ChatID = msg.Chat.ID
		limitOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/limitOrderMenu.txt", []byte(fmt.Sprintf("%d %s", limitOrderMenu.ChatID, limitOrderMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}
		return nil
	})
	b.Handle(&types.BtnBuyLimitOrder, func(ctx tele.Context) error {
		text := "Place a buy limit order"
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(types.MenuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		msg, err := b.Send(ctx.Chat(), text, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		createOrderMenu.ChatID = msg.Chat.ID
		createOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)
		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/createOrderMenu.txt", []byte(fmt.Sprintf("%d %s", createOrderMenu.ChatID, createOrderMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}
		// Adjust global limit order
		globalLimitOrder.Direction = "buy"

		return nil
	})
	b.Handle(&types.BtnSellLimitOrder, func(ctx tele.Context) error {
		text := "Place a sell limit order"
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		msg, err := b.Send(ctx.Chat(), text, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		createOrderMenu.ChatID = msg.Chat.ID
		createOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)
		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/createOrderMenu.txt", []byte(fmt.Sprintf("%d %s", createOrderMenu.ChatID, createOrderMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}
		// Adjust global limit order
		globalLimitOrder.Direction = "sell"

		return nil
	})
	b.Handle(&types.BtnAmount, func(c tele.Context) error {
		*currentStep = "limitAmount"
		return c.Send("Enter the amount to buy", tele.ForceReply)
	})
	b.Handle(&types.BtnPrice, func(c tele.Context) error {
		*currentStep = "limitPrice"
		return c.Send("Enter the price to buy", tele.ForceReply)
	})
	b.Handle(&types.BtnToken, func(c tele.Context) error {
		*currentStep = "limitToken"
		return c.Send("Enter the token to buy", tele.ForceReply)
	})
	b.Handle(&types.BtnConfirmOrder, func(c tele.Context) error {
		*currentStep = "confirmOrder"
		orderOverview := client.ToMessage(*globalLimitOrder, true)
		return c.Send(orderOverview, menuConfirmOrder)
	})
	b.Handle(&types.BtnConfirmLimitOrder, func(c tele.Context) error {
		// TODO: Perform the limit order logic here
		txHash, err := client.PlaceSpotOrder(globalLimitOrder.DenomIn, globalLimitOrder.DenomOut, globalLimitOrder.Amount, globalLimitOrder.Price)
		if err != nil {
			return c.Send("Error placing limit order", menu)
		}
		return c.Send("Successfully send order, check txhash here: "+txHash, menu)
	})
	// Handle active orders
	b.Handle(&types.BtnActiveOrders, func(c tele.Context) error {
		// TODO: fix hardcoded market id, should fetch from all markets
		orders, err := client.GetActiveOrders("0xfbd55f13641acbb6e69d7b59eb335dabe2ecbfea136082ce2eedaba8a0c917a3")
		if err != nil {
			return c.Send("Error fetching active orders")
		}
		msgs := []string{}
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

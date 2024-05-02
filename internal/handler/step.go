package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	tele "gopkg.in/telebot.v3"
)

func HandleStep(b *tele.Bot, client *internal.Client, currentStep *string, menuSendToken, menuLimitOrder, menuCreateLimitOrder *tele.ReplyMarkup, globalLimitOrder *types.LimitOrderInfo, selectedAmount, selectedToken, recipientAddress *string, globalMenu *tele.StoredMessage, createOrderMenu *tele.StoredMessage) {
	b.Handle(tele.OnText, func(c tele.Context) error {
		// Check if the user is entering a custom amount
		if *currentStep == "customAmount" {
			*selectedAmount = c.Text()
			menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
			return c.Send(fmt.Sprintf("Selected amount: %s", *selectedAmount), menuSendToken)
		} else if *currentStep == "recipientAddress" { // Check if the user is entering a recipient addres
			*recipientAddress = c.Text()
			fmt.Println("Recipient address: ", recipientAddress)
			err := b.Delete(c.Message().ReplyTo)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
			types.BtnRecipientSection.Text = "Recipient:" + *recipientAddress
			menuSendToken.InlineKeyboard[6][0] = *types.BtnRecipientSection.Inline()

			// load the global menu from database
			redisCtx := context.Background()
			sendTokenMenu, err := client.GetRedisInstance().HGetAll(redisCtx, "sendTokenMenu").Result()
			if err != nil {
				return err
			}
			globalMenu.ChatID, err = strconv.ParseInt(sendTokenMenu["chatID"], 10, 64)
			if err != nil {
				return err
			}
			globalMenu.MessageID = sendTokenMenu["messageID"]
			_, err = b.EditReplyMarkup(globalMenu, menuSendToken)
			if err != nil {
				return err
			}
			return nil
		} else if *currentStep == "limitAmount" {
			amount, err := strconv.ParseFloat(c.Text(), 64)
			if err != nil {
				return c.Send("Invalid amount")
			}
			globalLimitOrder.Amount = amount
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		} else if *currentStep == "limitPrice" {
			price, err := strconv.ParseFloat(c.Text(), 64)
			if err != nil {
				return c.Send("Invalid price")
			}
			globalLimitOrder.Price = price
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err = b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		} else if *currentStep == "limitToken" {
			globalLimitOrder.DenomOut = c.Text()
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err := b.EditReplyMarkup(createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		} else if *currentStep == "cancelOrder" {
			orderId := c.Text()
			marketId, err := client.GetRedisInstance().HGet(context.Background(), "orders", orderId).Result()
			if err != nil {
				return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
			}
			txhash, err := client.CancelOrder(marketId, orderId)

			if err != nil {
				return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
			}
			err = client.GetRedisInstance().HDel(context.Background(), client.GetAddress(), orderId).Err()
			if err != nil {
				return c.Send(fmt.Sprintf("Error cancelling order: %s", err), types.Menu)
			}
			*currentStep = ""
			return c.Send(fmt.Sprintf("Order cancelled with tx hash: %s", txhash), types.Menu)
		}
		return c.Send("Invalid input", types.Menu)

	})

}

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
		} else if *currentStep == "limitAmount" || *currentStep == "limitPrice" || *currentStep == "limitToken" {
			return HandleLimitStep(b, c, createOrderMenu, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, *currentStep)
		} else if *currentStep == "cancelOrder" {
			return HandleCancelLimitOrderStep(b, c, client, globalLimitOrder)
		}
		return c.Send("Invalid input", types.Menu)
	})

}

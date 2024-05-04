package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	tele "gopkg.in/telebot.v3"
)

func HandlerTransferToken(b *tele.Bot, client *internal.Client, menuSendToken *tele.ReplyMarkup, btnInlineAtom, btnInlineInj, btnTenDollar, btnFiftyDollar, btnHundredDollar, btnTwoHundredDollar, btnFiveHundredDollar, btnCustomAmount, btnRecipientSection, btnCustomToken *tele.Btn, selectedToken, selectedAmount, currentStep, recipientAddress *string, globalMenu *tele.StoredMessage) {
	// Handle the "Send Tokens" button click
	rdb := client.GetRedisInstance()
	b.Handle(&types.BtnSendToken, func(c tele.Context) error {
		msg, err := b.Send(c.Chat(), "Select the token to send:", menuSendToken)
		if err != nil {
			return err
		}

		globalMenu.ChatID = msg.Chat.ID
		globalMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in Redis using HSET
		ctx := context.Background()
		err = rdb.HSet(ctx, "sendTokenMenu", "chatID", fmt.Sprintf("%d", globalMenu.ChatID)).Err()
		if err != nil {
			return err
		}

		err = rdb.HSet(ctx, "sendTokenMenu", "messageID", globalMenu.MessageID).Err()
		if err != nil {
			return err
		}

		return nil
	})
	// Handle inline button clicks for token selection
	b.Handle(btnInlineAtom, func(c tele.Context) error {
		*selectedToken = "atom"
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][0] = internal.AddGreenTick(*btnInlineAtom.Inline())
		return c.Edit("Selected token: ATOM", menuSendToken)
	})

	b.Handle(btnInlineInj, func(c tele.Context) error {
		*selectedToken = "inj"
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][1] = internal.AddGreenTick(*btnInlineInj.Inline())
		return c.Edit("Selected token: INJ", menuSendToken)
	})

	// Handle amount button clicks
	b.Handle(btnTenDollar, func(c tele.Context) error {
		*selectedAmount = "1"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][0] = internal.AddGreenTick(*btnTenDollar.Inline())
		return c.Edit("Selected amount: 10 "+strings.ToUpper(*selectedToken), menuSendToken)
	})

	b.Handle(btnFiftyDollar, func(c tele.Context) error {
		*selectedAmount = "5"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][1] = internal.AddGreenTick(*btnFiftyDollar.Inline())

		return c.Edit("Selected amount: 50 "+strings.ToUpper(*selectedToken), menuSendToken)
	})

	b.Handle(btnHundredDollar, func(c tele.Context) error {
		*selectedAmount = "10"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][2] = internal.AddGreenTick(*btnHundredDollar.Inline())
		return c.Edit("Selected amount: 100 "+strings.ToUpper(*selectedToken), menuSendToken)
	})

	b.Handle(btnTwoHundredDollar, func(c tele.Context) error {
		*selectedAmount = "50"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][0] = internal.AddGreenTick(*btnTwoHundredDollar.Inline())
		return c.Edit("Selected amount: 50 "+strings.ToUpper(*selectedToken), menuSendToken)
	})

	b.Handle(btnFiveHundredDollar, func(c tele.Context) error {
		*selectedAmount = "100"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][1] = internal.AddGreenTick(*btnFiveHundredDollar.Inline())
		return c.Edit("Selected amount: 100 "+strings.ToUpper(*selectedToken), menuSendToken)
	})

	b.Handle(btnCustomAmount, func(c tele.Context) error {
		// Prompt the user to enter a custom amount
		*currentStep = "customAmount"
		return c.Send("Please enter the custom amount:")
	})
	b.Handle(btnCustomToken, func(c tele.Context) error {
		// Prompt the user to enter a custom token
		*currentStep = "customToken"
		return c.Send("Please enter the custom token:")
	})

	b.Handle(btnRecipientSection, func(c tele.Context) error {
		// Prompt the user to enter a recipient address
		*currentStep = "recipientAddress"
		return c.Send("Please enter the recipient address:", tele.ForceReply)
	})
	// Handle the "Send" button click
	b.Handle(&types.BtnSend, func(c tele.Context) error {
		// Sanity check to ensure all required fields are filled
		if *selectedToken == "" || *selectedAmount == "" || *recipientAddress == "" {
			return c.Send("Please fill in all required fields", menuSendToken)
		}
		selectedAmount, err := strconv.ParseFloat(*selectedAmount, 64)
		fmt.Println(selectedAmount)
		if err != nil {
			return c.Send("Invalid amount", menuSendToken)
		}
		// Trim whitespace from the recipient address
		*recipientAddress = strings.TrimSpace(*recipientAddress)
		txHash, err := client.TransferToken(*recipientAddress, selectedAmount/100, *selectedToken)
		if err != nil {
			return c.Send("Error sending token", menuSendToken)
		}

		// TODO: Perform the token sending logic here
		// Use the selected token, amount, and recipient address
		return c.Send(fmt.Sprintf("Sent %f %s to %s, with tx hash %s", selectedAmount, *selectedToken, *recipientAddress, txHash), types.Menu)
	})
}

func HandleTransferStep(b *tele.Bot, client *internal.Client, c tele.Context, menuSendToken *tele.ReplyMarkup, selectedAmount, selectedToken, recipientAddress *string, globalMenu *tele.StoredMessage, currentStep *string) error {
	switch *currentStep {
	case "customAmount":
		*selectedAmount = c.Text()
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, *selectedAmount, *selectedToken)
		return c.Send(fmt.Sprintf("Selected amount: %s", *selectedAmount), menuSendToken)
	case "recipientAddress":
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
	}
	return c.Send("Invalid input", types.Menu)
}

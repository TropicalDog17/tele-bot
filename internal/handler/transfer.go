package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func HandlerTransferToken(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, clients map[string]internal.BotClient, menuSendToken *tele.ReplyMarkup, btnSend, btnSendToken, btnInlineAtom, btnInlineInj, btnTenDollar, btnFiftyDollar, btnHundredDollar, btnTwoHundredDollar, btnFiveHundredDollar, btnCustomAmount, btnRecipientSection, btnCustomToken tele.Btn, transferInfo *types.TransferInfo, currentStep *string, globalMenu *tele.StoredMessage) {
	// Handle the "Send Tokens" button click
	authRoute.Handle(&btnSendToken, func(c tele.Context) error {
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		rdb := client.GetRedisInstance()
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
	authRoute.Handle(&btnInlineAtom, func(c tele.Context) error {
		transferInfo.SelectedToken = "atom"
		fmt.Println("btnInlineAtom clicked")
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][0] = internal.AddGreenTick(*btnInlineAtom.Inline())
		return c.Edit("Selected token: ATOM", menuSendToken)
	})

	authRoute.Handle(&btnInlineInj, func(c tele.Context) error {
		transferInfo.SelectedToken = "inj"
		fmt.Println("btnInlineInj clicked")
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][1] = internal.AddGreenTick(*btnInlineInj.Inline())
		return c.Edit("Selected token: INJ", menuSendToken)
	})

	// Handle amount button clicks
	authRoute.Handle(&btnTenDollar, func(c tele.Context) error {
		transferInfo.SelectedAmount = "1"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][0] = internal.AddGreenTick(*btnTenDollar.Inline())
		return c.Edit("Selected amount: 1 "+strings.ToUpper(transferInfo.SelectedToken), menuSendToken)
	})

	authRoute.Handle(&btnFiftyDollar, func(c tele.Context) error {
		transferInfo.SelectedAmount = "5"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][1] = internal.AddGreenTick(*btnFiftyDollar.Inline())

		return c.Edit("Selected amount: 5 "+strings.ToUpper(transferInfo.SelectedToken), menuSendToken)
	})

	authRoute.Handle(&btnHundredDollar, func(c tele.Context) error {
		transferInfo.SelectedAmount = "10"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][2] = internal.AddGreenTick(*btnHundredDollar.Inline())
		return c.Edit("Selected amount: 10 "+strings.ToUpper(transferInfo.SelectedToken), menuSendToken)
	})

	authRoute.Handle(&btnTwoHundredDollar, func(c tele.Context) error {
		transferInfo.SelectedAmount = "20"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][0] = internal.AddGreenTick(*btnTwoHundredDollar.Inline())
		return c.Edit("Selected amount: 20 "+strings.ToUpper(transferInfo.SelectedToken), menuSendToken)
	})

	authRoute.Handle(&btnFiveHundredDollar, func(c tele.Context) error {
		transferInfo.SelectedAmount = "50"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][1] = internal.AddGreenTick(*btnFiveHundredDollar.Inline())
		return c.Edit("Selected amount: 50 "+strings.ToUpper(transferInfo.SelectedToken), menuSendToken)
	})

	authRoute.Handle(&btnCustomAmount, func(c tele.Context) error {
		// Prompt the user to enter a custom amount
		*currentStep = "customAmount"
		return c.Send("Please enter the custom amount:")
	})
	authRoute.Handle(&btnCustomToken, func(c tele.Context) error {
		// Prompt the user to enter a custom token
		*currentStep = "customToken"
		return c.Send("Please enter the custom token:")
	})

	authRoute.Handle(&btnRecipientSection, func(c tele.Context) error {
		// Prompt the user to enter a recipient address
		*currentStep = "recipientAddress"
		return c.Send("Please enter the recipient address:", tele.ForceReply)
	})
	// Handle the "Send" button click
	authRoute.Handle(&btnSend, func(c tele.Context) error {
		client, ok := clients[c.Callback().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		// Sanity check to ensure all required fields are filled
		if transferInfo.SelectedToken == "" || transferInfo.SelectedAmount == "" || transferInfo.RecipientAddress == "" {
			missingField := ""
			if transferInfo.SelectedToken == "" {
				missingField = "token"
			} else if transferInfo.SelectedAmount == "" {
				missingField = "amount"
			} else if transferInfo.RecipientAddress == "" {
				missingField = "recipient address"
			}
			return c.Send("Please fill in all required fields, missing "+missingField, menuSendToken)
		}
		selectedAmount, err := strconv.ParseFloat(transferInfo.SelectedAmount, 64)
		if err != nil {
			return c.Send("Invalid amount", menuSendToken)
		}
		// Trim whitespace from the recipient address
		transferInfo.RecipientAddress = strings.TrimSpace(transferInfo.RecipientAddress)
		txHash, err := client.TransferToken(transferInfo.RecipientAddress, selectedAmount/100, transferInfo.SelectedToken)
		if err != nil {
			return c.Send("Error sending token", menuSendToken)
		}

		// TODO: Perform the token sending logic here
		// Use the selected token, amount, and recipient address
		return c.Send(fmt.Sprintf("Sent %f %s to %s, with tx hash %s", selectedAmount, transferInfo.SelectedToken, transferInfo.RecipientAddress, txHash), types.Menu)
	})
}

func HandleTransferStep(b *tele.Bot, localizer *i18n.Localizer, client internal.BotClient, c tele.Context, menuSendToken *tele.ReplyMarkup, transferInfo *types.TransferInfo, globalMenu *tele.StoredMessage, currentStep *string) error {
	switch *currentStep {
	case "customAmount":
		transferInfo.SelectedAmount = c.Text()
		// menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, transferInfo.SelectedAmount, transferInfo.SelectedToken)
		return c.Send(fmt.Sprintf("Selected amount: %s", transferInfo.SelectedAmount), menuSendToken)
	case "recipientAddress":
		transferInfo.RecipientAddress = c.Text()
		fmt.Println("Recipient address: ", transferInfo.RecipientAddress)
		err := b.Delete(c.Message().ReplyTo)
		if err != nil {
			return err
		}
		err = c.Delete()
		if err != nil {
			return err
		}
		btnRecipientSection := types.BtnRecipientSection(localizer, transferInfo)
		if transferInfo.RecipientAddress == "" {
			btnRecipientSection.Text = "Recipient: Not set"
		} else {
			btnRecipientSection.Text = "Recipient: " + transferInfo.RecipientAddress
		}
		menuSendToken.InlineKeyboard[6][0] = *btnRecipientSection.Inline()
		return c.Edit("Recipent set: "+transferInfo.RecipientAddress, menuSendToken)
	case "customToken":
		fmt.Println("Custom token: ", c.Text())
		transferInfo.SelectedToken = strings.ToLower(c.Text())
		if transferInfo.SelectedToken != "atom" && transferInfo.SelectedToken != "inj" && transferInfo.SelectedToken != "usdt" {
			return c.Send("Unsupported token", menuSendToken)
		}
		// menuSendToken.InlineKeyboard = internal.ModifyCustomTokenButton(menuSendToken.InlineKeyboard, transferInfo.SelectedToken)
		return c.Send(fmt.Sprintf("Selected token: %s", transferInfo.SelectedToken), menuSendToken)
	}

	return c.Send("Invalid input", types.Menu)
}

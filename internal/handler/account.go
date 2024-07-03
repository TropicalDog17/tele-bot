package handler

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	types "github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/skip2/go-qrcode"
	tele "gopkg.in/telebot.v3"
)

func HandleAddressQr(b *tele.Bot, authRoute *tele.Group, clients map[string]internal.BotClient) {
	authRoute.Handle(&tele.Btn{Unique: "qr"}, func(c tele.Context) error {
		client, ok := clients[c.Callback().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		// Get the address
		address := client.GetAddress()

		// Generate the QR code in memory
		qr, err := qrcode.New(address, qrcode.Medium)
		if err != nil {
			fmt.Printf("could not generate QRCode: %v", err)
			return err
		}

		// Convert the QR code to a byte slice
		qrBytes, err := qr.PNG(256)
		if err != nil {
			fmt.Printf("could not convert QRCode to PNG: %v", err)
			return err
		}

		// Create a byte reader from the QR code byte slice
		reader := bytes.NewReader(qrBytes)

		// Send the QR code directly from memory
		photo := &tele.Photo{File: tele.FromReader(reader)}
		_, err = b.Send(c.Chat(), photo)
		if err != nil {
			return err
		}

		return nil
	})
}

func HandleAccountDetails(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, clients map[string]internal.BotClient, btnShowAccount tele.Btn) {
	// Show account
	authRoute.Handle(&btnShowAccount, func(c tele.Context) error {
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		accountDetails := &tele.ReplyMarkup{}
		address := client.GetAddress()
		balances, err := client.GetBalances(address, []string{"atom", "inj"})
		if err != nil {
			return c.Send("Error fetching balances" + err.Error())
		}
		rows := []tele.Row{}
		totalBalanceInUsd := 0.0
		for denom, balance := range balances {
			usdPrice, found := client.GetPrice(denom)
			var balanceInUsd float64
			if !found {
				balanceInUsd = 0
			} else {
				balanceInUsd = balance * usdPrice
			}
			totalBalanceInUsd += balanceInUsd
			rows = append(rows, accountDetails.Row(accountDetails.Data(fmt.Sprintf("💰💰 %s: %.3f(%.3f $)", strings.ToUpper(denom), balance, balanceInUsd), "balance", "balance")))
		}
		rows = append(rows, accountDetails.Row(accountDetails.Data(fmt.Sprintf("💸💸 Total Balance: %.3f $ 💸💸", totalBalanceInUsd), "totalBalance", "totalBalance")))
		rows = append(rows, accountDetails.Row(accountDetails.Data("📱 Show QR for address", "qr", "qr")))
		accountDetails.Inline(rows...)
		// Message contain the account address
		explorerUrl := os.Getenv("EXPLORER_URL")
		message := fmt.Sprintf("Account: [%s](%s/injective/account/%s)", address, explorerUrl, address)
		return c.Send(message, accountDetails, tele.ModeMarkdown)
	})
}

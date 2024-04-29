package handler

import (
	"fmt"
	"os"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	tele "gopkg.in/telebot.v3"
)

func HandleAddressQr(b *tele.Bot, client *internal.Client) {
	b.Handle(&tele.Btn{Unique: "qr"}, func(c tele.Context) error {
		// Get the address
		address := client.GetAddress()
		// Address
		qrc, err := qrcode.New(address)
		if err != nil {
			fmt.Printf("could not generate QRCode: %v", err)
			return err
		}

		w, err := standard.New("../temp/repo-qrcode.jpeg")
		if err != nil {
			fmt.Printf("standard.New failed: %v", err)
			return err
		}
		defer w.Close()

		// save file
		if err = qrc.Save(w); err != nil {
			fmt.Printf("could not save image: %v", err)
		}
		photo := &tele.Photo{File: tele.FromDisk("../temp/repo-qrcode.jpeg")}
		_, err = b.Send(c.Chat(), photo)
		if err != nil {
			return err
		}
		// Remove the file after sending
		return os.Remove("../temp/repo-qrcode.jpeg")
	})
}

func HandleAccountDetails(b *tele.Bot, client *internal.Client) {
	// Show account
	b.Handle(&types.BtnShowAccount, func(c tele.Context) error {
		accountDetails := &tele.ReplyMarkup{}
		address := client.GetAddress()
		balances, err := client.GetBalances(address, []string{"atom", "inj"})
		if err != nil {
			return c.Send("Error fetching balances")
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
			rows = append(rows, accountDetails.Row(accountDetails.Data(fmt.Sprintf("%s: %.3f(%.3f $)", denom, balance, balanceInUsd), "balance", "balance")))
		}
		rows = append(rows, accountDetails.Row(accountDetails.Data(fmt.Sprintf("Total Balance: %.3f $", totalBalanceInUsd), "totalBalance", "totalBalance")))
		rows = append(rows, accountDetails.Row(accountDetails.Data("Show QR for address", "qr", "qr")))
		accountDetails.Inline(rows...)
		// Message contain the account address

		return c.Send("Account: "+address, accountDetails)
	})
}

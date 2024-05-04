package handler

import (
	"github.com/TropicalDog17/tele-bot/internal"
	tele "gopkg.in/telebot.v3"
)

var btnOnboard = tele.ReplyMarkup{}
var btnConfirmMnemonic = btnOnboard.Data("Confirm Mnemonic", "confirm_mnemonic")

func HandleOnboard(b internal.Bot, client internal.BotClient) {
	btnOnboard.Inline(
		btnOnboard.Row(btnConfirmMnemonic),
	)
	b.Handle("/start", func(c tele.Context) error {
		text := "Welcome to the TropicalDog17 bot! 🐶\n\nI am a bot that can help you with your trading needs. I can provide you with the latest cryptocurrency prices, help you place limit orders, and more.\n\nTo get started, type /help to see a list of available commands."
		return c.Reply(text, &btnOnboard, tele.ModeHTML)
	})
}

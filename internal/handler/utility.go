package handler

import (
	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	tele "gopkg.in/telebot.v3"
)

// Handler for utility buttons like menu, back, close, etc.
func UtilityHandler(b *tele.Bot, client *internal.Client, currentStep *string) {
	b.Handle(&types.BtnBack, func(c tele.Context) error {
		if *currentStep == "confirmOrder" {
			*currentStep = ""
		}
		return c.Send("Back to main menu", types.Menu)
	})
}

package handler

import (
	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func HandlePriceAlert(b internal.Bot, localizer *i18n.Localizer, currentStep *string, btnPriceAlert, btnCreatePriceAlert, btnViewPriceAlert, btnDeletePriceAlert, btnUpdatePriceAlert tele.Btn) {
	b.Handle(&btnPriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Price alert", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// CRUD operations for price alerts
	// Create price alert
	b.Handle(&btnCreatePriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Create price alert", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// View price alerts
	b.Handle(&btnViewPriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "View price alert", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Delete price alert
	b.Handle(&btnDeletePriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Delete price alert", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// Update price alert
	b.Handle(&btnUpdatePriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Update price alert", nil)
		if err != nil {
			return err
		}
		return nil
	})

}

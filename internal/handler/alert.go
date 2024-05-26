package handler

import (
	"errors"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

func HandlePriceAlert(b internal.Bot, localizer *i18n.Localizer, currentStep *string, btnPriceAlert, btnCreatePriceAlert, btnViewPriceAlert, btnDeletePriceAlert, btnUpdatePriceAlert tele.Btn) {
	b.Handle(&btnPriceAlert, func(c tele.Context) error {
		return c.Send("Price alert hihi", types.PriceAlertMenu(localizer))
	})

	// CRUD operations for price alerts
	// Create price alert
	b.Handle(&btnCreatePriceAlert, func(c tele.Context) error {
		*currentStep = "createPriceAlert"
		_, err := b.Send(c.Chat(), "Please input the alert in the format <<symbol>> <<condition>> <<value>>", tele.ForceReply)
		if err != nil {
			return err
		}
		return nil
	})

	// View price alerts
	b.Handle(&btnViewPriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "View price alert")
		if err != nil {
			return err
		}
		return nil
	})

	// Delete price alert
	b.Handle(&btnDeletePriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Delete price alert")
		if err != nil {
			return err
		}
		return nil
	})

	// Update price alert
	b.Handle(&btnUpdatePriceAlert, func(c tele.Context) error {
		_, err := b.Send(c.Chat(), "Update price alert")
		if err != nil {
			return err
		}
		return nil
	})

}

func HandleCreateAlertStep(b internal.Bot, c tele.Context, localizer *i18n.Localizer, currentStep *string) error {
	switch *currentStep {
	case "createPriceAlert":
		_, err := ParseCreateAlertInput(c.Text())
		if err != nil {
			return c.Send("Invalid input. Please try again", types.PriceAlertMenu(localizer))
		}
		_, err = b.Send(c.Chat(), "Alert created", types.PriceAlertMenu(localizer))
		if err != nil {
			return err
		}
		*currentStep = ""
		return nil
	}
	return nil
}

type Condition int

const (
	Greater Condition = iota
	Less
	PercentageGreater
	PercentageLess
)

type Alert struct {
	Value     string
	Condition Condition
	Symbol    string
}

// Let's say input is form of <<symbol>> <<condition>> <<value>>
func ParseCreateAlertInput(input string) (*Alert, error) {
	var alert *Alert
	parts := strings.Fields(input)

	if len(parts) != 3 {
		return nil, errors.New("invalid input")
	}

	symbol := parts[0]
	condition := parts[1]
	value := parts[2]

	if symbol == "" || condition == "" || value == "" {
		return nil, errors.New("invalid input")
	}

	supportedTokens := map[string]bool{
		"BTC":  true,
		"ETH":  true,
		"ATOM": true,
	}
	if !supportedTokens[strings.ToUpper(symbol)] {
		return nil, errors.New("unsupported token")
	}

	switch condition {
	case "up":
		if strings.HasSuffix(value, "%") {
			value = strings.TrimSuffix(value, "%")
			alert = &Alert{
				Value:     value,
				Condition: PercentageGreater,
				Symbol:    symbol,
			}
		} else {
			alert = &Alert{
				Value:     value,
				Condition: Greater,
				Symbol:    symbol,
			}
		}
	case "down":
		if strings.HasSuffix(value, "%") {
			value = strings.TrimSuffix(value, "%")
			alert = &Alert{
				Value:     value,
				Condition: PercentageLess,
				Symbol:    symbol,
			}
		} else {
			alert = &Alert{
				Value:     value,
				Condition: Less,
				Symbol:    symbol,
			}
		}
	default:
		return nil, errors.New("invalid input")
	}

	return alert, nil
}

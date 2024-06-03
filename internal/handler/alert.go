package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/google/uuid"
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
		endpoint := os.Getenv("ALERT_CRUD_ENDPOINT")
		alerts, err := ViewAlert(endpoint+"alerts/user/", c.Chat().Recipient())
		if err != nil {
			return c.Send("Error viewing alert"+err.Error(), types.PriceAlertMenu(localizer))
		}
		var message string
		priceMap := make(map[string]float32)
		changeMap := make(map[string]float32)
		if len(alerts.Alerts) == 0 {
			message = "No alerts found"
		} else {
			message += fmt.Sprintf("You have %d active price alerts:\n", len(alerts.Alerts))
			for _, alert := range alerts.Alerts {
				if priceMap[alert.Symbol] == 0 {
					priceMap[alert.Symbol] = 0.3202
				}
				if changeMap[alert.Symbol] == 0 {
					changeMap[alert.Symbol] = 0.3202
				}
				message += FormatAlert(alert, localizer, priceMap[alert.Symbol], changeMap[alert.Symbol])
			}
		}
		_, err = b.Send(c.Chat(), message, types.PriceAlertMenu(localizer), tele.ModeMarkdown)
		if err != nil {
			return err
		}
		return nil
	})

	// Delete price alert
	b.Handle(&btnDeletePriceAlert, func(c tele.Context) error {
		*currentStep = "deletePriceAlert"
		_, err := b.Send(c.Chat(), "Please input the alert ID to delete", tele.ForceReply)
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

func HandleAlertStep(b internal.Bot, c tele.Context, localizer *i18n.Localizer, currentStep *string) error {
	switch *currentStep {
	case "createPriceAlert":
		alert, err := ParseCreateAlertInput(c.Text())
		if err != nil {
			return c.Send("Invalid input. Please try again", types.PriceAlertMenu(localizer))
		}
		// Save alert via API
		endpoint := os.Getenv("ALERT_CRUD_ENDPOINT") + "alert"
		err = CreateAlert(endpoint, alert, c.Chat().Recipient())
		if err != nil {
			return c.Send("Error creating alert"+err.Error(), types.PriceAlertMenu(localizer))
		}
		_, err = b.Send(c.Chat(), "Alert created", types.PriceAlertMenu(localizer))
		if err != nil {
			return err
		}
		*currentStep = ""
		return nil
	case "deletePriceAlert":
		alertId := c.Text()
		if !ParseAlertId(alertId) {
			return c.Send("Invalid alert ID", types.PriceAlertMenu(localizer))
		}

		// delete alert via API
		endpoint := os.Getenv("ALERT_CRUD_ENDPOINT") + "alert/"
		err := DeleteAlert(endpoint, alertId)
		if err != nil {
			return c.Send("Error deleting alert"+err.Error(), types.PriceAlertMenu(localizer))
		}
		_, err = b.Send(c.Chat(), "Alert deleted", types.PriceAlertMenu(localizer))
		if err != nil {
			return err
		}
		*currentStep = ""

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

type CreateAlertRequest struct {
	UserID    string
	Value     float32
	Condition string
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

func ConvertAlertToCreateAlertRequest(alert *Alert, userID string) (*CreateAlertRequest, error) {
	var condition string
	switch alert.Condition {
	case Greater:
		condition = "PRICE_ABOVE"
	case Less:
		condition = "PRICE_BELOW"
	case PercentageGreater:
		condition = "PRICE_PERCENT_CHANGE_ABOVE"
	case PercentageLess:
		condition = "PRICE_PERCENT_CHANGE_BELOW"
	default:
		return nil, errors.New("invalid condition")
	}

	value, err := strconv.ParseFloat(alert.Value, 32)
	if err != nil {
		return nil, err
	}

	return &CreateAlertRequest{
		UserID:    userID,
		Value:     float32(value),
		Condition: condition,
		Symbol:    alert.Symbol,
	}, nil
}

// http://localhost:8080/alert?userId=user1&symbol=BTC&value=66900.000000&condition=PRICE_BELOW
func CreateAlert(endpoint string, alert *Alert, userID string) error {
	// Send alert to server
	httpClient := &http.Client{}
	createAlertReq, err := ConvertAlertToCreateAlertRequest(alert, userID)
	if err != nil {
		return err
	}
	endpoint = fmt.Sprintf("%s?userId=%s&symbol=%s&value=%f&condition=%s", endpoint, createAlertReq.UserID, createAlertReq.Symbol, createAlertReq.Value, createAlertReq.Condition)
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("error creating alert")
	}

	return nil
}

func ViewAlert(endpoint string, userID string) (*types.Alerts, error) {
	// Send alert to server
	httpClient := &http.Client{}
	endpoint = fmt.Sprintf("%s%s", endpoint, userID)
	fmt.Println("view alert endpoint: ", endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error viewing alert: %d", resp.StatusCode)
	}
	// parse response
	alerts := &types.Alerts{}
	err = json.NewDecoder(resp.Body).Decode(alerts)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// You have 1 active price alert:

// Alert id: UeIUXHnVI3vOWdBP2sTlA
// 🔗 Bitcoin BTC/USDT
// ⏰ Alert me when price is over $2.00
// Current price in USD: $0.3202

func FormatAlert(alert *types.Alert, localizer *i18n.Localizer, currentPriceInUSD float32, change float32) string {
	condition := ""
	switch alert.Condition {
	case types.Condition_PRICE_ABOVE:
		condition = "over"
	case types.Condition_PRICE_BELOW:
		condition = "below"
	case types.Condition_PRICE_PERCENT_CHANGE_ABOVE:
		condition = "up"
	case types.Condition_PRICE_PERCENT_CHANGE_BELOW:
		condition = "down"
	}
	fmt.Println("condition: ", condition)
	name := ""
	switch strings.ToUpper(alert.Symbol) {
	case "BTC":
		name = "Bitcoin"
	case "ETH":
		name = "Ethereum"
	case "ATOM":
		name = "Cosmos"
	}
	if alert.Condition == types.Condition_PRICE_PERCENT_CHANGE_ABOVE || alert.Condition == types.Condition_PRICE_PERCENT_CHANGE_BELOW {
		fmt.Println("huhu")
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "AlertMessagePercentage",
				Other: "*🚨 Alert ID:* {{.AlertID}}\n" +
					"*📈 {{.Name}} {{.Symbol}}/USDT*\n" +
					"⏰ Alert me when price is *{{.Condition}} {{.Value}}% last 24h*\n" +
					"💰 Current price in USD: ${{.CurrentPrice}}",
			},
			TemplateData: map[string]interface{}{
				"AlertID":      alert.Id,
				"Name":         name,
				"Symbol":       strings.ToUpper(alert.Symbol),
				"Condition":    condition,
				"Value":        fmt.Sprintf("%.2f", alert.Value),
				"CurrentPrice": fmt.Sprintf("%.2f", currentPriceInUSD),
			},
		})
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "AlertMessage",
			Other: "*🚨 Alert ID:* {{.AlertID}}\n" +
				"*📈 {{.Name}} {{.Symbol}}/USDT*\n" +
				"⏰ Alert me when price is *{{.Condition}} ${{.Value}}*\n" +
				"💰 Current price in USD: ${{.CurrentPrice}}",
		},
		TemplateData: map[string]interface{}{
			"AlertID":      alert.Id,
			"Name":         name,
			"Symbol":       strings.ToUpper(alert.Symbol),
			"Condition":    condition,
			"Value":        fmt.Sprintf("%.2f", alert.Value),
			"CurrentPrice": fmt.Sprintf("%.2f", currentPriceInUSD),
		},
	})
}

// check if valid uuid
func ParseAlertId(input string) bool {
	_, err := uuid.Parse(input)
	if err != nil {
		return false
	}
	return true
}

func DeleteAlert(endpoint, alertId string) error {
	httpClient := &http.Client{}
	endpoint = fmt.Sprintf("%s?alertId=%s", endpoint, alertId)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("error deleting alert")
	}
	return nil
}

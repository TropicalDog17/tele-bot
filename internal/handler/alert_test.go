package handler_test

import (
	"testing"

	"github.com/TropicalDog17/tele-bot/internal/handler"
)

func TestParseInputAlert(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected *handler.Alert
	}{
		{
			"valid absolute value alert",
			"BTC up 5000",
			&handler.Alert{
				Symbol:    "BTC",
				Condition: handler.Up,
				Value:     "5000",
			},
		},
		{
			"valid percentage value alert up",
			"ETH down 5%",
			&handler.Alert{
				Symbol:    "ETH",
				Condition: handler.PercentageLess,
				Value:     "5",
			},
		},
		{
			"valid percentage value alert down",
			"ETH up 100%",
			&handler.Alert{
				Symbol:    "ETH",
				Condition: handler.PercentageGreater,
				Value:     "100",
			},
		},
		{
			"valid absolute value alert above",
			"ETH over 4000",
			&handler.Alert{
				Symbol:    "ETH",
				Condition: handler.Above,
				Value:     "4000",
			},
		},
		{
			"valid absolute value alert below",
			"BTC under 50000",
			&handler.Alert{
				Symbol:    "BTC",
				Condition: handler.Below,
				Value:     "50000",
			},
		},
		{
			"invalid absolute percentage alert below",
			"BTC under 50%",
			nil,
		},
		{
			"invalid input",
			"btc up",
			nil,
		},
		{
			"invalid value",
			"ATOM ldjsd 10%",
			nil,
		},
		{
			"unsupported token",
			"okdfhkfh up 10",
			nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			alert, err := handler.ParseCreateAlertInput(tc.input)
			if err != nil {
				if tc.expected != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if alert.Symbol != tc.expected.Symbol {
				t.Errorf("expected symbol %s, got %s", tc.expected.Symbol, alert.Symbol)
			}
			if alert.Condition != tc.expected.Condition {
				t.Errorf("expected condition %v, got %v", tc.expected.Condition, alert.Condition)
			}
			if alert.Value != tc.expected.Value {
				t.Errorf("expected value %s, got %s", tc.expected.Value, alert.Value)
			}
		})
	}
}

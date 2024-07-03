package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

type MarketData struct {
	Change   float64 `json:"change"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	MarketID string  `json:"marketId"`
	Open     float64 `json:"open"`
	Price    float64 `json:"price"`
	Volume   float64 `json:"volume"`
}

type DisplayData struct {
	Ticker string `json:"ticker"`
	Change string `json:"change"`
	Price  string `json:"price"`
	Volume string `json:"volume"`
}

var LinkToHelix = "\n\n[View on Helix](https://helixapp.com/markets/?type=spot)"

func HandleViewMarket(b internal.Bot, localizer *i18n.Localizer, btnViewMarket, btnBiggestGainer24h, btnBiggestLoser24h tele.Btn, btnBiggestVolume24h tele.Btn) {
	b.Handle(&btnViewMarket, func(c tele.Context) error {
		return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "market_quick_look",
			DefaultMessage: &i18n.Message{
				ID:    "market_quick_look",
				Other: "Here you can have a quick look at the market the last 24h",
			},
		}), types.ViewMarketMenu(localizer))
	})

	b.Handle(&btnBiggestGainer24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "error_fetching_data",
				DefaultMessage: &i18n.Message{
					ID:    "error_fetching_data",
					Other: "Error fetching data: {{.ErrorMessage}}",
				},
				TemplateData: map[string]interface{}{
					"ErrorMessage": err.Error(),
				},
			}), types.Menu)
		}
		gainers := GetTopNBiggestGainer(data, 5)
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "biggest_gainers_24h",
			DefaultMessage: &i18n.Message{
				ID:    "biggest_gainers_24h",
				Other: "Here are the biggest gainers in the last 24h 📈📈📈 \n ",
			},
		})
		return c.Send(text+DisplayDataToString(gainers, localizer)+LinkToHelix, types.Menu, types.ViewMarketMenu(localizer), tele.ModeMarkdown)
	})

	b.Handle(&btnBiggestLoser24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "error_fetching_data",
				DefaultMessage: &i18n.Message{
					ID:    "error_fetching_data",
					Other: "Error fetching data: {{.ErrorMessage}}",
				},
				TemplateData: map[string]interface{}{
					"ErrorMessage": err.Error(),
				},
			}), types.Menu)
		}
		losers := GetTopNBiggestLoser(data, 5)
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "biggest_losers_24h",
			DefaultMessage: &i18n.Message{
				ID:    "biggest_losers_24h",
				Other: "Here are the biggest losers in the last 24h 📉📉📉 \n ",
			},
		})
		return c.Send(text+DisplayDataToString(losers, localizer)+LinkToHelix, types.Menu, types.ViewMarketMenu(localizer), tele.ModeMarkdown)
	})

	b.Handle(&btnBiggestVolume24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "error_fetching_data",
				DefaultMessage: &i18n.Message{
					ID:    "error_fetching_data",
					Other: "Error fetching data: {{.ErrorMessage}}",
				},
				TemplateData: map[string]interface{}{
					"ErrorMessage": err.Error(),
				},
			}), types.Menu)
		}
		volume := GetTopNBiggestVolume(data, 5)
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "biggest_volume_24h",
			DefaultMessage: &i18n.Message{
				ID:    "biggest_volume_24h",
				Other: "Here are the biggest volume in the last 24h 📊📊📊 \n ",
			},
		})
		return c.Send(text+DisplayDataToString(volume, localizer)+LinkToHelix, types.Menu, types.ViewMarketMenu(localizer), tele.ModeMarkdown)
	})
}

func FetchMarketsDataLast24h() ([]MarketData, error) {
	url := "https://sentry.exchange.grpc-web.injective.network/api/chronos/v1/spot/market_summary_all?resolution=24h"
	// Fetch data from the exchange
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	// Unmarshal the response body

	var data []MarketData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetTopNBiggestGainer(data []MarketData, n int) []MarketData {
	// Sort the data by the change in descending order
	sort.Slice(data, func(i, j int) bool {
		return data[i].Change > data[j].Change
	})

	// Return the top N biggest gainers
	if n > len(data) {
		n = len(data)
	}
	return data[:n]
}

func GetTopNBiggestLoser(data []MarketData, n int) []MarketData {
	// Sort the data by the change in ascending order
	sort.Slice(data, func(i, j int) bool {
		return data[i].Change < data[j].Change
	})

	// Return the top N biggest losers
	if n > len(data) {
		n = len(data)
	}
	return data[:n]
}

func GetTopNBiggestVolume(data []MarketData, n int) []MarketData {
	// Sort the data by the volume in descending order
	sort.Slice(data, func(i, j int) bool {
		return data[i].Volume > data[j].Volume
	})

	// Return the top N biggest volume
	if n > len(data) {
		n = len(data)
	}
	return data[:n]
}

func (m MarketData) Display() DisplayData {
	ticker, _ := MarketIDToTicker(m.MarketID)
	return DisplayData{
		Ticker: ticker,
		Change: formatFloat(m.Change),
		Price:  formatFloat(m.Price),
		Volume: formatFloat(m.Volume),
	}
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

type Market struct {
	Ticker string `json:"ticker"`
}

type Response struct {
	Market Market `json:"market"`
}

func MarketIDToTicker(marketID string) (string, error) {
	url := "https://sentry.exchange.grpc-web.injective.network/api/exchange/spot/v1/markets/" + marketID
	// Fetch data from the exchange
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return resp.Market.Ticker, nil
	// Unmarshal the response body

}
func (d DisplayData) String(localizer *i18n.Localizer) string {
	icon := "🟢⬆️"
	if d.Change[0] == '-' {
		icon = "🔴⬇️"
	}

	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "display_data_string",
		DefaultMessage: &i18n.Message{
			ID: "display_data_string",
			Other: `{{.Ticker}}

Change: {{.Change}}% {{.Icon}}

Price: ${{.Price}} 📊 Volume: ${{.Volume}}`,
		},
		TemplateData: map[string]interface{}{
			"Ticker": d.Ticker,
			"Change": d.Change,
			"Icon":   icon,
			"Price":  d.Price,
			"Volume": d.Volume,
		},
	})
}
func DisplayDataToString(data []MarketData, localizer *i18n.Localizer) string {
	var result string
	for i, d := range data {
		var prefix string
		switch i {
		case 0:
			prefix = "🥇 "
		case 1:
			prefix = "🥈 "
		case 2:
			prefix = "🥉 "
		default:
			prefix = fmt.Sprintf("%d. ", i+1)
		}
		result += prefix + d.Display().String(localizer) + "\n\n"
	}
	return result

}
func MockFetchData24h() ([]MarketData, error) {
	// read from tests/markets.json
	buffer, err := os.Open("internal/handler/tests/markets.json")
	if err != nil {
		return nil, err
	}
	defer buffer.Close()
	jsonData, err := io.ReadAll(buffer)
	if err != nil {
		return nil, err
	}
	var data []MarketData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

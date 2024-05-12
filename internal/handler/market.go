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

func HandleViewMarket(b internal.Bot) {
	b.Handle(&types.BtnViewMarket, func(c tele.Context) error {
		return c.Send("Here you can have a quick look at the market the last 24h", types.MenuViewMarket)
	})
	b.Handle(&types.BtnBiggestGainer24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send("Error fetching data"+err.Error(), types.Menu)
		}
		gainers := GetTopNBiggestGainer(data, 5)
		text := "Here are the biggest gainers in the last 24h 📈📈📈 \n "
		return c.Send(text+DisplayDataToString(gainers)+LinkToHelix, types.Menu, types.MenuViewMarket, tele.ModeMarkdown)
	})
	b.Handle(&types.BtnBiggestLoser24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send("Error fetching data"+err.Error(), types.Menu)
		}
		losers := GetTopNBiggestLoser(data, 5)
		text := "Here are the biggest losers in the last 24h 📉📉📉 \n "
		return c.Send(text+DisplayDataToString(losers)+LinkToHelix, types.Menu, types.MenuViewMarket, tele.ModeMarkdown)
	})
	b.Handle(&types.BtnBiggestVolume24h, func(c tele.Context) error {
		data, err := FetchMarketsDataLast24h()
		if err != nil {
			return c.Send("Error fetching data"+err.Error(), types.Menu)
		}
		volume := GetTopNBiggestVolume(data, 5)
		text := "Here are the biggest volume in the last 24h 📊📊📊 \n "
		return c.Send(text+DisplayDataToString(volume)+LinkToHelix, types.Menu, types.MenuViewMarket, tele.ModeMarkdown)
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
func (d DisplayData) String() string {
	icon := "🟢⬆️"
	if d.Change[0] == '-' {
		icon = "🔴⬇️"
	}
	return fmt.Sprintf("%s \n\n Change: %s%% %s \n\n Price: $%s 📊 Volume: $%s",
		d.Ticker, d.Change, icon, d.Price, d.Volume)
}
func DisplayDataToString(data []MarketData) string {
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
		result += prefix + d.Display().String() + "\n\n"
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

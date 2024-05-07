package main

import (
	"log"
	"os"

	"github.com/TropicalDog17/tele-bot/config"
	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/handler"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var (
	selectedToken    string
	selectedAmount   string
	recipientAddress string
	currentStep      string
	globalMenu       tele.StoredMessage
	limitOrderMenu   tele.StoredMessage
	createOrderMenu  tele.StoredMessage
)

var (
	globalLimitOrder = types.NewLimitOrderInfo()
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	menu, menuSendToken, menuLimitOrder, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders := types.InitializeUI()[0], types.InitializeUI()[1], types.InitializeUI()[2], types.InitializeUI()[3], types.InitializeUI()[4], types.InitializeUI()[5]
	pref := config.NewBotPref(os.Getenv("TELEGRAM_TOKEN"))
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	client := internal.NewClient()

	// On start command
	handler.HandleOnboard(b, client, &currentStep)

	// On reply button pressed (message)
	b.Handle(&types.BtnViewBalances, func(c tele.Context) error {
		// Unimplemented
		return c.Send("View Balances", menu)
	})

	b.Handle("/menu", func(c tele.Context) error {
		return c.Send("Menu", menu)
	})
	handler.HandleAccountDetails(b, client)
	handler.HandleAddressQr(b, client)

	// Handle the "Limit Order" flow
	handler.HandleLimitOrder(b, client, &limitOrderMenu, &createOrderMenu, globalLimitOrder, &currentStep, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders)
	handler.UtilityHandler(b, client, &currentStep)

	// Handle the transfer token flow
	handler.HandlerTransferToken(b, client, menuSendToken, &types.BtnInlineAtom, &types.BtnInlineInj, &types.BtnTenDollar, &types.BtnFiftyDollar, &types.BtnHundredDollar, &types.BtnTwoHundredDollar, &types.BtnFiveHundredDollar, &types.BtnCustomAmount, &types.BtnRecipientSection, &types.BtnCustomToken, &selectedToken, &selectedAmount, &currentStep, &recipientAddress, &globalMenu)

	handler.HandleStep(b, client, utils.Utils{}, &currentStep, menuSendToken, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, &selectedAmount, &selectedToken, &recipientAddress, &globalMenu, &createOrderMenu)

	b.Start()
}

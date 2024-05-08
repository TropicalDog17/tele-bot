package main

import (
	"log"
	"os"

	"github.com/TropicalDog17/tele-bot/config"
	"github.com/TropicalDog17/tele-bot/internal"
	clienttypes "github.com/TropicalDog17/tele-bot/internal/client"
	"github.com/TropicalDog17/tele-bot/internal/handler"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	memguard "github.com/awnumar/memguard"
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
var clients = make(map[string]internal.BotClient)

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
	b.Use(clientMiddleware)
	// On start command
	handler.HandleOnboard(b, clients, &currentStep)

	// On reply button pressed (message)
	b.Handle(&types.BtnViewBalances, func(c tele.Context) error {
		// Unimplemented
		return c.Send("View Balances", menu)
	})

	b.Handle("/menu", func(c tele.Context) error {

		return c.Send("Menu", menu)
	})
	handler.HandleAccountDetails(b, clients)
	handler.HandleAddressQr(b, clients)

	// Handle the "Limit Order" flow
	handler.HandleLimitOrder(b, clients, &limitOrderMenu, &createOrderMenu, globalLimitOrder, &currentStep, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders)
	handler.UtilityHandler(b, &currentStep)

	// Handle the transfer token flow
	handler.HandlerTransferToken(b, clients, menuSendToken, &types.BtnInlineAtom, &types.BtnInlineInj, &types.BtnTenDollar, &types.BtnFiftyDollar, &types.BtnHundredDollar, &types.BtnTwoHundredDollar, &types.BtnFiveHundredDollar, &types.BtnCustomAmount, &types.BtnRecipientSection, &types.BtnCustomToken, &selectedToken, &selectedAmount, &currentStep, &recipientAddress, &globalMenu)
	handler.HandleViewMarket(b)

	handler.HandleStep(b, clients, utils.Utils{}, &currentStep, menuSendToken, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, &selectedAmount, &selectedToken, &recipientAddress, &globalMenu, &createOrderMenu)
	b.Start()
}

var notWaitingForPassword = make(map[string]bool)

func clientMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		username := c.Sender().Username
		_, ok := clients[username]
		if !ok {
			// Client doesn't exist, check if waiting for password
			if notWaitingForPassword[username] {
				// User has provided the password
				password := c.Text()

				// Create a new LockedBuffer with the password
				pwdBuffer := memguard.NewBufferFromBytes([]byte(password))

				// Create a new Client with the provided password
				client, err := clienttypes.NewClient(c.Bot(), username, pwdBuffer, &currentStep)
				if err != nil {
					// Password is invalid
					_, _ = c.Bot().Send(c.Recipient(), "Invalid password. Please re-enter your password")
					return nil
				}
				clients[username] = client

				// Reset the waiting flag
				notWaitingForPassword[username] = true

				// Clean up the password buffer
				pwdBuffer.Destroy()
				_ = c.Delete()
				// Proceed with the next handler
				return c.Send("Password accepted. You can perform your action again", types.Menu)
			} else {
				// Client doesn't exist and not waiting for password
				// Send password request message
				_, _ = c.Bot().Send(c.Recipient(), "Session expired! Please enter your password")

				// Set the waiting flag for the user
				notWaitingForPassword[username] = true

				// Initialize the current step for the user
				currentStep = "askPassword"

				// Stop further processing
				return nil
			}
		}

		// Client exists, proceed with the next handler
		return next(c)
	}
}

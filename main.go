package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/TropicalDog17/tele-bot/config"
	"github.com/TropicalDog17/tele-bot/internal"
	clienttypes "github.com/TropicalDog17/tele-bot/internal/client"
	"github.com/TropicalDog17/tele-bot/internal/database"
	"github.com/TropicalDog17/tele-bot/internal/handler"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	memguard "github.com/awnumar/memguard"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var redisInstance internal.RedisClient
var notWaitingForPassword = make(map[string]bool)

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
	redisInstance = database.NewRedisInstance()
	menu, menuSendToken, menuLimitOrder, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders := types.InitializeUI()[0], types.InitializeUI()[1], types.InitializeUI()[2], types.InitializeUI()[3], types.InitializeUI()[4], types.InitializeUI()[5]
	pref := config.NewBotPref(os.Getenv("TELEGRAM_TOKEN"))
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	authRoute := b.Group()
	authRoute.Use(clientMiddleware)
	// On start command
	handler.HandleOnboard(b, clients, &currentStep)

	// On reply button pressed (message)
	authRoute.Handle(&types.BtnViewBalances, func(c tele.Context) error {
		// Unimplemented
		return c.Send("View Balances", menu)
	})

	authRoute.Handle("/menu", func(c tele.Context) error {
		return c.Send("Menu", menu)
	})
	handler.HandleAccountDetails(b, authRoute, clients)
	handler.HandleAddressQr(b, authRoute, clients)

	// Handle the "Limit Order" flow
	handler.HandleLimitOrder(b, authRoute, clients, &limitOrderMenu, &createOrderMenu, globalLimitOrder, &currentStep, menu, menuCreateLimitOrder, menuConfirmOrder, menuActiveOrders)
	handler.UtilityHandler(b, authRoute, &currentStep)

	// Handle the transfer token flow
	handler.HandlerTransferToken(b, authRoute, clients, menuSendToken, &types.BtnInlineAtom, &types.BtnInlineInj, &types.BtnTenDollar, &types.BtnFiftyDollar, &types.BtnHundredDollar, &types.BtnTwoHundredDollar, &types.BtnFiveHundredDollar, &types.BtnCustomAmount, &types.BtnRecipientSection, &types.BtnCustomToken, &selectedToken, &selectedAmount, &currentStep, &recipientAddress, &globalMenu)
	handler.HandleViewMarket(b)

	handler.HandleStep(b, authRoute, clients, utils.Utils{}, &currentStep, menuSendToken, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, &selectedAmount, &selectedToken, &recipientAddress, &globalMenu, &createOrderMenu)
	b.Start()
}

var authSteps = []string{
	"customAmount", "recipientAddress", "limitAmount", "limitPrice", "limitToken", "cancelOrder", "confirmOrder", types.BtnSendToken.Text, types.BtnLimitOrder.Text, types.BtnViewBalances.Text, types.BtnShowAccount.Text,
	types.BtnActiveOrders.Text, types.BtnCancelOrder.Text, types.BtnBack.Text, types.BtnMenu.Text, types.BtnInlineAtom.Text, types.BtnInlineInj.Text, types.BtnTenDollar.Text, types.BtnFiftyDollar.Text, types.BtnHundredDollar.Text, types.BtnTwoHundredDollar.Text, types.BtnFiveHundredDollar.Text, types.BtnCustomAmount.Text, types.BtnRecipientSection.Text, types.BtnCustomToken.Text,
}

func clientMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {

		username := c.Sender().Username
		_, ok := clients[username]
		if !ok {
			// Client doesn't exist, check if waiting for password
			if notWaitingForPassword[username] && currentStep == "askPassword" {
				// User has provided the password
				password := c.Text()

				// Create a new LockedBuffer with the password
				pwdBuffer := memguard.NewBufferFromBytes([]byte(password))

				// Create a new Client with the provided password
				client, err := clienttypes.NewClient(c.Bot(), username, pwdBuffer, redisInstance, &currentStep)
				if err != nil {
					// Password is invalid
					_, _ = c.Bot().Send(c.Recipient(), "Invalid password. Please re-enter your password")
					return nil
				}

				fmt.Println("Client created successfully")
				clients[username] = client

				// Reset the waiting flag
				notWaitingForPassword[username] = true
				redisInstance.HSet(context.Background(), username, "client", true)
				// Clean up the password buffer
				pwdBuffer.Destroy()

				_ = c.Delete()
				// Retrieve the Session expired! Please enter your password message
				chatID := redisInstance.HGet(context.Background(), username, "currentExpiredChatId").Val()
				msgID := redisInstance.HGet(context.Background(), username, "currentExpiredMsgId").Val()

				chatIDInt, _ := strconv.ParseInt(chatID, 10, 64)
				_ = c.Bot().Delete(tele.StoredMessage{
					ChatID:    chatIDInt,
					MessageID: msgID,
				})

				// Proceed with the next handler
				return c.Send("Password accepted. You can perform your action again", types.Menu)
			} else {

				if currentStep != "" && !slices.Contains(authSteps, currentStep) {
					return next(c)
				}
				fmt.Println("Current step: ", c.Message().Text)
				if !slices.Contains(authSteps, c.Message().Text) {
					return next(c)
				}
				// Client doesn't exist and not waiting for password
				// Check if any credentials exist for the user
				haveCreds := redisInstance.HExists(context.Background(), username, "salt").Val()
				fmt.Println(haveCreds)
				if !haveCreds {
					isFirstTime := !redisInstance.HExists(context.Background(), username, "client").Val()
					fmt.Println(c.Message().Text)
					if isFirstTime && c.Message().Text == "/start" {
						return next(c)
					} else if isFirstTime {
						_, _ = c.Bot().Send(c.Recipient(), "Please type /start to start the bot")
						return nil
					}
				}

				msg, _ := c.Bot().Send(c.Recipient(), "Session expired! Please enter your password")
				redisInstance.HSet(context.Background(), username, "currentExpiredChatId", msg.Chat.ID)
				redisInstance.HSet(context.Background(), username, "currentExpiredMsgId", msg.ID)
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

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

var transferInfo = &types.TransferInfo{}
var redisInstance internal.RedisClient
var notWaitingForPassword = make(map[string]bool)
var bundle = i18n.NewBundle(language.English)
var localizer = i18n.NewLocalizer(bundle, "en")
var authRoute *tele.Group
var (
	currentStep     string
	globalMenu      tele.StoredMessage
	limitOrderMenu  tele.StoredMessage
	createOrderMenu tele.StoredMessage
)

var clients = make(map[string]internal.BotClient)
var globalLimitOrder = types.NewLimitOrderInfo()

func main() {
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("active.vi.toml")
	bundle.MustLoadMessageFile("active.en.toml")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redisInstance = database.NewRedisInstance()
	pref := config.NewBotPref(os.Getenv("TELEGRAM_TOKEN"))
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	authRoute = b.Group()
	authRoute.Use(clientMiddleware)
	// On start command
	b.Use(languageMiddleware)
	handler.HandleOnboard(b, clients, &currentStep)

	authRoute.Handle("/menu", func(c tele.Context) error {
		return c.Send("Menu", types.MainMenu(localizer))
	})

	SetHandlerForBot(b, localizer, authRoute, globalLimitOrder, transferInfo)
	b.Start()
}

func SetHandlerForBot(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, globalLimitOrder *types.LimitOrderInfo, transferInfo *types.TransferInfo) {
	handler.HandleAccountDetails(b, localizer, authRoute, clients, types.BtnShowAccount(localizer))
	handler.HandleAddressQr(b, authRoute, clients)

	// Handle the "Limit Order" flow
	handler.HandleLimitOrder(b, localizer, authRoute, clients, &limitOrderMenu, &createOrderMenu, globalLimitOrder, types.BtnLimitOrder(localizer), types.BtnBuyLimitOrder(localizer), types.BtnSellLimitOrder(localizer), types.BtnAmount(localizer), types.BtnPrice(localizer), types.BtnToken(localizer), types.BtnPayWith(localizer), types.BtnConfirmOrder(localizer), types.BtnConfirmLimitOrder(localizer), types.BtnActiveOrders(localizer), types.BtnCancelOrder(localizer), &currentStep, types.MainMenu(localizer), types.CreateLimitOrderMenu(localizer), types.ConfirmOrderMenu(localizer), types.ActiveOrdersMenu(localizer))
	// handler.UtilityHandler(b, localizer, authRoute, &currentStep)

	// Handle the transfer token flow
	handler.HandlerTransferToken(b, localizer, authRoute, clients, types.SendTokenMenu(localizer, transferInfo), types.BtnSend(localizer), types.BtnSendToken(localizer), types.BtnInlineAtom(localizer), types.BtnInlineInj(localizer), types.BtnTenDollar(localizer), types.BtnFiftyDollar(localizer), types.BtnHundredDollar(localizer), types.BtnTwoHundredDollar(localizer), types.BtnFiveHundredDollar(localizer), types.BtnCustomAmount(localizer), types.BtnRecipientSection(localizer, transferInfo), types.BtnCustomToken(localizer), transferInfo, &currentStep, &globalMenu)
	handler.HandleViewMarket(b, localizer, types.BtnViewMarket(localizer), types.BtnBiggestGainer24h(localizer), types.BtnBiggestLoser24h(localizer), types.BtnBiggestVolume24h(localizer))
	handler.HandleSettings(b, localizer, authRoute, clients, types.ViewSettingsMenu(localizer), types.BtnSettings(localizer), types.BtnChangeLanguage(localizer), &currentStep)
	handler.HandleStep(b, localizer, authRoute, clients, utils.Utils{}, &currentStep, types.SendTokenMenu(localizer, transferInfo), types.LimitOrderMenu(localizer), types.CreateLimitOrderMenu(localizer), globalLimitOrder, transferInfo, &globalMenu, &createOrderMenu)
	handler.HandlePriceAlert(b, localizer, &currentStep, types.BtnPriceAlert(localizer), types.BtnCreatePriceAlert(localizer), types.BtnViewPriceAlert(localizer), types.BtnDeletePriceAlert(localizer), types.BtnUpdatePriceAlert(localizer))
}

var authSteps = []string{
	"customAmount", "recipientAddress", "limitAmount", "limitPrice", "limitToken", "payWithToken", "cancelOrder", "confirmOrder", types.BtnSendToken(localizer).Text, types.BtnLimitOrder(localizer).Text, types.BtnShowAccount(localizer).Text,
	types.BtnActiveOrders(localizer).Text, types.BtnCancelOrder(localizer).Text, types.BtnBack(localizer).Text, types.BtnMenu(localizer).Text, types.BtnInlineAtom(localizer).Text, types.BtnInlineInj(localizer).Text, types.BtnTenDollar(localizer).Text, types.BtnFiftyDollar(localizer).Text, types.BtnHundredDollar(localizer).Text, types.BtnTwoHundredDollar(localizer).Text, types.BtnFiveHundredDollar(localizer).Text, types.BtnCustomAmount(localizer).Text, types.BtnRecipientSection(localizer, transferInfo).Text, types.BtnCustomToken(localizer).Text, types.BtnSettings(localizer).Text,
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
					fmt.Println("Error creating client: ", err)
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
				// "Password accepted. You can perform your action again"
				return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "passwordAccepted",
						Other: "Password accepted. You can perform your action again",
					},
				}), types.MainMenu(localizer))
			} else {

				if currentStep != "" && !slices.Contains(authSteps, currentStep) {
					return next(c)
				}
				fmt.Println("Current step: ", currentStep)
				fmt.Println("Auth steps: ", c.Message().Text)
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

				msg, _ := c.Bot().Send(c.Recipient(), localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "sessionExpired",
						Other: "Session expired! Please enter your password",
					},
				}))
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

func languageMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	authSteps = []string{
		"customAmount", "recipientAddress", "limitAmount", "limitPrice", "limitToken", "payWithToken", "cancelOrder", "confirmOrder", types.BtnSendToken(localizer).Text, types.BtnLimitOrder(localizer).Text, types.BtnShowAccount(localizer).Text,
		types.BtnActiveOrders(localizer).Text, types.BtnCancelOrder(localizer).Text, types.BtnBack(localizer).Text, types.BtnMenu(localizer).Text, types.BtnInlineAtom(localizer).Text, types.BtnInlineInj(localizer).Text, types.BtnTenDollar(localizer).Text, types.BtnFiftyDollar(localizer).Text, types.BtnHundredDollar(localizer).Text, types.BtnTwoHundredDollar(localizer).Text, types.BtnFiveHundredDollar(localizer).Text, types.BtnCustomAmount(localizer).Text, types.BtnRecipientSection(localizer, transferInfo).Text, types.BtnCustomToken(localizer).Text, types.BtnSettings(localizer).Text,
	}
	return func(c tele.Context) error {
		redisInstance = database.NewRedisInstance()
		username := c.Sender().Username

		// Check if the user has set a language
		language := redisInstance.HGet(context.Background(), username, "language").Val()
		if language == "" {
			// Language is not set
			localizer = i18n.NewLocalizer(bundle, "en-US")
			fmt.Println("Language not set")
		}
		if language == "en" {
			localizer = i18n.NewLocalizer(bundle, "en-US")
			fmt.Println("set eng")
		}
		if language == "vi" {
			localizer = i18n.NewLocalizer(bundle, "vi-VN")
			fmt.Println("set vi")
		}
		SetHandlerForBot(c.Bot(), localizer, authRoute, globalLimitOrder, transferInfo)
		// Language is set, proceed with the next handler
		return next(c)
	}
}

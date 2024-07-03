package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

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
var globalLimitOrder *types.LimitOrderInfo

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
	globalLimitOrder = types.NewLimitOrderInfo()
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
		ctx := context.Background()

		if _, ok := clients[username]; !ok {
			if notWaitingForPassword[username] && currentStep == "askPassword" {
				password := c.Text()
				pwdBuffer := memguard.NewBufferFromBytes([]byte(password))
				defer pwdBuffer.Destroy()

				client, err := clienttypes.NewClient(c.Bot(), username, pwdBuffer, redisInstance, &currentStep)
				if err != nil {
					return c.Send("Invalid password. Please re-enter your password")
				}

				clients[username] = client
				notWaitingForPassword[username] = true

				pipe := redisInstance.Pipeline()
				pipe.HSet(ctx, username, "client", "true")
				chatIDCmd := pipe.HGet(ctx, username, "currentExpiredChatId")
				msgIDCmd := pipe.HGet(ctx, username, "currentExpiredMsgId")
				_, err = pipe.Exec(ctx)
				if err != nil {
					return fmt.Errorf("Redis pipeline error: %w", err)
				}

				_ = c.Delete()
				chatID, _ := strconv.Atoi(chatIDCmd.Val())
				_ = c.Bot().Delete(tele.StoredMessage{
					ChatID:    int64(chatID),
					MessageID: msgIDCmd.Val(),
				})

				return c.Send(localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "passwordAccepted",
						Other: "Password accepted. You can perform your action again",
					},
				}), types.MainMenu(localizer))
			}

			if currentStep != "" && !slices.Contains(authSteps, currentStep) {
				return next(c)
			}

			messageText := c.Message().Text
			if !slices.Contains(authSteps, messageText) {
				return next(c)
			}

			pipe := redisInstance.Pipeline()
			haveCredsCmd := pipe.HExists(ctx, username, "salt")
			isFirstTimeCmd := pipe.HExists(ctx, username, "client")
			_, err := pipe.Exec(ctx)
			if err != nil {
				return fmt.Errorf("Redis pipeline error: %w", err)
			}

			if !haveCredsCmd.Val() {
				if !isFirstTimeCmd.Val() && strings.EqualFold(messageText, "/start") {
					return next(c)
				} else if !isFirstTimeCmd.Val() {
					return c.Send("Please type /start to start the bot")
				}
			}

			msg, _ := c.Bot().Send(c.Recipient(), localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "sessionExpired",
					Other: "Session expired! Please enter your password",
				},
			}))

			pipe = redisInstance.Pipeline()
			pipe.HSet(ctx, username, "currentExpiredChatId", msg.Chat.ID)
			pipe.HSet(ctx, username, "currentExpiredMsgId", msg.ID)
			_, err = pipe.Exec(ctx)
			if err != nil {
				return fmt.Errorf("Redis pipeline error: %w", err)
			}

			notWaitingForPassword[username] = true
			currentStep = "askPassword"
			return nil
		}

		return next(c)
	}
}

func languageMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	authSteps = []string{
		"customAmount", "recipientAddress", "limitAmount", "limitPrice", "limitToken", "payWithToken", "cancelOrder", "confirmOrder", types.BtnSendToken(localizer).Text, types.BtnLimitOrder(localizer).Text, types.BtnShowAccount(localizer).Text,
		types.BtnActiveOrders(localizer).Text, types.BtnCancelOrder(localizer).Text, types.BtnBack(localizer).Text, types.BtnMenu(localizer).Text, types.BtnInlineAtom(localizer).Text, types.BtnInlineInj(localizer).Text, types.BtnTenDollar(localizer).Text, types.BtnFiftyDollar(localizer).Text, types.BtnHundredDollar(localizer).Text, types.BtnTwoHundredDollar(localizer).Text, types.BtnFiveHundredDollar(localizer).Text, types.BtnCustomAmount(localizer).Text, types.BtnRecipientSection(localizer, transferInfo).Text, types.BtnCustomToken(localizer).Text, types.BtnSettings(localizer).Text,
	}
	var currentLanguage string

	return func(c tele.Context) error {
		redisInstance = database.NewRedisInstance()
		username := c.Sender().Username

		// Check if the user has set a language
		language := redisInstance.HGet(context.Background(), username, "language").Val()

		// Only set a new localizer if the language has changed
		if language != currentLanguage {
			var newLocalizer *i18n.Localizer

			switch language {
			case "":
				newLocalizer = i18n.NewLocalizer(bundle, "en-US")
				currentLanguage = "en"
				fmt.Println("Language not set, defaulting to English")
			case "en":
				newLocalizer = i18n.NewLocalizer(bundle, "en-US")
				currentLanguage = "en"
				fmt.Println("Set to English")
			case "vi":
				newLocalizer = i18n.NewLocalizer(bundle, "vi-VN")
				currentLanguage = "vi"
				fmt.Println("Set to Vietnamese")
			default:
				newLocalizer = i18n.NewLocalizer(bundle, "en-US")
				currentLanguage = "en"
				fmt.Println("Unexpected language value, defaulting to English")
			}

			// Only update localizer and set handler if a new localizer was created
			if newLocalizer != nil {
				localizer = newLocalizer
				SetHandlerForBot(c.Bot(), localizer, authRoute, globalLimitOrder, transferInfo)
			}
		}

		// Proceed with the next handler
		return next(c)
	}
}

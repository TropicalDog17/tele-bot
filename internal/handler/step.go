package handler

import (
	"fmt"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/TropicalDog17/tele-bot/internal/client"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	tele "gopkg.in/telebot.v3"
)

func HandleStep(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, clients map[string]internal.BotClient, utils utils.UtilsInterface, currentStep *string, menuSendToken, menuLimitOrder, menuCreateLimitOrder *tele.ReplyMarkup, globalLimitOrder *types.LimitOrderInfo, transferInfo *types.TransferInfo, globalMenu *tele.StoredMessage, createOrderMenu *tele.StoredMessage) {

	authRoute.Handle(tele.OnText, func(c tele.Context) error {
		if *currentStep == "addPassword" || *currentStep == "sendMnemonic" || *currentStep == "confirmMnemonic" || *currentStep == "receiveMnemonicWords" {
			client := client.NewTempClient()
			return HandleOnboardStep(b, c, client, clients, utils, currentStep)
		}
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		fmt.Println("Current step: ", *currentStep)
		if *currentStep == "customAmount" || *currentStep == "recipientAddress" || *currentStep == "customToken" {
			return HandleTransferStep(b, localizer, client, c, menuSendToken, transferInfo, globalMenu, currentStep)
		} else if *currentStep == "limitAmount" || *currentStep == "limitPrice" || *currentStep == "limitToken" || *currentStep == "payWithToken" {
			return HandleLimitStep(b, c, client, createOrderMenu, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, currentStep)
		} else if *currentStep == "cancelOrder" {
			return HandleCancelLimitOrderStep(b, c, client, globalLimitOrder)
		} else if *currentStep == "changeLanguage" || *currentStep == "userInputLanguage" || *currentStep == "changeCurrency" || *currentStep == "changePassword" || *currentStep == "deletePassword" {
			return HandleSettingsStep(b, localizer, c, client, currentStep)
		}
		return c.Send("Invalid input", types.Menu)
	})

}

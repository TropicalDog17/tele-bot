package handler

import (
	"github.com/TropicalDog17/tele-bot/internal"

	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/awnumar/memguard"
	tele "gopkg.in/telebot.v3"
)

var passwordChan = make(chan *memguard.LockedBuffer)

func HandleStep(b *tele.Bot, clients map[string]internal.BotClient, utils utils.UtilsInterface, currentStep *string, menuSendToken, menuLimitOrder, menuCreateLimitOrder *tele.ReplyMarkup, globalLimitOrder *types.LimitOrderInfo, selectedAmount, selectedToken, recipientAddress *string, globalMenu *tele.StoredMessage, createOrderMenu *tele.StoredMessage) {

	b.Handle(tele.OnText, func(c tele.Context) error {
		client, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		if *currentStep == "customAmount" || *currentStep == "recipentAddress" {
			return HandleTransferStep(b, client, c, menuSendToken, selectedAmount, selectedToken, recipientAddress, globalMenu, currentStep)
		} else if *currentStep == "limitAmount" || *currentStep == "limitPrice" || *currentStep == "limitToken" {
			return HandleLimitStep(b, c, createOrderMenu, menuLimitOrder, menuCreateLimitOrder, globalLimitOrder, currentStep)
		} else if *currentStep == "cancelOrder" {
			return HandleCancelLimitOrderStep(b, c, client, globalLimitOrder)
		} else if *currentStep == "addPassword" || *currentStep == "sendMnemonic" || *currentStep == "confirmMnemonic" || *currentStep == "receiveMnemonicWords" {
			return HandleOnboardStep(b, c, client, utils, currentStep)
		}
		return c.Send("Invalid input", types.Menu)
	})

}

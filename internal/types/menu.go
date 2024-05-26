package types

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

type TransferInfo struct {
	SelectedToken    string
	SelectedAmount   string
	RecipientAddress string
}

var (
	// Universal markup builders.
	Menu                 = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuSendToken        = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuLimitOrder       = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuCreateLimitOrder = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuConfirmOrder     = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuActiveOrders     = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuViewMarket       = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuViewSettings     = &tele.ReplyMarkup{ResizeKeyboard: true}
	Selector             = &tele.ReplyMarkup{}
)

func MainMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	Menu.Reply(
		Menu.Row(BtnShowAccount(localizer), BtnSettings(localizer)),
		Menu.Row(BtnSendToken(localizer), BtnLimitOrder(localizer)),
		Menu.Row(BtnViewMarket(localizer), BtnPriceAlert(localizer)),
	)
	return Menu
}

func SendTokenMenu(localizer *i18n.Localizer, info *TransferInfo) *tele.ReplyMarkup {
	MenuSendToken.Inline(
		MenuSendToken.Row(BtnBack(localizer), BtnMenu(localizer)),
		MenuSendToken.Row(BtnTokenSection(localizer, info)),
		MenuSendToken.Row(BtnInlineAtom(localizer), BtnInlineInj(localizer), BtnCustomToken(localizer)),
		MenuSendToken.Row(BtnAmountSection(localizer, info)),
		MenuSendToken.Row(BtnTenDollar(localizer), BtnFiftyDollar(localizer), BtnHundredDollar(localizer)),
		MenuSendToken.Row(BtnTwoHundredDollar(localizer), BtnFiveHundredDollar(localizer), BtnCustomAmount(localizer)),
		MenuSendToken.Row(BtnRecipientSection(localizer, info)),
		MenuSendToken.Row(BtnSend(localizer)),
	)
	switch info.SelectedToken {
	case "atom":
		MenuSendToken.InlineKeyboard = RemoveGreenTickToken(MenuSendToken.InlineKeyboard)
		fmt.Println("Selected token: ", info.SelectedToken)
		MenuSendToken.InlineKeyboard[2][0] = AddGreenTick(*BtnInlineAtom(localizer).Inline())
	case "inj":
		MenuSendToken.InlineKeyboard = RemoveGreenTickToken(MenuSendToken.InlineKeyboard)
		fmt.Println("Selected token: ", info.SelectedToken)
		MenuSendToken.InlineKeyboard[2][1] = AddGreenTick(*BtnInlineInj(localizer).Inline())
	}

	if info.SelectedAmount != "" {
		MenuSendToken.InlineKeyboard = RemoveGreenTickForAmount(MenuSendToken.InlineKeyboard)
		MenuSendToken.InlineKeyboard = ModifyAmountToTransferButton(MenuSendToken.InlineKeyboard, localizer, info)
	}
	if info.RecipientAddress != "" {
		MenuSendToken.InlineKeyboard[6][0].Text = "Recipient: " + info.RecipientAddress
	}

	return MenuSendToken
}

func LimitOrderMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuLimitOrder.Inline(
		MenuLimitOrder.Row(BtnActiveOrders(localizer)),
		MenuLimitOrder.Row(BtnBuyLimitOrder(localizer), BtnSellLimitOrder(localizer)),
		MenuLimitOrder.Row(BtnBack(localizer)),
	)
	return MenuLimitOrder
}

func CreateLimitOrderMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuCreateLimitOrder.Inline(
		MenuCreateLimitOrder.Row(BtnBack(localizer)),
		MenuCreateLimitOrder.Row(BtnToken(localizer)),
		MenuCreateLimitOrder.Row(BtnAmount(localizer)),
		MenuCreateLimitOrder.Row(BtnPayWith(localizer)),
		MenuCreateLimitOrder.Row(BtnPrice(localizer)),
		MenuCreateLimitOrder.Row(BtnConfirmOrder(localizer)),
	)
	return MenuCreateLimitOrder
}

func ActiveOrdersMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuActiveOrders.Inline(
		MenuActiveOrders.Row(BtnCancelOrder(localizer)),
	)
	return MenuActiveOrders
}

func ConfirmOrderMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuConfirmOrder.Inline(
		MenuConfirmOrder.Row(BtnConfirmLimitOrder(localizer), BtnClose(localizer)),
	)
	return MenuConfirmOrder
}

func ViewSettingsMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuViewSettings.Inline(
		MenuViewSettings.Row(BtnChangeLanguage(localizer)),
		MenuViewSettings.Row(BtnChangeDefaultLimitPair(localizer)),
		MenuViewSettings.Row(BtnDisablePassword(localizer)),
		MenuViewSettings.Row(BtnChangePassword(localizer)),
	)
	return MenuViewSettings
}

func ViewMarketMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuViewMarket.Inline(
		MenuViewMarket.Row(BtnBiggestGainer24h(localizer)),
		MenuViewMarket.Row(BtnBiggestLoser24h(localizer)),
		MenuViewMarket.Row(BtnBiggestVolume24h(localizer)),
	)
	return MenuViewMarket
}

func AddGreenTick(btn tele.InlineButton) tele.InlineButton {
	btn.Text = "✅ " + btn.Text
	return btn
}
func RemoveGreenTickToken(keyboard [][]tele.InlineButton) [][]tele.InlineButton {
	for i := 0; i < len(keyboard[2]); i++ {
		if keyboard[2][i].Text[0:3] == "✅" {
			keyboard[2][i].Text = keyboard[2][i].Text[3:]
		}
	}
	return keyboard
}
func RemoveGreenTickForAmount(keyboard [][]tele.InlineButton) [][]tele.InlineButton {
	for i := 0; i < len(keyboard[4]); i++ {

		if keyboard[4][i].Text[0:3] == "✅" {
			keyboard[4][i].Text = keyboard[4][i].Text[3:]
		}
	}
	for i := 0; i < len(keyboard[5]); i++ {
		if keyboard[5][i].Text[0:3] == "✅" {
			keyboard[5][i].Text = keyboard[5][i].Text[3:]
		}
	}
	return keyboard
}

func ModifyAmountToTransferButton(keyboard [][]tele.InlineButton, localizer *i18n.Localizer, info *TransferInfo) [][]tele.InlineButton {
	if info.SelectedToken != "" {
		keyboard[3][0].Text = localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Transfer", Other: "Transfer"}}) + " " + info.SelectedAmount + " " + info.SelectedToken
	}
	switch info.SelectedAmount {
	case "1":
		keyboard = RemoveGreenTickForAmount(keyboard)
		keyboard[4][0] = AddGreenTick(*BtnTenDollar(localizer).Inline())
	case "5":
		keyboard = RemoveGreenTickForAmount(keyboard)
		keyboard[4][1] = AddGreenTick(*BtnFiftyDollar(localizer).Inline())
	case "10":
		keyboard = RemoveGreenTickForAmount(keyboard)
		keyboard[4][2] = AddGreenTick(*BtnHundredDollar(localizer).Inline())
	case "20":
		keyboard = RemoveGreenTickForAmount(keyboard)
		keyboard[5][0] = AddGreenTick(*BtnTwoHundredDollar(localizer).Inline())
	case "50":
		keyboard = RemoveGreenTickForAmount(keyboard)
		keyboard[5][1] = AddGreenTick(*BtnFiveHundredDollar(localizer).Inline())
	default:
		keyboard[5][2].Text = "Custom amount: " + info.SelectedAmount
	}
	return keyboard
}

func ModifyCustomTokenButton(keyboard [][]tele.InlineButton, localizer *i18n.Localizer, info *TransferInfo) [][]tele.InlineButton {
	keyboard = RemoveGreenTickToken(keyboard)
	if info.SelectedToken == "atom" {
		keyboard[2][0] = AddGreenTick(*BtnInlineAtom(localizer).Inline())
	} else if info.SelectedToken == "inj" {
		keyboard[2][1] = AddGreenTick(*BtnInlineInj(localizer).Inline())
	} else {
		if info.SelectedToken != "" {
			keyboard[2][2].Text = "Custom token: " + info.SelectedToken
		}
	}

	return keyboard
}

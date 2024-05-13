package types

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	tele "gopkg.in/telebot.v3"
)

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

func SendTokenMenu(localizer *i18n.Localizer) *tele.ReplyMarkup {
	MenuSendToken.Inline(
		MenuSendToken.Row(BtnBack(localizer), BtnMenu(localizer)),
		MenuSendToken.Row(BtnTokenSection(localizer)),
		MenuSendToken.Row(BtnInlineAtom(localizer), BtnInlineInj(localizer), BtnCustomToken(localizer)),
		MenuSendToken.Row(BtnAmountSection(localizer)),
		MenuSendToken.Row(BtnTenDollar(localizer), BtnFiftyDollar(localizer), BtnHundredDollar(localizer)),
		MenuSendToken.Row(BtnTwoHundredDollar(localizer), BtnFiveHundredDollar(localizer), BtnCustomAmount(localizer)),
		MenuSendToken.Row(BtnRecipientSection(localizer)),
		MenuSendToken.Row(BtnSend(localizer)),
	)
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

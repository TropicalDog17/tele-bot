package types

import tele "gopkg.in/telebot.v3"

var (
	// Universal markup builders.
	Menu                 = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuSendToken        = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuLimitOrder       = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuCreateLimitOrder = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuConfirmOrder     = &tele.ReplyMarkup{ResizeKeyboard: true}
	MenuActiveOrders     = &tele.ReplyMarkup{ResizeKeyboard: true}
	Selector             = &tele.ReplyMarkup{}
)

func InitializeUI() []*tele.ReplyMarkup {
	Menu.Reply(
		Menu.Row(BtnViewBalances, BtnSettings),
		Menu.Row(BtnSendToken, BtnShowAccount),
		Menu.Row(BtnLimitOrder, BtnSpotOrder),
	)
	MenuSendToken.Inline(
		MenuSendToken.Row(BtnBack, BtnMenu),
		MenuSendToken.Row(BtnTokenSection),
		MenuSendToken.Row(BtnInlineAtom, BtnInlineInj),
		MenuSendToken.Row(BtnAmountSection),
		MenuSendToken.Row(BtnTenDollar, BtnFiftyDollar, BtnHundredDollar),
		MenuSendToken.Row(BtnTwoHundredDollar, BtnFiveHundredDollar, BtnCustomAmount),
		MenuSendToken.Row(BtnRecipientSection),
		MenuSendToken.Row(BtnSend),
	)
	MenuLimitOrder.Inline(
		MenuLimitOrder.Row(BtnActiveOrders),
		MenuLimitOrder.Row(BtnBuyLimitOrder, BtnSellLimitOrder),
		MenuLimitOrder.Row(BtnBack),
	)
	MenuCreateLimitOrder.Inline(
		MenuCreateLimitOrder.Row(BtnBack),
		MenuCreateLimitOrder.Row(BtnToken),
		MenuCreateLimitOrder.Row(BtnAmount),
		MenuCreateLimitOrder.Row(BtnPayWith),
		MenuCreateLimitOrder.Row(BtnPrice),
		MenuCreateLimitOrder.Row(BtnConfirmOrder),
	)
	MenuActiveOrders.Inline(
		MenuActiveOrders.Row(BtnCancelOrder),
	)
	MenuConfirmOrder.Inline(
		MenuConfirmOrder.Row(BtnConfirmLimitOrder, BtnClose),
	)

	return []*tele.ReplyMarkup{Menu, MenuSendToken, MenuLimitOrder, MenuCreateLimitOrder, MenuConfirmOrder, MenuActiveOrders}
}

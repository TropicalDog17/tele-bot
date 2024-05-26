package types

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/telebot.v3"
)

func BtnViewBalances(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ViewBalances",
			Other: "ℹ View Balances",
		},
	}))
}

func BtnSettings(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Settings",
			Other: "⚙ Settings",
		},
	}))
}

func BtnSendToken(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SendToken",
			Other: "💸 Send Token",
		},
	}))
}

func BtnShowAccount(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ShowAccount",
			Other: "👤 Show Account",
		},
	}))
}

func BtnLimitOrder(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "LimitOrder",
			Other: "🚀 Limit Order",
		},
	}))
}

func BtnSpotOrder(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SpotOrder",
			Other: "📊 Spot Order",
		},
	}))
}

func BtnViewMarket(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ViewMarket",
			Other: "📈 View Market",
		},
	}))
}

func BtnPriceAlert(localizer *i18n.Localizer) telebot.Btn {
	return Menu.Text(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PriceAlert",
			Other: "🔔 Price Alert",
		},
	}))
}

func BtnInlineAtom(localizer *i18n.Localizer) telebot.Btn {
	return Selector.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InlineAtom",
			Other: "ATOM",
		},
	}), "atom", "atom")
}

func BtnInlineInj(localizer *i18n.Localizer) telebot.Btn {
	return Selector.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "InlineInj",
			Other: "INJ",
		},
	}), "inj", "inj")
}

func BtnCustomToken(localizer *i18n.Localizer) telebot.Btn {
	return Selector.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CustomToken",
			Other: "Custom Token",
		},
	}), "customToken", "customToken")
}

func BtnMenu(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Menu",
			Other: "Menu",
		},
	}), "menu")
}

func BtnSend(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Send",
			Other: "Send",
		},
	}), "send", "send")
}

func BtnBack(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Back",
			Other: "Back",
		},
	}), "btnBack")
}

func BtnConfirmLimitOrder(localizer *i18n.Localizer) telebot.Btn {
	return MenuConfirmOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ConfirmLimitOrder",
			Other: "Confirm",
		},
	}), "confirmLimitOrder", "confirm")
}

func BtnClose(localizer *i18n.Localizer) telebot.Btn {
	return MenuConfirmOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Close",
			Other: "Close",
		},
	}), "close", "close")
}

func BtnCancelOrder(localizer *i18n.Localizer) telebot.Btn {
	return MenuActiveOrders.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CancelOrder",
			Other: "Cancel Order",
		},
	}), "cancelOrder", "cancel")
}

func BtnToken(localizer *i18n.Localizer) telebot.Btn {
	return MenuCreateLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Token",
			Other: "🪙 Token",
		},
	}), "limitToken", "token")
}

func BtnAmount(localizer *i18n.Localizer) telebot.Btn {
	return MenuCreateLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Amount",
			Other: "📈 Amount",
		},
	}), "limitAmount", "amount")
}

func BtnPrice(localizer *i18n.Localizer) telebot.Btn {
	return MenuCreateLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Price",
			Other: "💰 Price",
		},
	}), "limitPrice", "price")
}

func BtnConfirmOrder(localizer *i18n.Localizer) telebot.Btn {
	return MenuCreateLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ConfirmOrder",
			Other: "✅ Confirm Order",
		},
	}), "confirmOrder", "confirm")
}

func BtnPayWith(localizer *i18n.Localizer) telebot.Btn {
	return MenuCreateLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "PayWith",
			Other: "💳 Pay With",
		},
	}), "payWith", "payWith")
}

func BtnActiveOrders(localizer *i18n.Localizer) telebot.Btn {
	return MenuLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ActiveOrders",
			Other: "💸 Active Orders",
		},
	}), "activeOrders", "active")
}

func BtnBuyLimitOrder(localizer *i18n.Localizer) telebot.Btn {
	return MenuLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BuyLimitOrder",
			Other: "📈 Buy",
		},
	}), "buyLimit", "buy")
}

func BtnSellLimitOrder(localizer *i18n.Localizer) telebot.Btn {
	return MenuLimitOrder.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SellLimitOrder",
			Other: "📉 Sell",
		},
	}), "sellLimit", "sell")
}

//////////////////////////
// Send Token Section
//////////////////////////

func BtnTokenSection(localizer *i18n.Localizer, info *TransferInfo) telebot.Btn {

	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TokenSection",
			Other: "🪙 Token: ",
		},
	}), "tokenSection")
}

func BtnAmountSection(localizer *i18n.Localizer, info *TransferInfo) telebot.Btn {

	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "AmountSection",
			Other: "💰 Amount: ",
		},
	}), "amountSection")

}

func BtnRecipientSection(localizer *i18n.Localizer, info *TransferInfo) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "RecipientSection",
			Other: "👤 Recipient: ",
		},
	}), "recipientSection")

}

func BtnTenDollar(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TenDollar",
			Other: "💵 $1",
		},
	}), "btnTenDollar", "1")
}

func BtnFiftyDollar(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "FiftyDollar",
			Other: "💰 $5",
		},
	}), "btnFiftyDollar", "5")
}

func BtnHundredDollar(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "HundredDollar",
			Other: "💸 $10",
		},
	}), "btnHundredDollar", "10")
}

func BtnTwoHundredDollar(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "TwoHundredDollar",
			Other: "🪙 $20",
		},
	}), "btnTwoHundredDollar", "20")
}

func BtnFiveHundredDollar(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "FiveHundredDollar",
			Other: "💶 $50",
		},
	}), "btnFiveHundredDollar", "50")
}

func BtnCustomAmount(localizer *i18n.Localizer) telebot.Btn {
	return MenuSendToken.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CustomAmount",
			Other: "🎛️ Custom Amount",
		},
	}), "btnCustomAmount", "")
}

func BtnBiggestVolume24h(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewMarket.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BiggestVolume24h",
			Other: "📊 Biggest Volume 24h",
		},
	}), "biggestVolume24h", "biggestVolume24h")
}

func BtnBiggestGainer24h(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewMarket.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BiggestGainer24h",
			Other: "🚀 Biggest Gainer 24h",
		},
	}), "biggestGainer24h", "biggestGainer24h")
}

func BtnBiggestLoser24h(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewMarket.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "BiggestLoser24h",
			Other: "📉 Biggest Loser 24h",
		},
	}), "biggestLoser24h", "biggestLoser24h")
}

func BtnChangeLanguage(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewSettings.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangeLanguage",
			Other: "🌐 Change Language",
		},
	}), "changeLanguage", "changeLanguage")
}

func BtnChangeDefaultLimitPair(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewSettings.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangeDefaultLimitPair",
			Other: "🪙 Change Default Limit Pair",
		},
	}), "changeDefaultLimitPair", "changeDefaultLimitPair")
}

func BtnChangePassword(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewSettings.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ChangePassword",
			Other: "🔑 Change Password",
		},
	}), "changePassword", "changePassword")
}

func BtnDisablePassword(localizer *i18n.Localizer) telebot.Btn {
	return MenuViewSettings.Data(localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "DisablePassword",
			Other: "🔒 Disable Password",
		},
	}), "disablePassword", "disablePassword")
}

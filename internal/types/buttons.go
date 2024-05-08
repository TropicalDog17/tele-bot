package types

var (
	// Reply buttons.
	BtnViewBalances = Menu.Text("ℹ View Balances")
	BtnSettings     = Menu.Text("⚙ Settings")
	BtnSendToken    = Menu.Text("💸 Send Token")
	BtnShowAccount  = Menu.Text("👤 Show Account")
	BtnLimitOrder   = Menu.Text("🚀 Limit Order")
	BtnSpotOrder    = Menu.Text("📊 Spot Order")
	BtnViewMarket   = Menu.Text("📈 View Market")
	BtnPriceAlert   = Menu.Text("🔔 Price Alert")
	BtnInlineAtom   = Selector.Data("ATOM", "atom", "atom")
	BtnInlineInj    = Selector.Data("INJ", "inj", "inj")
	BtnCustomToken  = Selector.Data("Custom Token", "customToken", "customToken")
	BtnMenu         = MenuSendToken.Data("Menu", "menu")

	BtnSend              = MenuSendToken.Data("Send", "send", "send")
	BtnBack              = MenuSendToken.Data("Back", "btnBack")
	BtnToken             = MenuCreateLimitOrder.Data("Token", "limitToken", "token")
	BtnAmount            = MenuCreateLimitOrder.Data("Amount", "limitAmount", "amount")
	BtnPrice             = MenuCreateLimitOrder.Data("Price", "limitPrice", "price")
	BtnConfirmOrder      = MenuCreateLimitOrder.Data("Confirm Order", "confirmOrder", "confirm")
	BtnConfirmLimitOrder = MenuConfirmOrder.Data("Confirm", "confirmLimitOrder", "confirm")
	BtnClose             = MenuConfirmOrder.Data("Close", "close", "close")
	BtnCancelOrder       = MenuActiveOrders.Data("Cancel Order", "cancelOrder", "cancel")
	BtnPayWith           = MenuCreateLimitOrder.Data("Pay With", "payWith", "payWith")
)

// Limit order buttons.
var (
	BtnActiveOrders   = MenuLimitOrder.Data("💸 Active Orders", "activeOrders", "active")
	BtnBuyLimitOrder  = MenuLimitOrder.Data("📈 Buy", "buyLimit", "buy")
	BtnSellLimitOrder = MenuLimitOrder.Data("📉 Sell", "sellLimit", "sell")
)

// Send token buttons.
var (
	BtnTokenSection     = MenuSendToken.Data("💠 Token Section 💠 ", "tokenSection")
	BtnAmountSection    = MenuSendToken.Data("💰 Amount Section 💰", "amountSection")
	BtnRecipientSection = MenuSendToken.Data("📨 Enter Recipient Address: 📨 ", "recipient", "recipient")
)

// Amount button
var (
	BtnTenDollar         = MenuSendToken.Data("💵 $10", "btnTenDollar", "10")
	BtnFiftyDollar       = MenuSendToken.Data("💰 $50", "btnFiftyDollar", "50")
	BtnHundredDollar     = MenuSendToken.Data("💸 $100", "btnHundredDollar", "100")
	BtnTwoHundredDollar  = MenuSendToken.Data("🪙 $200", "btnTwoHundredDollar", "200")
	BtnFiveHundredDollar = MenuSendToken.Data("💶 $500", "btnFiveHundredDollar", "500")
	BtnCustomAmount      = MenuSendToken.Data("🎛️ Custom Amount", "btnCustomAmount", "")
)

// View market button

var (
	BtnBiggestVolume24h = MenuViewMarket.Data("📊 Biggest Volume 24h", "biggestVolume24h", "biggestVolume24h")
	BtnBiggestGainer24h = MenuViewMarket.Data("🚀 Biggest Gainer 24h", "biggestGainer24h", "biggestGainer24h")
	BtnBiggestLoser24h  = MenuViewMarket.Data("📉 Biggest Loser 24h", "biggestLoser24h", "biggestLoser24h")
)

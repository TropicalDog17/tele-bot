package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/handler"
	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)

var (
	selectedToken    string
	selectedAmount   string
	recipientAddress string
	currentStep      string
	globalMenu       tele.StoredMessage
	limitOrderMenu   tele.StoredMessage
	createOrderMenu  tele.StoredMessage
	promptMsg        tele.StoredMessage
)

var (
	// Universal markup builders.
	menu                 = &tele.ReplyMarkup{ResizeKeyboard: true}
	menuSendToken        = &tele.ReplyMarkup{ResizeKeyboard: true}
	menuLimitOrder       = &tele.ReplyMarkup{ResizeKeyboard: true}
	menuCreateLimitOrder = &tele.ReplyMarkup{ResizeKeyboard: true}
	menuConfirmOrder     = &tele.ReplyMarkup{ResizeKeyboard: true}
	menuActiveOrders     = &tele.ReplyMarkup{ResizeKeyboard: true}
	selector             = &tele.ReplyMarkup{}
	// Reply buttons.
	btnViewBalances = menu.Text("ℹ View Balances")
	btnSettings     = menu.Text("⚙ Settings")
	btnSendToken    = menu.Text("💸 Send Token")
	btnShowAccount  = menu.Text("👤 Show Account")
	btnLimitOrder   = menu.Text("📈 Limit Order")
	btnSpotOrder    = menu.Text("📊 Spot Order")

	btnInlineAtom        = selector.Data("ATOM", "atom", "atom")
	btnInlineInj         = selector.Data("INJ", "inj", "inj")
	btnMenu              = menuSendToken.Data("Menu", "menu")
	btnTokenSection      = menuSendToken.Data("---Token Section---", "tokenSection")
	btnAmountSection     = menuSendToken.Data("---Amount Section---", "amountSection")
	btnRecipientSection  = menuSendToken.Data("Enter Recipient Address:", "recipient", "recipient")
	btnSend              = menuSendToken.Data("Send", "send", "send")
	btnTenDollar         = menuSendToken.Data("$10", "btnTenDollar", "10")
	btnFiftyDollar       = menuSendToken.Data("$50", "btnFiftyDollar", "50")
	btnHundredDollar     = menuSendToken.Data("$100", "btnHundredDollar", "100")
	btnTwoHundredDollar  = menuSendToken.Data("$200", "btnTwoHundredDollar", "200")
	btnFiveHundredDollar = menuSendToken.Data("$500", "btnFiveHundredDollar", "500")
	btnCustomAmount      = menuSendToken.Data("Custom Amount", "btnCustomAmount", "")
	btnBack              = menuSendToken.Data("Back", "btnBack")
	btnBuyLimitOrder     = menuLimitOrder.Data("📈 Buy", "buyLimit", "buy")
	btnSellLimitOrder    = menuLimitOrder.Data("📉 Sell", "sellLimit", "sell")
	btnActiveOrders      = menuLimitOrder.Data("💸 Active Orders", "activeOrders", "active")
	btnToken             = menuCreateLimitOrder.Data("Token", "limitToken", "token")
	btnAmount            = menuCreateLimitOrder.Data("Amount", "limitAmount", "amount")
	btnPrice             = menuCreateLimitOrder.Data("Price", "limitPrice", "price")
	btnConfirmOrder      = menuCreateLimitOrder.Data("Confirm Order", "confirmOrder", "confirm")
	btnConfirmLimitOrder = menuConfirmOrder.Data("Confirm", "confirmLimitOrder", "confirm")
	btnClose             = menuConfirmOrder.Data("Close", "close", "close")
)

var (
	globalLimitOrder = types.NewLimitOrderInfo()
)

func main() {
	client := internal.NewClient()

	menu.Reply(
		menu.Row(btnViewBalances, btnSettings),
		menu.Row(btnSendToken, btnShowAccount),
		menu.Row(btnLimitOrder, btnSpotOrder),
	)
	menuSendToken.Inline(
		menuSendToken.Row(btnBack, btnMenu),
		menuSendToken.Row(btnTokenSection),
		menuSendToken.Row(btnInlineAtom, btnInlineInj),
		menuSendToken.Row(btnAmountSection),
		menuSendToken.Row(btnTenDollar, btnFiftyDollar, btnHundredDollar),
		menuSendToken.Row(btnTwoHundredDollar, btnFiveHundredDollar, btnCustomAmount),
		menuSendToken.Row(btnRecipientSection),
		menuSendToken.Row(btnSend),
	)
	menuLimitOrder.Inline(
		menuLimitOrder.Row(btnActiveOrders),
		menuLimitOrder.Row(btnBuyLimitOrder, btnSellLimitOrder),
		menuLimitOrder.Row(btnBack),
	)
	menuCreateLimitOrder.Inline(
		menuCreateLimitOrder.Row(btnBack),
		menuCreateLimitOrder.Row(btnToken),
		menuCreateLimitOrder.Row(btnAmount),
		menuCreateLimitOrder.Row(btnPrice),
		menuCreateLimitOrder.Row(btnConfirmOrder),
	)
	menuConfirmOrder.Inline(
		menuConfirmOrder.Row(btnConfirmLimitOrder, btnClose),
	)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	fmt.Println(pref.Token)
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!", menu)
	})
	// On reply button pressed (message)
	b.Handle(&btnViewBalances, func(c tele.Context) error {
		fmt.Println("Button pressed")
		balances, err := client.GetBalances("inj1gxv7rs9q60qyjtaxrmu0pgwvatm6smyk4cz9d0", []string{"atom", "inj"})
		if err != nil {
			return c.Send("Error fetching balances")
		}
		return c.Send(fmt.Sprintf("Balances: %v", balances))
	})

	b.Handle("/menu", func(c tele.Context) error {
		return c.Send("Menu", menu)
	})
	// Show account
	b.Handle(&btnShowAccount, func(c tele.Context) error {
		accountDetails := &tele.ReplyMarkup{}
		address := client.GetAddress()
		balances, err := client.GetBalances(address, []string{"atom", "inj"})
		if err != nil {
			return c.Send("Error fetching balances")
		}
		rows := []tele.Row{}
		for denom, balance := range balances {
			usdPrice, found := client.GetPrice(denom)
			var balanceInUsd float64
			if !found {
				balanceInUsd = 0
			} else {
				balanceInUsd = balance * usdPrice
			}
			rows = append(rows, accountDetails.Row(accountDetails.Data(fmt.Sprintf("%s: %.3f %.3f", denom, balance, balanceInUsd), "balance", "balance")))
		}
		rows = append(rows, accountDetails.Row(accountDetails.Data("Show QR for address", "qr", "qr")))
		accountDetails.Inline(rows...)
		// Message contain the account address

		return c.Send("Account: "+address, accountDetails)
	})
	handler.HandleAddressQr(b, client)
	// Handle the "Send Tokens" button click
	b.Handle(&btnSendToken, func(c tele.Context) error {
		msg, err := b.Send(c.Chat(), "Select the token to send:", menuSendToken)
		if err != nil {
			return err
		}
		globalMenu.ChatID = msg.Chat.ID
		globalMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/sendTokenMenu.txt", []byte(fmt.Sprintf("%d %s", globalMenu.ChatID, globalMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}
		return nil
	})

	b.Handle(&btnLimitOrder, func(c tele.Context) error {
		text := "📊 Limit Orders\n\nBuy or Sell tokens automatically at your desired price.\n1. Choose to Buy or Sell.\n2. Choose the Token to Buy or Sell.\n3. Select the amount to Buy or Sell.\n4. Set your target buy or sell price.\n5. Pick an expiry time for the order.\n6. Click Create Order and Review, Confirm.\n\nTo manage or view unfilled orders, click Active Orders."
		msg, err := b.Send(c.Chat(), text, menuLimitOrder)
		if err != nil {
			return err
		}
		limitOrderMenu.ChatID = msg.Chat.ID
		limitOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)

		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/limitOrderMenu.txt", []byte(fmt.Sprintf("%d %s", limitOrderMenu.ChatID, limitOrderMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}
		return nil
	})
	b.Handle(&btnBuyLimitOrder, func(ctx tele.Context) error {
		text := "Place a buy limit order"
		menuCreateLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
		msg, err := b.Send(ctx.Chat(), text, menuCreateLimitOrder)
		if err != nil {
			return err
		}
		createOrderMenu.ChatID = msg.Chat.ID
		createOrderMenu.MessageID = fmt.Sprintf("%d", msg.ID)
		// Store chat ID and message ID in a file for future reference
		err = os.WriteFile("db/createOrderMenu.txt", []byte(fmt.Sprintf("%d %s", createOrderMenu.ChatID, createOrderMenu.MessageID)), fs.FileMode(0644))
		if err != nil {
			return err
		}

		return nil
	})
	b.Handle(&btnAmount, func(c tele.Context) error {
		currentStep = "limitAmount"
		return c.Send("Enter the amount to buy", tele.ForceReply)
	})
	b.Handle(&btnPrice, func(c tele.Context) error {
		currentStep = "limitPrice"
		return c.Send("Enter the price to buy", tele.ForceReply)
	})
	b.Handle(&btnToken, func(c tele.Context) error {
		currentStep = "limitToken"
		return c.Send("Enter the token to buy", tele.ForceReply)
	})
	b.Handle(&btnConfirmOrder, func(c tele.Context) error {
		currentStep = "confirmOrder"
		orderOverview := client.ToMessage(*globalLimitOrder)
		return c.Send(orderOverview, menuConfirmOrder)
	})
	b.Handle(&btnConfirmLimitOrder, func(c tele.Context) error {
		// TODO: Perform the limit order logic here
		txHash, err := client.PlaceSpotOrder(globalLimitOrder.DenomIn, globalLimitOrder.DenomOut, globalLimitOrder.Amount, globalLimitOrder.Price)
		if err != nil {
			return c.Send("Error placing limit order", menu)
		}
		return c.Send("Successfully send order, check txhash here: "+txHash, menu)
	})
	b.Handle(&btnBack, func(c tele.Context) error {
		if currentStep == "confirmOrder" {
			currentStep = ""
		}
		return c.Send("Back to main menu", menu)
	})

	// Handle active orders
	b.Handle(&btnActiveOrders, func(c tele.Context) error {
		// TODO: fix hardcoded market id
		orders, err := client.GetActiveOrders("0xfbd55f13641acbb6e69d7b59eb335dabe2ecbfea136082ce2eedaba8a0c917a3")
		if err != nil {
			return c.Send("Error fetching active orders")
		}
		msgs := []string{}
		if len(orders) == 0 {
			return c.Send("No active orders", menu)
		}
		for _, order := range orders {
			msgs = append(msgs, client.ToMessage(order))
		}

		return c.Send(strings.Join(msgs, "\n\n"), menu)
	})

	// Handle inline button clicks for token selection
	b.Handle(&btnInlineAtom, func(c tele.Context) error {
		selectedToken = "atom"
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][0] = internal.AddGreenTick(*btnInlineAtom.Inline())
		return c.Edit("Selected token: ATOM", menuSendToken)
	})

	b.Handle(&btnInlineInj, func(c tele.Context) error {
		selectedToken = "inj"
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickToken(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[2][1] = internal.AddGreenTick(*btnInlineInj.Inline())
		return c.Edit("Selected token: INJ", menuSendToken)
	})

	// Handle amount button clicks
	b.Handle(&btnTenDollar, func(c tele.Context) error {
		selectedAmount = "10"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][0] = internal.AddGreenTick(*btnTenDollar.Inline())
		return c.Edit("Selected amount: $10", menuSendToken)
	})

	b.Handle(&btnFiftyDollar, func(c tele.Context) error {
		selectedAmount = "50"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][1] = internal.AddGreenTick(*btnFiftyDollar.Inline())

		return c.Edit("Selected amount: $50", menuSendToken)
	})

	b.Handle(&btnHundredDollar, func(c tele.Context) error {
		selectedAmount = "100"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[4][2] = internal.AddGreenTick(*btnHundredDollar.Inline())
		return c.Edit("Selected amount: $100", menuSendToken)
	})

	b.Handle(&btnTwoHundredDollar, func(c tele.Context) error {
		selectedAmount = "200"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][0] = internal.AddGreenTick(*btnTwoHundredDollar.Inline())
		return c.Edit("Selected amount: $200", menuSendToken)
	})

	b.Handle(&btnFiveHundredDollar, func(c tele.Context) error {
		selectedAmount = "500"
		menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
		menuSendToken.InlineKeyboard = internal.RemoveGreenTickForAmount(menuSendToken.InlineKeyboard)
		menuSendToken.InlineKeyboard[5][1] = internal.AddGreenTick(*btnFiveHundredDollar.Inline())
		return c.Edit("Selected amount: $500", menuSendToken)
	})

	b.Handle(&btnCustomAmount, func(c tele.Context) error {
		// Prompt the user to enter a custom amount
		currentStep = "customAmount"
		return c.Send("Please enter the custom amount:")
	})

	b.Handle(&btnRecipientSection, func(c tele.Context) error {
		// Prompt the user to enter a recipient address
		currentStep = "recipientAddress"
		return c.Send("Please enter the recipient address:", tele.ForceReply)
	})
	b.Handle(tele.OnText, func(c tele.Context) error {
		// Check if the user is entering a custom amount
		if currentStep == "customAmount" {
			selectedAmount = c.Text()
			menuSendToken.InlineKeyboard = internal.ModifyAmountToTransferButton(menuSendToken.InlineKeyboard, selectedAmount, selectedToken)
			return c.Send(fmt.Sprintf("Selected amount: %s", selectedAmount), menuSendToken)
		} else if currentStep == "recipientAddress" { // Check if the user is entering a recipient addres
			recipientAddress = c.Text()
			fmt.Println("Recipient address: ", recipientAddress)
			err := b.Delete(c.Message().ReplyTo)
			if err != nil {
				return err
			}
			err = c.Delete()
			if err != nil {
				return err
			}
			btnRecipientSection.Text = "Recipient:" + recipientAddress
			menuSendToken.InlineKeyboard[6][0] = *btnRecipientSection.Inline()

			// load the global menu from the file
			data, err := os.ReadFile("db/sendTokenMenu.txt")
			if err != nil {
				return err
			}
			_, err = fmt.Sscanf(string(data), "%d %s", &globalMenu.ChatID, &globalMenu.MessageID)
			if err != nil {
				return err
			}
			_, err = b.EditReplyMarkup(&globalMenu, menuSendToken)
			if err != nil {
				return err
			}
			return nil
		} else if currentStep == "limitAmount" {
			globalLimitOrder.Amount, err = strconv.ParseFloat(c.Text(), 64)
			if err != nil {
				return c.Send("Invalid amount")
			}
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err := b.EditReplyMarkup(&createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		} else if currentStep == "limitPrice" {
			globalLimitOrder.Price, err = strconv.ParseFloat(c.Text(), 64)
			if err != nil {
				return c.Send("Invalid price")
			}
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err := b.EditReplyMarkup(&createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		} else if currentStep == "limitToken" {
			globalLimitOrder.DenomOut = c.Text()
			menuLimitOrder.InlineKeyboard = internal.ModifyLimitOrderMenu(menuCreateLimitOrder.InlineKeyboard, globalLimitOrder)
			_, err := b.EditReplyMarkup(&createOrderMenu, menuCreateLimitOrder)
			if err != nil {
				return err
			}
			return internal.DeleteInputMessage(b, c)
		}

		return nil
	})

	// Handle the "Send" button click
	b.Handle(&btnSend, func(c tele.Context) error {
		// Sanity check to ensure all required fields are filled
		if selectedToken == "" || selectedAmount == "" || recipientAddress == "" {
			return c.Send("Please fill in all required fields", menuSendToken)
		}
		selectedAmount, err := strconv.ParseFloat(selectedAmount, 64)
		fmt.Println(selectedAmount)
		if err != nil {
			return c.Send("Invalid amount", menuSendToken)
		}
		// Trim whitespace from the recipient address
		recipientAddress = strings.TrimSpace(recipientAddress)
		txHash, err := client.TransferToken(recipientAddress, selectedAmount/100, selectedToken)
		if err != nil {
			return c.Send("Error sending token", menuSendToken)
		}

		// TODO: Perform the token sending logic here
		// Use the selected token, amount, and recipient address
		return c.Send(fmt.Sprintf("Sent %f %s to %s, with tx hash %s", selectedAmount, selectedToken, recipientAddress, txHash), menu)
	})

	b.Start()
}

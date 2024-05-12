package handler

import (
	"context"
	"fmt"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/TropicalDog17/tele-bot/internal/client"
	"github.com/TropicalDog17/tele-bot/internal/database"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/awnumar/memguard"
	tele "gopkg.in/telebot.v3"
)

var btnOnboard = tele.ReplyMarkup{}
var btnConfirmMnemonic = btnOnboard.Data("Confirm Mnemonic", "confirm_mnemonic")
var randomIndexes [3]int

// in-memory storage, for confirmation step
var mnemonic *memguard.LockedBuffer
var password *memguard.LockedBuffer

func HandleOnboard(b internal.Bot, clients map[string]internal.BotClient, currentStep *string) {
	btnOnboard.Inline(
		btnOnboard.Row(btnConfirmMnemonic),
	)
	b.Handle("/start", func(c tele.Context) error {
		client := clients[c.Message().Sender.Username]
		return HandleStart(c, client, currentStep)
	})

}

func HandleStart(c tele.Context, botClient internal.BotClient, step *string) error {
	*step = "addPassword"
	text := "Welcome to the TropicalDog17 bot! 🐶\n\nI am a bot that can help you with your trading needs. I can provide you with the latest cryptocurrency prices, help you place limit orders, and more.\n\n To start, please provide a password"
	fmt.Println("step: ", *step)
	return c.Reply(text)

}

func HandleOnboardStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, clients map[string]internal.BotClient, utils utils.UtilsInterface, step *string) error {
	switch *step {
	case "addPassword":
		HandleStorePassword(b, c, botClient, utils, step)
	case "sendMnemonic":
		HandleSendMnemonicStep(b, c, botClient, step)
	case "confirmMnemonic":
		return HandleConfirmMnemonicStep(b, c, botClient, mnemonic, step)
	case "receiveMnemonicWords":
		HandleReceiveMnemonicWords(b, c, clients, utils, mnemonic, randomIndexes, step)
	}
	return nil
}

func HandleAddPassword(b *tele.Bot, c tele.Context, step *string) {
	*step = "addPassword"
	_, _ = b.Send(c.Chat(), "Please enter your password")
}

func HandleStorePassword(b internal.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, step *string) {
	*step = "sendMnemonic"
	// Generate random mnemonic - 24 words
	randomMnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		_ = c.Reply("Error generating mnemonic")
		return
	}
	mnemonic = memguard.NewBufferFromBytes([]byte(randomMnemonic))
	password = memguard.NewBufferFromBytes([]byte(c.Text()))
	msg1, err := HandleStorePrivateKey(b, c, botClient, utils, mnemonic, password, step)
	if err != nil {
		_ = c.Reply("Error storing private key")
	}
	b.Handle(&btnConfirmMnemonic, func(c tele.Context) error {
		_ = b.Delete(msg1)
		return HandleConfirmMnemonicStep(b, c, botClient, mnemonic, step)
	})
}

func HandleStorePrivateKey(b internal.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, mnemonic, password *memguard.LockedBuffer, step *string) (*tele.Message, error) {
	encryptedMnemonic, salt, err := utils.GetEncryptedMnemonic(mnemonic.String(), password.String())

	// Destroy password from memory
	password.Destroy()
	if err != nil {
		_, _ = b.Send(c.Chat(), "Error getting encrypted mnemonic")
		return nil, err
	}

	redisClient := botClient.GetRedisInstance()
	ctx := context.Background()
	// Store encrypted mnemonic and salt in Redis
	redisClient.HSet(ctx, c.Message().Sender.Username, "encryptedMnemonic", encryptedMnemonic)
	redisClient.HSet(ctx, c.Message().Sender.Username, "salt", salt)
	msg1, err := b.Send(c.Chat(), "Here is the mnemonic:\n"+mnemonic.String()+"\nPlease store it in a safe place, as it will be used to recover your account")

	if err != nil {
		_ = c.Reply("Error sending mnemonic")
	}

	err = c.Send("Please store it and keep in a safe place", &btnOnboard)
	if err != nil {
		_ = c.Reply("Error storing private key")
	}
	return msg1, err
}

func HandleSendMnemonicStep(b internal.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	*step = "confirmMnemonic"
	_, _ = b.Send(c.Chat(), "Mnemonic sent! Now go to confirmation part!")
}

// TODO: implement confirmation mnemonic
func HandleConfirmMnemonicStep(b internal.Bot, c tele.Context, botClient internal.BotClient, mnemonic *memguard.LockedBuffer, step *string) error {
	*step = "receiveMnemonicWords"
	randomIndexes = utils.GetRandomIndexesForTesting(len(utils.SplitMnemonic(mnemonic.String())))
	text := fmt.Sprintf("please enter these missing words, seperate by space \n %s", utils.GenerateMissedWordsMnemonicFromIndexes(mnemonic.String(), randomIndexes))
	return c.Reply(text)
}

func HandleReceiveMnemonicWords(b internal.Bot, c tele.Context, clients map[string]internal.BotClient, utils utils.UtilsInterface, mnemonic *memguard.LockedBuffer, randomIndexes [3]int, step *string) {
	providedWords := utils.SplitMnemonic(c.Text())
	if len(providedWords) != 3 {
		_, _ = b.Send(c.Chat(), "Please provide 3 words")
		*step = "confirmMnemonic"
		return
	}

	result, _ := utils.MnemonicChallenge(mnemonic.String(), randomIndexes, [3]string{providedWords[0], providedWords[1], providedWords[2]})
	if result {
		_ = c.Reply("Mnemonic confirmed!")
		// hooks
		_ = AfterMnemonicConfirmed(b, c, clients, mnemonic, step)
	} else {
		_ = c.Reply("Mnemonic not confirmed, please try again")
		*step = "confirmMnemonic"
	}
}

// delete mnemonic from memory
func AfterMnemonicConfirmed(b internal.Bot, c tele.Context, clients map[string]internal.BotClient, mnemonic *memguard.LockedBuffer, step *string) error {
	privateKey, err := utils.DerivePrivateKeyFromMnemonic(mnemonic.String())
	privateKeyBuffer := memguard.NewBufferFromBytes([]byte(utils.ECDSAToString(privateKey)))
	if err != nil {
		_ = c.Reply("Error deriving private key from mnemonic")
		return err
	}
	*step = ""
	username := c.Message().Sender.Username
	redisInstance := database.NewRedisInstance()
	client, err := client.NewClientFromPrivateKey(b, username, privateKeyBuffer, redisInstance, step)
	privateKeyBuffer.Destroy()
	if err != nil {
		return c.Reply("Error creating client" + err.Error())
	}
	clients[username] = client

	mnemonic.Destroy()

	return c.Reply("Account created successfully! You can now start using the bot. Type /help to see a list of available commands or /menu to see the main menu")
}

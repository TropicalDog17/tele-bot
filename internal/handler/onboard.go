package handler

import (
	"context"
	"fmt"

	"github.com/TropicalDog17/tele-bot/internal"
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

func HandleOnboard(b internal.Bot, client internal.BotClient, currentStep *string) {
	btnOnboard.Inline(
		btnOnboard.Row(btnConfirmMnemonic),
	)
	b.Handle("/start", func(c tele.Context) error {
		return HandleStart(c, client, currentStep)
	})

}

func HandleStart(c tele.Context, botClient internal.BotClient, step *string) error {
	*step = "addPassword"
	text := "Welcome to the TropicalDog17 bot! 🐶\n\nI am a bot that can help you with your trading needs. I can provide you with the latest cryptocurrency prices, help you place limit orders, and more.\n\nTo get started, type /help to see a list of available commands.\n To start, please provide a password"
	return c.Reply(text, &btnOnboard, tele.ModeHTML)

}

func HandleOnboardStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, step *string) error {
	switch *step {
	case "addPassword":
		HandleStorePassword(b, c, botClient, utils, step)
	case "sendMnemonic":
		HandleSendMnemonicStep(b, c, botClient, step)
	case "confirmMnemonic":
		HandleConfirmMnemonicStep(b, c, botClient, mnemonic.String(), step)
	case "receiveMnemonicWords":
		HandleReceiveMnemonicWords(b, c, botClient, utils, mnemonic.String(), randomIndexes, step)
	}
	return nil
}

func HandleAddPassword(b *tele.Bot, c tele.Context, step *string) {
	*step = "addPassword"
	_, _ = b.Send(c.Chat(), "Please enter your password")
}

func HandleStorePassword(b internal.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, step *string) {
	*step = "sendMnemonic"
	mnemonic = memguard.NewBufferFromBytes([]byte("pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"))
	password = memguard.NewBufferFromBytes([]byte(c.Text()))
	HandleStorePrivateKey(b, c, botClient, utils, mnemonic, password, step)
	HandleConfirmMnemonicStep(b, c, botClient, mnemonic.String(), step)
}

func HandleStorePrivateKey(b internal.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, mnemonic, password *memguard.LockedBuffer, step *string) {
	// Generate mnemonic
	// TODO: generate this randomly, securely
	// mnemonic := memguard.NewBufferFromBytes([]byte("pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"))
	encryptedMnemonic, salt, err := utils.GetEncryptedMnemonic(mnemonic.String(), password.String())

	// Destroy password from memory
	password.Destroy()
	if err != nil {
		_, _ = b.Send(c.Chat(), "Error getting encrypted mnemonic")
		return
	}

	redisClient := botClient.GetRedisInstance()
	ctx := context.Background()
	// Store encrypted mnemonic and salt in Redis
	redisClient.HSet(ctx, c.Message().Sender.Username, "encryptedMnemonic", encryptedMnemonic)
	redisClient.HSet(ctx, c.Message().Sender.Username, "salt", salt)
	err = c.Send("Mnemonic stored! Please confirm it", &btnConfirmMnemonic)
	if err != nil {
		_ = c.Send("Error storing private key")
	}
}

func HandleSendMnemonicStep(b internal.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	*step = "confirmMnemonic"
	_, _ = b.Send(c.Chat(), "Mnemonic sent! Now go to confirmation part!")
}

// TODO: implement confirmation mnemonic
func HandleConfirmMnemonicStep(b internal.Bot, c tele.Context, botClient internal.BotClient, mnemonic string, step *string) {
	*step = "receiveMnemonicWords"
	randomIndexes = utils.GetRandomIndexesForTesting(len(utils.SplitMnemonic(mnemonic)))
	text := fmt.Sprintf("please enter these missing words, seperate by space \n %s", utils.GenerateMissedWordsMnemonicFromIndexes(mnemonic, randomIndexes))
	_ = c.Send(text)
}

func HandleReceiveMnemonicWords(b internal.Bot, c tele.Context, botClient internal.BotClient, utils utils.UtilsInterface, mnemonic string, randomIndexes [3]int, step *string) {
	providedWords := utils.SplitMnemonic(c.Text())
	if len(providedWords) != 3 {
		_, _ = b.Send(c.Chat(), "Please provide 3 words")
		*step = "confirmMnemonic"
		return
	}

	result, _ := utils.MnemonicChallenge(mnemonic, randomIndexes, [3]string{providedWords[0], providedWords[1], providedWords[2]})
	if result {
		_ = c.Send("Mnemonic confirmed!")
	} else {
		_ = c.Send("Mnemonic not confirmed, please try again")
		*step = "confirmMnemonic"
	}
}
func HandleStoreMnemonicStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	_, _ = b.Send(c.Chat(), "Mnemonic stored!")
}

// delete mnemonic from memory
func AfterMnemonicConfirmed(b *tele.Bot, c tele.Context, exchangeClient internal.ExchangeClient, step *string) {
	privateKey, err := utils.DerivePrivateKeyFromMnemonic(mnemonic.String())
	if err != nil {
		_, _ = b.Send(c.Chat(), "Error deriving private key from mnemonic")
		return
	}
	*step = "addPassword"
	exchangeClient.GetChainClient().AdjustKeyringFromPrivateKey(utils.ECDSAToString(privateKey))
	mnemonic.Destroy()

	_, _ = b.Send(c.Chat(), "Private key derived from mnemonic and set as default keyring!")
}

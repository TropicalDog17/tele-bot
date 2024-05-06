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

// in-memory storage for mnemonic, for confirmation step
var mnemonic *memguard.LockedBuffer

func HandleOnboard(b internal.Bot, client internal.BotClient, currentStep *string) {
	btnOnboard.Inline(
		btnOnboard.Row(btnConfirmMnemonic),
	)
	b.Handle("/start", func(c tele.Context) error {
		*currentStep = "addPassword"
		text := "Welcome to the TropicalDog17 bot! 🐶\n\nI am a bot that can help you with your trading needs. I can provide you with the latest cryptocurrency prices, help you place limit orders, and more.\n\nTo get started, type /help to see a list of available commands.\n To start, please provide a password"
		return c.Reply(text, &btnOnboard, tele.ModeHTML)
	})

}

func HandleOnboardStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) error {
	switch *step {
	case "addPassword":
		HandleStorePassword(b, c, botClient, step)
	case "sendMnemonic":
		HandleSendMnemonicStep(b, c, botClient, step)
	case "confirmMnemonic":
		HandleConfirmMnemonicStep(b, c, botClient, step)
	case "receiveMnemonicWords":
		HandleReceiveMnemonicWords(b, c, botClient, step)
	}
	return nil
}

func HandleAddPassword(b *tele.Bot, c tele.Context, step *string) {
	*step = "addPassword"
	_, _ = b.Send(c.Chat(), "Please enter your password")
}

func HandleStorePassword(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	*step = "sendMnemonic"
	redisClient := botClient.GetRedisInstance()
	ctx := context.Background()
	redisClient.HSet(ctx, c.Message().Sender.Username, "password", c.Text())
	HandleStorePrivateKey(b, c, botClient, step)
	HandleConfirmMnemonicStep(b, c, botClient, step)
}

func HandleStorePrivateKey(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	// Generate mnemonic
	mnemonic := memguard.NewBufferFromBytes([]byte("pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"))
	password := "12345678"

	encryptedMnemonic, salt, err := utils.GetEncryptedMnemonic(mnemonic.String(), password)

	if err != nil {
		_, _ = b.Send(c.Chat(), "Error getting encrypted mnemonic")
		return
	}

	redisClient := botClient.GetRedisInstance()
	ctx := context.Background()

	// Store encrypted mnemonic and salt in Redis
	redisClient.HSet(ctx, c.Message().Sender.Username, "encryptedMnemonic", encryptedMnemonic)
	redisClient.HSet(ctx, c.Message().Sender.Username, "salt", salt)
	_, _ = b.Send(c.Chat(), "Private key encrypted with password and stored safely in Redis!")
}

func HandleSendMnemonicStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	*step = "confirmMnemonic"
	_, _ = b.Send(c.Chat(), "Mnemonic sent! Now go to confirmation part!")
}

func HandleRequestUserEnterPassword(b *tele.Bot, c tele.Context, step *string) {
	*step = "addPassword"
	_, _ = b.Send(c.Chat(), "Please enter your password, at least 16 characters long")
}

// TODO: implement confirmation mnemonic
func HandleConfirmMnemonicStep(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	*step = "receiveMnemonicWords"
	randomIndexes = utils.GetRandomIndexesForTesting(len(utils.SplitMnemonic(mnemonic.String())))
	text := fmt.Sprintf("please enter these missing words, seperate by space \n %s", utils.GenerateMissedWordsMnemonicFromIndexes(mnemonic.String(), randomIndexes))
	_, _ = b.Send(c.Chat(), text)
}

func HandleReceiveMnemonicWords(b *tele.Bot, c tele.Context, botClient internal.BotClient, step *string) {
	providedWords := utils.SplitMnemonic(c.Text())
	if len(providedWords) != 3 {
		_, _ = b.Send(c.Chat(), "Please provide 3 words")
		*step = "confirmMnemonic"
		return
	}

	result, _ := utils.MnemonicChallenge(mnemonic.String(), randomIndexes, [3]string{providedWords[0], providedWords[1], providedWords[2]})
	if result {
		_, _ = b.Send(c.Chat(), "Mnemonic confirmed!")
	} else {
		_, _ = b.Send(c.Chat(), "Mnemonic incorrect! Please enter again.")
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

package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/awnumar/memguard"
	tele "gopkg.in/telebot.v3"
)

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
			fmt.Println(len("✅"))
			keyboard[5][i].Text = keyboard[5][i].Text[3:]
		}
	}
	time.Sleep(1 * time.Second)
	return keyboard
}

func ModifyAmountToTransferButton(keyboard [][]tele.InlineButton, amount, denom string) [][]tele.InlineButton {
	if denom != "" {
		keyboard[3][0].Text = "Transfer " + amount + " " + denom
		return keyboard
	}
	return keyboard
}

const (
	selectTokenToBuyButtonLabel = "Select Token to Buy"
	enterAmountToBuyButtonLabel = "Enter Amount to Buy"
	selectTokenToPayButtonLabel = "Select Token to Pay"
	setPriceButtonLabel         = "Set Price"
)

func formatTokenToBuyButtonLabel(denomIn string) string {
	if denomIn == "" {
		return "🪙 " + selectTokenToBuyButtonLabel
	}
	return fmt.Sprintf("🛒 Buy: %s", strings.ToUpper(denomIn))
}

func formatAmountToBuyButtonLabel(amount float64, denomIn string) string {
	if amount == 0 {
		return "💰 " + enterAmountToBuyButtonLabel
	}
	return fmt.Sprintf("💸 Buy Amount: %.2f %s", amount, strings.ToUpper(denomIn))
}

func formatTokenToPayButtonLabel(denomOut string) string {
	if denomOut == "" {
		return "💳 " + selectTokenToPayButtonLabel
	}
	return fmt.Sprintf("💸 Pay With: %s", strings.ToUpper(denomOut))
}

func formatPriceButtonLabel(price float64, denomOut, denomIn string) string {
	if price == 0 {
		return "💲 " + setPriceButtonLabel
	}
	return fmt.Sprintf("💰 Price: %.2f %s per %s", price, strings.ToUpper(denomOut), strings.ToUpper(denomIn))
}

func ModifyLimitOrderMenu(keyboard [][]tele.InlineButton, orderInfo *types.LimitOrderInfo) [][]tele.InlineButton {
	keyboard[1][0].Text = formatTokenToBuyButtonLabel(orderInfo.DenomIn)
	keyboard[2][0].Text = formatAmountToBuyButtonLabel(orderInfo.Amount, orderInfo.DenomIn)
	keyboard[3][0].Text = formatTokenToPayButtonLabel(orderInfo.DenomOut)
	keyboard[4][0].Text = formatPriceButtonLabel(orderInfo.Price, orderInfo.DenomOut, orderInfo.DenomIn)
	return keyboard
}
func DeleteInputMessage(b *tele.Bot, c tele.Context) error {
	err := b.Delete(c.Message().ReplyTo)
	if err != nil {
		return err
	}
	return c.Delete()
}

// RetrievePrivateKeyFromRedis retrieves the private key from Redis and returns it as a LockedBuffer.
func RetrievePrivateKeyFromRedis(redisClient RedisClient, username string, password *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	// retrieve mnemonic
	ctx := context.Background()
	encryptedMnemonic, err := redisClient.HGet(ctx, username, "encryptedMnemonic").Result()
	if err != nil {
		return nil, err
	}
	salt, err := redisClient.HGet(ctx, username, "salt").Result()
	if err != nil {
		return nil, err
	}

	// Decrypt mnemonic
	key, err := utils.DeriveKeyFromSalt(password.String(), []byte(salt))
	if err != nil {
		return nil, fmt.Errorf("failed to derive key from salt: %w", err)
	}
	password.Destroy()

	decryptedMnemonic, err := utils.GetDecryptedMnemonic(key, encryptedMnemonic)
	if err != nil {
		memguard.WipeBytes(key)
		return nil, err
	}
	memguard.WipeBytes(key)

	// Create a LockedBuffer for the decrypted mnemonic
	mnemonicBuffer := memguard.NewBufferFromBytes([]byte(decryptedMnemonic))
	defer mnemonicBuffer.Destroy()

	// Derive the private key bytes from the mnemonic
	return utils.DerivePrivateKeyBufferFromMnemonic(mnemonicBuffer)
}

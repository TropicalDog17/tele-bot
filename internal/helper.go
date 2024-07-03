package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
			keyboard[5][i].Text = keyboard[5][i].Text[3:]
		}
	}
	return keyboard
}

func ModifyAmountToTransferButton(localizer *i18n.Localizer, keyboard [][]tele.InlineButton, amount, denom string) [][]tele.InlineButton {
	if denom != "" {
		keyboard[3][0].Text = localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "Transfer", Other: "Transfer"}}) + " " + amount + " " + denom
		return keyboard
	}
	return keyboard
}

func ModifyCustomTokenButton(keyboard [][]tele.InlineButton, denom string) [][]tele.InlineButton {
	keyboard = RemoveGreenTickToken(keyboard)
	if denom == "ATOM" {
		keyboard[2][0] = AddGreenTick(*types.BtnInlineAtom(&i18n.Localizer{}).Inline())
		return keyboard
	} else if denom == "INJ" {
		keyboard[2][1] = AddGreenTick(*types.BtnInlineInj(&i18n.Localizer{}).Inline())
	} else {
		if denom != "" {
			keyboard[2][2].Text = "Custom token: " + denom
			return keyboard
		}
	}
	return keyboard
}

func formatTokenToBuyButtonLabel(localizer *i18n.Localizer, denomIn string) string {
	if denomIn == "" {
		return "🪙 " + localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "SelectTokenToBuyButtonLabel", Other: "Select token to buy"},
		})
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "BuyTokenFormat", Other: "🛒 Buy: {{.Token}}"},
		TemplateData:   map[string]string{"Token": strings.ToUpper(denomIn)},
	})
}

func formatAmountToBuyButtonLabel(localizer *i18n.Localizer, amount float64, denomIn string) string {
	if amount == 0 {
		return "💰 " + localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "EnterAmountToBuyButtonLabel", Other: "Enter amount to buy"},
		})
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "BuyAmountFormat", Other: "💸 Buy Amount: {{.Amount}} {{.Token}}"},
		TemplateData: map[string]interface{}{
			"Amount": fmt.Sprintf("%.2f", amount),
			"Token":  strings.ToUpper(denomIn),
		},
	})
}

func formatTokenToPayButtonLabel(localizer *i18n.Localizer, denomOut string) string {
	if denomOut == "" {
		return "💳 " + localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "SelectTokenToPayButtonLabel", Other: "Select token to pay"},
		})
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "PayWithTokenFormat", Other: "💸 Pay With: {{.Token}}"},
		TemplateData:   map[string]string{"Token": strings.ToUpper(denomOut)},
	})
}

func formatPriceButtonLabel(localizer *i18n.Localizer, price float64, denomOut, denomIn string) string {
	if price == 0 {
		return "💲 " + localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ID: "SetPriceButtonLabel", Other: "Set price"},
		})
	}
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{ID: "PriceFormat", Other: "💰 Price: {{.Price}} {{.DenomOut}} per {{.DenomIn}}"},
		TemplateData: map[string]interface{}{
			"Price":    fmt.Sprintf("%.2f", price),
			"DenomOut": strings.ToUpper(denomOut),
			"DenomIn":  strings.ToUpper(denomIn),
		},
	})
}
func ModifyLimitOrderMenu(keyboard [][]tele.InlineButton, orderInfo *types.LimitOrderInfo, localizer *i18n.Localizer) [][]tele.InlineButton {
	keyboard[1][0].Text = formatTokenToBuyButtonLabel(localizer, orderInfo.DenomIn)
	keyboard[2][0].Text = formatAmountToBuyButtonLabel(localizer, orderInfo.Amount, orderInfo.DenomIn)
	keyboard[3][0].Text = formatTokenToPayButtonLabel(localizer, orderInfo.DenomOut)
	keyboard[4][0].Text = formatPriceButtonLabel(localizer, orderInfo.Price, orderInfo.DenomOut, orderInfo.DenomIn)
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
func RetrievePrivateKeyFromRedis(redisClient RedisClient, username string, password string) (string, error) {
	// retrieve mnemonic
	ctx := context.Background()
	encryptedMnemonic, err := redisClient.HGet(ctx, username, "encryptedMnemonic").Result()
	if err != nil {
		return "", err
	}
	salt, err := redisClient.HGet(ctx, username, "salt").Result()
	if err != nil {
		return "", err
	}

	// Decrypt mnemonic
	key, err := utils.DeriveKeyFromSalt(password, []byte(salt))
	if err != nil {
		return "", fmt.Errorf("failed to derive key from salt: %w", err)
	}

	decryptedMnemonic, err := utils.GetDecryptedMnemonic(key, encryptedMnemonic)
	if err != nil {
		return "", err
	}
	// Derive the private key bytes from the mnemonic
	return utils.DerivePrivateKeyBufferFromMnemonic(decryptedMnemonic)
}

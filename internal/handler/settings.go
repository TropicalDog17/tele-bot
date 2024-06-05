package handler

import (
	"context"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/TropicalDog17/tele-bot/internal"
	types "github.com/TropicalDog17/tele-bot/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	tele "gopkg.in/telebot.v3"
)

func HandleSettings(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, clients map[string]internal.BotClient, menuSetting *tele.ReplyMarkup, btnSetting, btnChangeLanguage tele.Btn, currentStep *string) {
	authRoute.Handle(&btnSetting, func(c tele.Context) error {
		_, ok := clients[c.Message().Sender.Username]
		if !ok {
			return c.Send("Client not found", types.Menu)
		}
		return c.Send("Settings", types.MenuViewSettings)
	})
	authRoute.Handle(&btnChangeLanguage, func(c tele.Context) error {
		return HandleChangeLanguage(b, c, clients[c.Sender().Username], currentStep)
	})
}

func HandleSettingsStep(b *tele.Bot, localizer *i18n.Localizer, c tele.Context, client internal.BotClient, currentStep *string) error {
	switch *currentStep {
	case "changeLanguage":
		return HandleChangeLanguage(b, c, client, currentStep)
	case "changeCurrency":
		return HandleChangeCurrency(b, c, client, currentStep)
	case "changePassword":
		return HandleChangePassword(b, c, client, currentStep)
	case "deletePassword":
		return HandleDeletePassword(b, c, client, currentStep)
	case "userInputLanguage":
		return HandleInputChangeLanguage(b, &localizer, c, client, currentStep)
	}

	return nil
}

func HandleChangeLanguage(b *tele.Bot, c tele.Context, client internal.BotClient, currentStep *string) error {
	*currentStep = "userInputLanguage"
	return c.Reply("Please input your language")
}
func HandleInputChangeLanguage(b *tele.Bot, localizer **i18n.Localizer, c tele.Context, client internal.BotClient, currentStep *string) error {
	*currentStep = ""
	lang := strings.ToLower(c.Text())
	if lang == "english" || lang == "en" {
		client.GetRedisInstance().HSet(context.Background(), c.Sender().Username, "language", "en")
		bundle := i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.MustLoadMessageFile("active.en.toml")
		*localizer = i18n.NewLocalizer(bundle, "en-US")
		return c.Reply("Language changed to English. Type /menu to view menu")
	} else if lang == "vietnamese" || lang == "vi" || lang == "vn" {
		bundle := i18n.NewBundle(language.Vietnamese)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.MustLoadMessageFile("active.vi.toml")
		client.GetRedisInstance().HSet(context.Background(), c.Sender().Username, "language", "vi")
		*localizer = i18n.NewLocalizer(bundle, "vi")
		return c.Reply("Ngôn ngữ chuyển sang Tiếng Việt thành công. Gõ /menu để xem menu")
	} else {
		return c.Reply("Invalid language")
	}
}

func HandleChangeCurrency(b *tele.Bot, c tele.Context, client internal.BotClient, currentStep *string) error {
	*currentStep = "changeCurrency"
	return c.Reply("Please select your currency")
}

func HandleChangePassword(b *tele.Bot, c tele.Context, client internal.BotClient, currentStep *string) error {
	*currentStep = "changePassword"
	return c.Reply("Please enter your new password")
}

func HandleDeletePassword(b *tele.Bot, c tele.Context, client internal.BotClient, currentStep *string) error {
	*currentStep = "deletePassword"
	return c.Reply("Please enter your password to delete it")
}

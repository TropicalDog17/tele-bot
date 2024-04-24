package internal

import (
	"fmt"
	"time"

	"github.com/TropicalDog17/tele-bot/internal/types"
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

func ModifyLimitOrderMenu(keyboard [][]tele.InlineButton, orderInfo *types.LimitOrderInfo) [][]tele.InlineButton {
	if orderInfo.DenomOut != "" {
		keyboard[1][0].Text = fmt.Sprintf("Token to Pay: %s", orderInfo.DenomOut)
	}
	if orderInfo.Amount != 0 {
		keyboard[2][0].Text = fmt.Sprintf("Amount: %f %s", orderInfo.Amount, orderInfo.DenomIn)
	}
	if orderInfo.Price != 0 {
		keyboard[3][0].Text = fmt.Sprintf("Price: %f", orderInfo.Price)
	}
	return keyboard
}

func DeleteInputMessage(b *tele.Bot, c tele.Context) error {
	err := b.Delete(c.Message().ReplyTo)
	if err != nil {
		return err
	}
	return c.Delete()
}

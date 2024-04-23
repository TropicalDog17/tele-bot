package internal

import (
	"fmt"
	"time"

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

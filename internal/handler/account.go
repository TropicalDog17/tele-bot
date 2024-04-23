package handler

import (
	"fmt"
	"os"

	"github.com/TropicalDog17/tele-bot/internal"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	tele "gopkg.in/telebot.v3"
)

func HandleAddressQr(b *tele.Bot, client *internal.Client) {
	b.Handle(&tele.Btn{Unique: "qr"}, func(c tele.Context) error {
		// Get the address
		address := client.GetAddress()
		// Address
		qrc, err := qrcode.New(address)
		if err != nil {
			fmt.Printf("could not generate QRCode: %v", err)
			return err
		}

		w, err := standard.New("../temp/repo-qrcode.jpeg")
		if err != nil {
			fmt.Printf("standard.New failed: %v", err)
			return err
		}
		defer w.Close()

		// save file
		if err = qrc.Save(w); err != nil {
			fmt.Printf("could not save image: %v", err)
		}
		photo := &tele.Photo{File: tele.FromDisk("../temp/repo-qrcode.jpeg")}
		_, err = b.Send(c.Chat(), photo)
		if err != nil {
			return err
		}
		// Remove the file after sending
		return os.Remove("../temp/repo-qrcode.jpeg")
	})
}

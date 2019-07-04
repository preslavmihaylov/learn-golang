package emails

import (
	"context"
	"time"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/config"
	"gopkg.in/mailgun/mailgun-go.v3"
)

const (
	lenslockedSupportEmail = "Lenslocked Support <support@lenslocked.myslidekit.com>"
)

var mailgunClient mailgun.Mailgun

func Setup(cfg config.MailgunConfig) {
	mailgunClient = mailgun.NewMailgun(cfg.Domain, cfg.APIKey)

	// setup rest api base to be the eu mailgun domain
	mailgunClient.SetAPIBase("https://api.eu.mailgun.net/v3")
}

func Send(recipient string, subject, msg string) error {
	m := mailgunClient.NewMessage(lenslockedSupportEmail, subject, msg, recipient)
	m.SetHtml(msg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mailgunClient.Send(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

package emails

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/preslavmihaylov/learn-golang/go-webdev/lenslocked.com/config"
	"gopkg.in/mailgun/mailgun-go.v3"
)

const (
	lenslockedSupportEmail    = "Lenslocked Support <support@lenslocked.myslidekit.com>"
	resetPasswordSubject      = "Instructions for resetting your password"
	resetPasswordTextTemplate = `Hi There!

It appears that you have requested a password reset. If this was you, please follow the link below:
%s

If you are asked for a token, please use the following value:
%s

If you didn't request a password reset you can safely ignore this email.

Best,
Lenslocked Support
	`
	resetPasswordHTMLTemplate = `Hi There!<br/>
<br/>
It appears that you have requested a password reset. If this was you, please follow the link below:<br/>
<a href="%s">%s</a><br/>
<br/>
If you are asked for a token, please use the following value:<br/>
%s<br/>
<br/>
If you didn't request a password reset you can safely ignore this email.<br/>
<br/>
Best,<br/>
Lenslocked Support<br/>
	`
	resetPasswordBaseURL = "https://lenslocked.myslidekit.com/reset_password"
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

func SendResetPasswordEmail(recipient, token string) error {
	urlValues := url.Values{}
	urlValues.Set("token", token)

	url := resetPasswordBaseURL + "?" + urlValues.Encode()
	txt := fmt.Sprintf(resetPasswordTextTemplate, url, token)
	html := fmt.Sprintf(resetPasswordHTMLTemplate, url, url, token)

	m := mailgunClient.NewMessage(
		lenslockedSupportEmail, resetPasswordSubject, txt, recipient)
	m.SetHtml(html)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, _, err := mailgunClient.Send(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

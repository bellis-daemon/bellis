package email

import (
	"fmt"
	"github.com/bellis-daemon/bellis/common/cache"
	mail "github.com/xhit/go-simple-mail/v2"
)

func SendCaptcha(email string) error {
	captcha, err := cache.CaptchaSet(email)
	if err != nil {
		return fmt.Errorf("cant set captcha for user: %w", err)
	}
	html, err := base().GenerateHTML(captchaEmail(captcha))
	if err != nil {
		return fmt.Errorf("cant generate email html: %w", err)
	}
	cl, err := tencentSmtpClient()
	if err != nil {
		return fmt.Errorf("cant connect to smtp: %w", err)
	}
	err = mail.NewMSG().
		SetFrom("Bellis Envoy <envoy@bellis.minoic.top>").
		SetReplyTo("minoic2020@gmail.com").
		AddTo(email).
		SetSubject("Bellis register captcha").
		SetBody(mail.TextHTML, html).
		Send(cl)
	if err != nil {
		return fmt.Errorf("cant send email via smtp: %w", err)
	}
	return nil
}

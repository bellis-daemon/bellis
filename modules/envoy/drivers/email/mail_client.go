package email

import (
	"time"

	"github.com/bellis-daemon/bellis/common/storage"
	mail "github.com/xhit/go-simple-mail/v2"
)

func tencentSmtpClient() (*mail.SMTPClient, error) {
	sv := &mail.SMTPServer{
		Authentication: mail.AuthPlain,
		Encryption:     mail.EncryptionTLS,
		Username:       "envoy@bellis.minoic.top",
		Password:       storage.Secret("tencent_smtp_password"),
		ConnectTimeout: 3 * time.Second,
		SendTimeout:    3 * time.Second,
		Host:           "gz-smtp.qcloudmail.com",
		Port:           465,
		KeepAlive:      false,
	}
	cl, err := sv.Connect()
	return cl, err
}

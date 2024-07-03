package email

import (
	"time"

	"github.com/bellis-daemon/bellis/common/storage"
	mail "github.com/xhit/go-simple-mail/v2"
)

func tencentSmtpClient() (*mail.SMTPClient, error) {
	sv := &mail.SMTPServer{
		Authentication: mail.AuthAuto,
		Encryption:     mail.EncryptionSSL,
		Username:       storage.Config().SMTPUsername,
		Password:       storage.Config().SMTPPassword,
		ConnectTimeout: 10 * time.Second,
		SendTimeout:    10 * time.Second,
		Host:           storage.Config().SMTPHostname,
		Port:           storage.Config().SMTPPort,
		KeepAlive:      false,
	}
	cl, err := sv.Connect()
	return cl, err
}

package email

import (
	"fmt"
	"testing"

	_ "github.com/bellis-daemon/bellis/common/storage"
	mail "github.com/xhit/go-simple-mail/v2"
)

func Test_tencentSmtpClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *mail.SMTPClient
		wantErr bool
	}{
		{
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tencentSmtpClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("tencentSmtpClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			msg := mail.NewMSG().
				SetFrom("Bellis Envoy <envoy@bellis.minoic.top>").
				AddTo("minoic2020@gmail.com").
				SetSubject("Bellis entity offline alert").
				SetBody(mail.TextPlain, "test email")
			fmt.Println(msg.GetFrom(),msg.GetMessage())
			err = msg.Send(got)				
			if (err != nil) != tt.wantErr {
				t.Errorf("tencentSmtpClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

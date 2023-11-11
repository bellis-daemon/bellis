package consumer

import (
	"context"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers/email"
	"github.com/spf13/cast"
)

func emailCaptcha() {
	stream.Register(common.CaptchaToEmail, func(ctx context.Context, message *redistream.Message) error {
		return email.SendCaptcha(cast.ToString(message.Values["Email"]))
	})
}

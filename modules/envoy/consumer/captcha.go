package consumer

import (
	"context"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/minoic/glgf"
)

func emailCaptcha() {
	stream.Register(common.CaptchaToEmail, func(ctx context.Context, message *redistream.Message) error {
		glgf.Debug(message)
		return nil
	})
}

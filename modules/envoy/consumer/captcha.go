package consumer

import (
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/minoic/glgf"
)

func emailCaptcha() {
	redistream.Instance().Register("CaptchaToEmail", func(message *redistream.Message) error {
		glgf.Debug(message)
		return nil
	})
}

package producer

import (
	"context"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
)

func EnvoyCaptchaToEmail(ctx context.Context, receiver string, captcha string) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: "CaptchaToEmail",
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"Receiver": receiver,
			"Captcha":  captcha,
		},
	}).Err()
}

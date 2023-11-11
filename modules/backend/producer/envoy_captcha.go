package producer

import (
	"context"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
)

func EnvoyCaptchaToEmail(ctx context.Context, email string) error {
	return storage.Redis().XAdd(ctx, &redis.XAddArgs{
		Stream: common.CaptchaToEmail,
		MaxLen: 256,
		Approx: true,
		Values: map[string]interface{}{
			"Email": email,
		},
	}).Err()
}

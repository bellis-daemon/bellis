package cache

import (
	"context"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

const captchaPrefix = "CAPTCHA_"
const captchaLength = 6

func CaptchaSet(email string) (string, error) {
	captcha := cryptoo.RandNum(captchaLength)
	err := storage.Redis().Set(context.Background(), captchaPrefix+email, captcha, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return captcha, nil
}

func CaptchaCheck(email string, captcha string) (bool, error) {
	c, err := storage.Redis().Get(context.Background(), captchaPrefix+email).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return c == captcha, nil
}

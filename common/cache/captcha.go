package cache

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/redis/go-redis/v9"
	"time"
)

const captchaPrefix = "CAPTCHA_"
const captchaLength = 6

func CaptchaSet(key string) (string, error) {
	captcha := cryptoo.RandNum(captchaLength)
	err := storage.Redis().Set(context.Background(), captchaPrefix+key, captcha, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return captcha, nil
}

func CaptchaCheck(key string, captcha string) (bool, error) {
	c, err := storage.Redis().Get(context.Background(), captchaPrefix+key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return c == captcha, nil
}

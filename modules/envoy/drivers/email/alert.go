package email

import (
	"github.com/bellis-daemon/bellis/common/models"
	"time"
)

func AlertOffline(entity *models.Application, policy *models.EnvoyEmail, msg string, offlineTime time.Time) error {
	panic("not implemented")
}

package drivers

import (
	"github.com/bellis-daemon/bellis/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EnvoyDriver interface {
	AlertOffline(entity *models.Application, msg string, offlineTime time.Time) error
	WithPolicy(policy any) EnvoyDriver
	WithPolicyId(policyId primitive.ObjectID) EnvoyDriver
}

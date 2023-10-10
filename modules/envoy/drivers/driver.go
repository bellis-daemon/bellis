package drivers

import (
	"github.com/bellis-daemon/bellis/common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnvoyDriver interface {
	AlertOffline(entity *models.Application, log *models.OfflineLog) error
	WithPolicy(policy any) EnvoyDriver
	WithPolicyId(policyId primitive.ObjectID) EnvoyDriver
}

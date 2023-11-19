package ntfy

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	policy *models.EnvoyNtfy
	ctx    context.Context
}

func (this *handler) AlertOffline(user *models.User, entity *models.Application, log *models.OfflineLog) error {
	//TODO implement me
	panic("implement me")
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	//TODO implement me
	panic("implement me")
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	//TODO implement me
	panic("implement me")
}

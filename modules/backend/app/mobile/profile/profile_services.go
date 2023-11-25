package profile

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func useNewPolicy(ctx context.Context, policy any) (*EnvoyPolicy, error) {
	user := midwares.GetUserFromCtx(ctx)
	var targetPolicyType models.EnvoyPolicyType
	var targetPolicy any
	switch policy.(type) {
	case *Gotify:
		targetPolicy = &models.EnvoyGotify{
			ID:    primitive.NewObjectID(),
			URL:   policy.(*Gotify).Url,
			Token: policy.(*Gotify).Token,
		}
		targetPolicyType = models.IsEnvoyGotify
	case *Email:
		targetPolicy = &models.EnvoyEmail{
			ID:      primitive.NewObjectID(),
			Address: policy.(*Email).Address,
		}
		targetPolicyType = models.IsEnvoyEmail
	case *Webhook:
		targetPolicy = &models.EnvoyWebhook{
			ID:       primitive.NewObjectID(),
			URL:      policy.(*Webhook).Url,
			Insecure: policy.(*Webhook).Insecure,
		}
		targetPolicyType = models.IsEnvoyWebhook
	default:
		return nil, errors.New("invalid policy type")
	}
	err := user.SetProfile(ctx, targetPolicyType, targetPolicy)
	if err != nil {
		return &EnvoyPolicy{}, status.Error(codes.Internal, err.Error())
	}
	return &EnvoyPolicy{
		PolicyID:     user.Envoy.PolicyID.Hex(),
		PolicyType:   int32(user.Envoy.PolicyType),
		OfflineAlert: user.Envoy.OfflineAlert,
		PredictAlert: user.Envoy.PredictAlert,
	}, nil
}

package profile

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func useNewPolicy(ctx context.Context, policy any) (*EnvoyPolicy, error) {
	user := midwares.GetUserFromCtx(ctx)
	err := storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// delete old policy
		if user.Envoy.PolicyID != primitive.NilObjectID {
			var coll *mongo.Collection
			switch user.Envoy.PolicyType {
			case models.IsEnvoyGotify:
				coll = storage.CEnvoyGotify
			case models.IsEnvoyEmail:
				coll = storage.CEnvoyEmail
			}
			_, err := coll.DeleteOne(sessionContext, bson.M{
				"_id": user.Envoy.PolicyID,
			})
			if err != nil {
				return err
			}
		}
		// create new policy
		switch policy.(type) {
		case *Gotify:
			one, err := storage.CEnvoyGotify.InsertOne(sessionContext, &models.EnvoyGotify{
				ID:    primitive.NewObjectID(),
				URL:   policy.(*Gotify).Url,
				Token: policy.(*Gotify).Token,
			})
			if err != nil {
				return err
			}
			user.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
			user.Envoy.PolicyType = models.IsEnvoyGotify
		case *Email:
			one, err := storage.CEnvoyEmail.InsertOne(sessionContext, &models.EnvoyEmail{
				ID:      primitive.NewObjectID(),
				Address: policy.(*Email).Address,
			})
			if err != nil {
				return err
			}
			user.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
			user.Envoy.PolicyType = models.IsEnvoyEmail
		case *Webhook:
			one, err := storage.CEnvoyEmail.InsertOne(sessionContext, &models.EnvoyWebhook{
				ID:       primitive.NewObjectID(),
				URL:      policy.(*Webhook).Url,
				Insecure: policy.(*Webhook).Insecure,
			})
			if err != nil {
				return err
			}
			user.Envoy.PolicyID = one.InsertedID.(primitive.ObjectID)
			user.Envoy.PolicyType = models.IsEnvoyEmail
		}
		// modify user model
		_, err := storage.CUser.ReplaceOne(sessionContext, bson.M{"_id": user.ID}, user)
		if err != nil {
			return err
		}
		return nil
	})
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

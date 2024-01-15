package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/generic"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sort"
	"sync"
	"time"
)

// implement ProfileServiceServer
type handler struct{}

func (h handler) CreateEnvoyPolicy(ctx context.Context, policy *EnvoyPolicy) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	if !user.UsageEnvoyAccessible() {
		return &emptypb.Empty{}, errors.New("the number of available envoy policies has reached the upper limit")
	}
	var targetPolicyType models.EnvoyPolicyType
	var targetPolicy any
	header := models.EnvoyHeader{
		UserID:       user.ID,
		CreatedAt:    time.Now(),
		OfflineAlert: policy.OfflineAlert,
		Sensitive:    int(policy.Sensitive),
	}
	switch models.EnvoyPolicyType(policy.PolicyType) {
	case models.IsEnvoyGotify:
		targetPolicy = &models.EnvoyGotify{
			EnvoyHeader: header,
			ID:          primitive.NewObjectID(),
			URL:         policy.PolicyContent.GetGotify().Url,
			Token:       policy.PolicyContent.GetGotify().Token,
		}
		targetPolicyType = models.IsEnvoyGotify
	case models.IsEnvoyEmail:
		targetPolicy = &models.EnvoyEmail{
			EnvoyHeader: header,
			ID:          primitive.NewObjectID(),
			Address:     policy.PolicyContent.GetEmail().Address,
		}
		targetPolicyType = models.IsEnvoyEmail
	case models.IsEnvoyWebhook:
		targetPolicy = &models.EnvoyWebhook{
			EnvoyHeader: header,
			ID:          primitive.NewObjectID(),
			URL:         policy.PolicyContent.GetWebhook().Url,
			Insecure:    policy.PolicyContent.GetWebhook().Insecure,
		}
		targetPolicyType = models.IsEnvoyWebhook
	case models.IsEnvoyTelegram:
		targetPolicy = &models.EnvoyTelegram{
			EnvoyHeader: header,
			ID:          primitive.NewObjectID(),
			ChatID:      policy.PolicyContent.GetTelegram().ChatId,
		}
	default:
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "invalid policy type")
	}
	err := storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		// create new policy
		var policyId primitive.ObjectID
		inserted, err := targetPolicyType.GetCollection().InsertOne(ctx, targetPolicy)
		if err != nil {
			return err
		}
		policyId = inserted.InsertedID.(primitive.ObjectID)
		// modify user model
		_, err = storage.CUser.UpdateOne(sessionContext, bson.M{"_id": user.ID}, bson.M{"$push": bson.M{"EnvoyPolicies": bson.M{
			"PolicyID":   policyId,
			"PolicyType": targetPolicyType,
		}}})
		if err != nil {
			return err
		}
		err = user.UsageEnvoyIncr(sessionContext, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) UpdateEnvoyPolicy(ctx context.Context, policy *EnvoyPolicy) (*emptypb.Empty, error) {
	id, err := primitive.ObjectIDFromHex(policy.PolicyID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "cant find specified policy")
	}
	user := midwares.GetUserFromCtx(ctx)
	err = storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		policyType := models.EnvoyPolicyType(policy.PolicyType)
		one, err := storage.CUser.UpdateOne(ctx, bson.M{"_id": user.ID, "EnvoyPolicies.PolicyID": id}, bson.M{"$set": bson.M{
			"EnvoyPolicies.$.PolicyType": policyType,
		}})
		if err != nil {
			return err
		}
		if one.ModifiedCount == 0 {
			return errors.New("cant find specified envoy policy")
		}
		set := bson.M{
			"Sensitive":    policy.Sensitive,
			"OfflineAlert": policy.OfflineAlert,
		}
		switch policyType {
		case models.IsEnvoyGotify:
			set["URL"] = policy.PolicyContent.GetGotify().Url
			set["Token"] = policy.PolicyContent.GetGotify().Token
		case models.IsEnvoyEmail:
			set["Address"] = policy.PolicyContent.GetEmail().Address
		case models.IsEnvoyTelegram:
			set["ChatID"] = policy.PolicyContent.GetTelegram().ChatId
		case models.IsEnvoyWebhook:
			set["URL"] = policy.PolicyContent.GetWebhook().Url
			set["Insecure"] = policy.PolicyContent.GetWebhook().Insecure
		}
		_, err = policyType.GetCollection().UpdateByID(ctx, id, bson.M{"$set": set})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (h handler) DeleteEnvoyPolicy(ctx context.Context, policy *EnvoyPolicy) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	id, err := primitive.ObjectIDFromHex(policy.PolicyID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "invalid policy id")
	}
	policyType := models.EnvoyPolicyType(policy.PolicyType)
	err = storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		updated, err := storage.CUser.UpdateByID(sessionContext, user.ID, bson.M{"$pull": bson.M{"EnvoyPolicies": bson.M{
			"PolicyID": id,
		}}})
		if err != nil {
			return err
		}
		if updated.ModifiedCount == 0 {
			return errors.New("policy does not exist in user envoy policies")
		}
		_, err = policyType.GetCollection().DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return err
		}
		err = user.UsageEnvoyIncr(sessionContext, -1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) GetUserLoginLogs(ctx context.Context, empty *emptypb.Empty) (*UserLoginLogs, error) {
	user := midwares.GetUserFromCtx(ctx)
	var logs []models.UserLoginLog
	find, err := storage.CUserLoginLog.Find(ctx, bson.M{"UserID": user.ID}, options.Find().SetSort(bson.M{"_id": -1}).SetLimit(7))
	if err != nil {
		return &UserLoginLogs{}, status.Error(codes.Internal, err.Error())
	}
	err = find.All(ctx, &logs)
	if err != nil {
		return &UserLoginLogs{}, status.Error(codes.Internal, err.Error())
	}
	return &UserLoginLogs{Logs: generic.SliceConvert(logs, func(s models.UserLoginLog) *UserLoginLog {
		return &UserLoginLog{
			LoginTime:  s.LoginTime.In(user.Timezone.Location()).Format(time.DateTime),
			Location:   s.Location,
			Device:     s.Device,
			DeticeType: s.DeviceType,
		}
	})}, nil
}

// GetEnvoyTelegramLink generates a unique Telegram link for the user to connect with the server's Telegram bot.
// It creates a unique captcha, associates it with the user ID in Redis, and returns the Telegram link for the user to initiate the connection.
func (h handler) GetEnvoyTelegramLink(ctx context.Context, empty *emptypb.Empty) (*EnvoyTelegramLink, error) {
	if storage.Config().TelegramBotName == "" {
		return &EnvoyTelegramLink{}, status.Error(codes.Internal, "telegram not supported on server")
	}
	user := midwares.GetUserFromCtx(ctx)
	captcha := "tg_" + cryptoo.RandString(24)
	err := storage.Redis().Set(ctx, captcha, user.ID.Hex(), time.Minute).Err()
	if err != nil {
		return &EnvoyTelegramLink{}, status.Error(codes.Internal, err.Error())
	}
	link := fmt.Sprintf("t.me/%s?start=%s", storage.Config().TelegramBotName, captcha)
	return &EnvoyTelegramLink{
		Url: link,
	}, nil
}

// ChangeSensitive updates the sensitivity level for the user's data in the storage.
// It modifies the sensitivity level for the user's data and returns an empty response if successful.
func (h handler) ChangeSensitive(ctx context.Context, sensitive *Sensitive) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	_, err := storage.CUser.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{
		"Envoy.Sensitive": sensitive.Level,
	}})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ChangePassword updates the user's password with the provided new password.
// It sets the new password for the user and returns an empty response or an error if the operation fails.
func (h handler) ChangePassword(ctx context.Context, password *NewPassword) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	err := user.SetPassword(ctx, password.Password)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// ChangeEmail updates the user's email with the provided new email address.
// It modifies the user's email in the storage and returns an empty response or an error if the operation fails.
func (h handler) ChangeEmail(ctx context.Context, email *NewEmail) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	_, err := storage.CUser.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"Email": email.Email}})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

// GetUserProfile retrieves the user's profile details including email, creation date, access level, and Envoy policy information.
// It fetches the user's policy content based on the policy type and returns the complete UserProfile object.
func (h handler) GetUserProfile(ctx context.Context, empty *emptypb.Empty) (*UserProfile, error) {
	user := midwares.GetUserFromCtx(ctx)
	ret := &UserProfile{
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Local().Format(time.DateTime),
		Level:     uint32(user.Level),
		Policies:  nil,
	}
	var wg sync.WaitGroup
	for i := range user.EnvoyPolicies {
		wg.Add(1)
		p := user.EnvoyPolicies[i]
		go func() {
			defer wg.Done()
			var content bson.M
			err := p.PolicyType.GetCollection().FindOne(ctx, bson.M{"_id": p.PolicyID}).Decode(&content)
			if err != nil {
				glgf.Error(err)
				return
			}
			var policyContent *EnvoyPolicyContent
			switch p.PolicyType {
			case models.IsEnvoyGotify:
				policyContent = &EnvoyPolicyContent{
					Content: &EnvoyPolicyContent_Gotify{
						Gotify: &Gotify{
							Url:   content["URL"].(string),
							Token: content["Token"].(string),
						},
					},
				}
			case models.IsEnvoyEmail:
				policyContent = &EnvoyPolicyContent{
					Content: &EnvoyPolicyContent_Email{
						Email: &Email{
							Address: content["Address"].(string),
						},
					},
				}
			case models.IsEnvoyWebhook:
				policyContent = &EnvoyPolicyContent{
					Content: &EnvoyPolicyContent_Webhook{
						Webhook: &Webhook{
							Url:      content["URL"].(string),
							Insecure: content["Insecure"].(bool),
						},
					},
				}
			case models.IsEnvoySMS:

			case models.IsEnvoyTelegram:
				policyContent = &EnvoyPolicyContent{
					Content: &EnvoyPolicyContent_Telegram{
						Telegram: &Telegram{
							ChatId: cast.ToInt64(content["ChatID"]),
						},
					},
				}
			}
			ret.Policies = append(ret.Policies, &EnvoyPolicy{
				PolicyID:      p.PolicyID.Hex(),
				PolicyType:    int32(p.PolicyType),
				Sensitive:     cast.ToInt32(content["Sensitive"]),
				OfflineAlert:  cast.ToBool(content["OfflineAlert"]),
				PolicyContent: policyContent,
			})
		}()
	}
	wg.Wait()
	sort.Slice(ret.Policies, func(i, j int) bool {
		return ret.Policies[i].PolicyID < ret.Policies[j].PolicyID
	})
	return ret, nil
}

func (h handler) NeedAuth() bool {
	return true
}

func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterProfileServiceServer(server, &handler{})
		return "Profile"
	})
}

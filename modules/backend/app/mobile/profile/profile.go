package profile

import (
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

// implement ProfileServiceServer
type handler struct{}

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

func (h handler) ChangePassword(ctx context.Context, password *NewPassword) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	err := user.SetPassword(ctx, password.Password)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) ChangeEmail(ctx context.Context, email *NewEmail) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	_, err := storage.CUser.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"Email": email.Email}})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) ChangeAlert(ctx context.Context, alert *Alert) (*emptypb.Empty, error) {
	user := midwares.GetUserFromCtx(ctx)
	_, err := storage.CUser.UpdateOne(ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"Envoy.OfflineAlert": alert.OfflineAlert,
			"Envoy.PredictAlert": alert.PredictAlert,
		}})
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) UseGotify(ctx context.Context, gotify *Gotify) (*EnvoyPolicy, error) {
	return useNewPolicy(ctx, gotify)
}

func (h handler) UseEmail(ctx context.Context, email *Email) (*EnvoyPolicy, error) {
	return useNewPolicy(ctx, email)
}

func (h handler) UseWebhook(ctx context.Context, webhook *Webhook) (*EnvoyPolicy, error) {
	return useNewPolicy(ctx, webhook)
}

func (h handler) GetUserProfile(ctx context.Context, empty *emptypb.Empty) (*UserProfile, error) {
	user := midwares.GetUserFromCtx(ctx)
	ret := &UserProfile{
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Local().Format(time.DateTime),
		Level:     uint32(user.Level),
		Envoy: &EnvoyPolicy{
			PolicyID:      user.Envoy.PolicyID.Hex(),
			PolicyType:    int32(user.Envoy.PolicyType),
			OfflineAlert:  user.Envoy.OfflineAlert,
			PredictAlert:  user.Envoy.PredictAlert,
			PolicyContent: &EnvoyPolicyContent{},
		},
	}
	switch user.Envoy.PolicyType {
	case models.IsEnvoyGotify:
		var policy models.EnvoyGotify
		err := storage.CEnvoyGotify.FindOne(ctx, bson.M{
			"_id": user.Envoy.PolicyID,
		}).Decode(&policy)
		if err != nil {
			glgf.Error(err)
			return ret, nil
		}
		ret.Envoy.PolicyContent.Content = &EnvoyPolicyContent_Gotify{
			Gotify: &Gotify{
				Url:   policy.URL,
				Token: policy.Token,
			},
		}
	case models.IsEnvoyEmail:
		var policy models.EnvoyEmail
		err := storage.CEnvoyEmail.FindOne(ctx, bson.M{
			"_id": user.Envoy.PolicyID,
		}).Decode(&policy)
		if err != nil {
			return ret, nil
		}
		ret.Envoy.PolicyContent.Content = &EnvoyPolicyContent_Email{
			Email: &Email{
				Address: policy.Address,
			},
		}
	case models.IsEnvoyWebhook:
	case models.IsEnvoySMS:
	case models.IsEnvoyTelegram:
		var policy models.EnvoyTelegram
		err := storage.CEnvoyTelegram.FindOne(ctx, bson.M{
			"_id": user.Envoy.PolicyType,
		}).Decode(&policy)
		if err != nil {
			return ret, nil
		}
		ret.Envoy.PolicyContent.Content = &EnvoyPolicyContent_Telegram{
			Telegram: &Telegram{
				ChatId: policy.ChatId,
			},
		}
	}
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

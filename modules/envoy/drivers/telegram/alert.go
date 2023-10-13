package telegram

import (
	"context"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	ctx    context.Context
	policy *models.EnvoyTelegram
}

func (this *handler) AlertOffline(entity *models.Application, log *models.OfflineLog) error {
	api, err := tgbotapi.NewBotAPI(storage.Firebase().ConfigGetString("telegram_bot_token"))
	if err != nil {
		return err
	}
	message := tgbotapi.NewMessage(this.policy.ChatId, log.OfflineMessage)
	_, err = api.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	this.policy = policy.(*models.EnvoyTelegram)
	return this
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	this.policy = new(models.EnvoyTelegram)
	err := storage.CEnvoyEmail.FindOne(this.ctx, bson.M{"_id": policyId}).Decode(this.policy)
	if err != nil {
		glgf.Error(err)
	}
	return this
}

func New(ctx context.Context) drivers.EnvoyDriver {
	return &handler{ctx: ctx}
}

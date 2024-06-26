package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	ctx    context.Context
	policy *models.EnvoyTelegram
}

func escapeCharacters(input string) string {
	charsToEscape := "_*[]()~`>#-|={}.!"
	escapedString := input

	for _, char := range charsToEscape {
		escapedString = strings.ReplaceAll(escapedString, string(char), "\\"+string(char))
	}

	return escapedString
}

func (this *handler) AlertOffline(user *models.User, entity *models.Application, log *models.OfflineLog) error {
	api, err := tgbotapi.NewBotAPIWithAPIEndpoint(storage.Config().TelegramBotToken, storage.Config().TelegramBotApiEndpoint+"/bot%s/%s")
	if err != nil {
		return err
	}
	message := tgbotapi.NewMessage(
		this.policy.ChatID,
		fmt.Sprintf("Bellis entity *OFFLINE* alert ⚠\n"+
			"The following is the information from this offline session:\n"+
			"*Entity name:* `%s`\n"+
			"*TimeZone:* `%s`\n"+
			"*Created time:* `%s`\n"+
			"*Offline time:* `%s`\n"+
			"*Offline message:* `%s`\n",
			entity.Name,
			user.Timezone,
			entity.CreatedAt.In(user.Timezone.Location()).Format(time.DateTime),
			log.OfflineTime.In(user.Timezone.Location()).Format(time.DateTime),
			log.OfflineMessage,
		),
	)
	message.ParseMode = tgbotapi.ModeMarkdown
	_, err = api.Send(message)
	if err != nil {
		return fmt.Errorf("%w via policy: %+v", err, this.policy)
	}
	return nil
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	this.policy = policy.(*models.EnvoyTelegram)
	return this
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	this.policy = new(models.EnvoyTelegram)
	err := storage.CEnvoyTelegram.FindOne(this.ctx, bson.M{"_id": policyId}).Decode(this.policy)
	if err != nil {
		glgf.Error(err)
	}
	return this
}

func (this *handler) PolicySnapShot() bson.M {
	ret := make(bson.M)
	_ = mapstructure.Decode(this.policy, &ret)
	return ret
}

func New(ctx context.Context) drivers.EnvoyDriver {
	return &handler{ctx: ctx}
}

package gotify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	gmodels "github.com/gotify/go-api-client/v2/models"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	policy *models.EnvoyGotify
	ctx    context.Context
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	this.policy = policy.(*models.EnvoyGotify)
	return this
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	this.policy = new(models.EnvoyGotify)
	err := storage.CEnvoyGotify.FindOne(this.ctx, bson.M{"_id": policyId}).Decode(this.policy)
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


func (this *handler) AlertOffline(user *models.User, entity *models.Application, log *models.OfflineLog) error {
	gotifyURL, err := url.Parse(this.policy.URL)
	if err != nil {
		return err
	}
	client := gotify.NewClient(gotifyURL, &http.Client{})
	params := message.NewCreateMessageParams()
	params.Body = &gmodels.MessageExternal{
		Title: "Offline alert - " + entity.Name,
		Message: fmt.Sprintf(
			"Your application <%s> just went offline at %s (%s), error message: %s",
			entity.Name,
			log.OfflineTime.In(user.Timezone.Location()).Format(time.DateTime),
			user.Timezone.Location().String(),
			log.OfflineMessage),
		Priority: 7,
		Date:     time.Now(),
	}
	_, err = client.Message.CreateMessage(params, auth.TokenAuth(this.policy.Token))
	if err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context) drivers.EnvoyDriver {
	return &handler{ctx: ctx}
}

package webhook

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	"github.com/minoic/glgf"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	policy *models.EnvoyWebhook
	ctx    context.Context
}

func (this *handler) AlertOffline(user *models.User, entity *models.Application, log *models.OfflineLog) error {
	parsedUrl, err := url.Parse(this.policy.URL)
	if err != nil {
		return err
	}
	if !this.policy.Insecure {
		if parsedUrl.Scheme != "https" {
			return errors.New("cant send alert to none https server without insecure option ")
		}
		dial, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", parsedUrl.Hostname(), parsedUrl.Port()), nil)
		if err != nil {
			return err
		}
		err = dial.VerifyHostname(parsedUrl.Hostname())
		if err != nil {
			return err
		}
		expire := dial.ConnectionState().PeerCertificates[0].NotAfter
		if expire.Before(time.Now()) {
			return errors.New("cant send alert to server with ssl certification without insecure option ")
		}
	}
	body := map[string]any{
		"EntityId":          entity.ID.Hex(),
		"EntityName":        entity.Name,
		"EntityDescription": entity.Description,
		"EntityCreatedAt":   entity.CreatedAt,
		"OfflineTime":       log.OfflineTime.Format(time.RFC3339),
		"OfflineMessage":    log.OfflineMessage,
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return err
	}
	_, err = http.Post(parsedUrl.String(), "application/json", &buf)
	if err != nil {
		return err
	}
	return nil
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	this.policy = policy.(*models.EnvoyWebhook)
	return this
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	this.policy = new(models.EnvoyWebhook)
	err := storage.CEnvoyWebhook.FindOne(this.ctx, bson.M{"_id": policyId}).Decode(this.policy)
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

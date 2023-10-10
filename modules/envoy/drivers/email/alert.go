package email

import (
	"context"
	"fmt"

	mail "github.com/xhit/go-simple-mail/v2"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/envoy/drivers"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type handler struct {
	policy *models.EnvoyEmail
	ctx    context.Context
}

func (this *handler) AlertOffline(entity *models.Application, log *models.OfflineLog) error {
	cl, err := tencentSmtpClient()
	if err != nil {
		return fmt.Errorf("cant connect to smtp server: %w", err)
	}
	var user models.User
	err = storage.CUser.FindOne(this.ctx, bson.M{"_id": entity.UserID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("cant find user in database: %w", err)
	}
	html, err := base().GenerateHTML(offlineEmail(&user, entity, log))
	if err != nil {
		return fmt.Errorf("cant generate email html: %w", err)
	}
	err = mail.NewMSG().
		SetFrom("bellis@email.minoic.top").
		AddTo(user.Email).
		SetSubject("Bellis entity offline alert").
		SetBody(mail.TextHTML, html).
		Send(cl)
	if err != nil {
		return fmt.Errorf("cant send email via smtp: %w", err)
	}
	return nil
}

func (this *handler) WithPolicy(policy any) drivers.EnvoyDriver {
	this.policy = policy.(*models.EnvoyEmail)
	return this
}

func (this *handler) WithPolicyId(policyId primitive.ObjectID) drivers.EnvoyDriver {
	this.policy = new(models.EnvoyEmail)
	err := storage.CEnvoyEmail.FindOne(this.ctx, bson.M{"_id": policyId}).Decode(this.policy)
	if err != nil {
		glgf.Error(err)
	}
	return this
}

func New(ctx context.Context) drivers.EnvoyDriver {
	return &handler{ctx: ctx}
}

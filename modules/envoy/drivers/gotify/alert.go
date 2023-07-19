package gotify

import (
	"fmt"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	gmodels "github.com/gotify/go-api-client/v2/models"
	"net/http"
	"net/url"
	"time"
)

func AlertOffline(entity *models.Application, policy *models.EnvoyGotify, msg string, offlineTime time.Time) error {
	gotifyURL, err := url.Parse(policy.URL)
	if err != nil {
		return err
	}
	client := gotify.NewClient(gotifyURL, &http.Client{})
	params := message.NewCreateMessageParams()
	params.Body = &gmodels.MessageExternal{
		Title:    "Offline alert - " + entity.Name,
		Message:  fmt.Sprintf("Your application <%s> just went offline at %s, error message: %s", entity.Name, offlineTime.Local().Format(time.DateTime), msg),
		Priority: 7,
	}
	_, err = client.Message.CreateMessage(params, auth.TokenAuth(policy.Token))
	if err != nil {
		return err
	}
	return nil
}

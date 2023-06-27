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

func AlertOffline(entity *models.Application, policy *models.EnvoyGotify, msg string) error {
	url, err := url.Parse(policy.URL)
	if err != nil {
		return err
	}
	client := gotify.NewClient(url, &http.Client{})
	params := message.NewCreateMessageParams()
	params.Body = &gmodels.MessageExternal{
		Title:    "Offline alert - " + entity.Name,
		Message:  fmt.Sprintf("Your application <%s> just went offline at %s, error message: %s", entity.Name, time.Now().Format(time.RFC3339), msg),
		Priority: 7,
	}
	_, err = client.Message.CreateMessage(params, auth.TokenAuth(policy.Token))
	if err != nil {
		return err
	}
	return nil
}

package services

import (
	"encoding/json"
	"net/http"

	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
)

func TelegramCallbackService() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var update tgbotapi.Update
		err := json.NewDecoder(ctx.Request.Body).Decode(&update)
		if err != nil || update.Message == nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				api, err := tgbotapi.NewBotAPIWithAPIEndpoint(storage.Config().GetString("telegram_bot_token"), storage.Config().GetString("telegram_bot_api_endpoint")+"bot%s/%s")
				if err != nil {
					glgf.Error(err)
					break
				}
				_, err = api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to Bellis envoy!"))
				if err != nil {
					glgf.Error(err)
					break
				}
			default:
				glgf.Warn("unknown command:", update.Message.Command())
			}
		}

	}
}

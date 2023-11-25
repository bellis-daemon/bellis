package services

import (
	"encoding/json"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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
		api, err := tgbotapi.NewBotAPIWithAPIEndpoint(storage.Config().TelegramBotToken, storage.Config().TelegramBotApiEndpoint+"/bot%s/%s")
		if err != nil {
			glgf.Error(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				captcha := update.Message.CommandArguments()
				glgf.Debug("telegram join: ", captcha, update.Message.Chat.ID)
				text := "Welcome to Bellis envoy!"
				func() {
					if val, err := storage.Redis().Get(ctx, captcha).Result(); err == nil {
						id, err := primitive.ObjectIDFromHex(val)
						if err != nil {
							glgf.Error(err)
							return
						}
						var user models.User
						err = storage.CUser.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
						if err != nil {
							glgf.Error(err)
							return
						}
						err = user.SetProfile(ctx, models.IsEnvoyTelegram, &models.EnvoyTelegram{
							ID:     primitive.NewObjectID(),
							ChatId: update.Message.Chat.ID,
						})
						if err != nil {
							glgf.Error(err)
							return
						}
						text = "Welcome to Bellis envoy, successfully bind to user: " + user.Email
					}
				}()
				_, err = api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
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

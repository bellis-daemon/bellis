package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

// TelegramCallbackService handles the callback from the Telegram bot.
// It decodes the incoming update, processes the command, and sends a response back to the user.
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
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		reply.ParseMode = tgbotapi.ModeMarkdown
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				captcha := update.Message.CommandArguments()
				glgf.Debug("telegram join: ", captcha, update.Message.Chat.ID)
				reply.Text = "Welcome to Bellis envoy!"
				if val, err := storage.Redis().Get(ctx, captcha).Result(); err == nil {
					id, err := primitive.ObjectIDFromHex(val)
					if err != nil {
						glgf.Warn(err)
						break
					}
					var user models.User
					err = storage.CUser.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
					if err != nil {
						glgf.Warn(err)
						break
					}
					err = user.SetProfile(ctx, models.IsEnvoyTelegram, &models.EnvoyTelegram{
						EnvoyHeader: models.EnvoyHeader{
							UserID:    user.ID,
							CreatedAt: time.Now(),
						},
						ID:     primitive.NewObjectID(),
						ChatID: update.Message.Chat.ID,
					})
					if err != nil {
						glgf.Warn(err)
						break
					}
					reply.Text = "Welcome to Bellis envoy, successfully bind to user: " + user.Email
				} else {
					glgf.Warn(err)
					break
				}
			case "status":
				entityName := update.Message.CommandArguments()
				var policy models.EnvoyTelegram
				err := storage.CEnvoyTelegram.FindOne(ctx, bson.M{"ChatID": update.Message.Chat.ID}).Decode(&policy)
				if err != nil {
					if errors.Is(err, mongo.ErrNoDocuments) {
						reply.Text = "Telegram user not registered to bellis, please register first."
					} else {
						glgf.Error(err)
						reply.Text = "Internal server error"
					}
					break
				}
				if entityName == "" || entityName == "all" {
					var entities []models.Application
					all, err := storage.CEntity.Find(ctx, bson.M{"UserID": policy.UserID})
					if err != nil {
						if errors.Is(err, mongo.ErrNoDocuments) {
							reply.Text = "No entity found."
						} else {
							glgf.Error(err)
							reply.Text = "Internal database error"
						}
						break
					}
					err = all.All(ctx, &entities)
					if err != nil {
						glgf.Error(err)
						reply.Text = "Internal server error"
						break
					}
					for i := range entities {
						reply.Text += getEntityStatusInline(ctx, &entities[i])
					}
				} else {
					var entity models.Application
					err := storage.CEntity.FindOne(ctx, bson.M{"Name": entityName, "UserID": policy.UserID}).Decode(&entity)
					if err != nil {
						if errors.Is(err, mongo.ErrNoDocuments) {
							reply.Text = "No entity found."
						} else {
							glgf.Error(err)
							reply.Text = "Internal database error"
						}
						break
					}
					reply.Text += getEntityStatusInline(ctx, &entity)
				}
			default:
				glgf.Warn("unknown command:", update.Message.Command())
			}
		}
		if reply.Text != "" {
			_, err = api.Send(reply)
			if err != nil {
				glgf.Error(err)
			}
		}
	}
}

func getEntityStatusInline(ctx context.Context, entity *models.Application) string {
	query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`
from(bucket: "backend")
  |> range(start: -10m)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_live")
`, entity.Scheme, entity.ID.Hex()))
	if err != nil {
		glgf.Error(err)
		return ""
	}
	for query.Next() {
		switch query.Record().Field() {
		case "c_live":
			if cast.ToBool(query.Record().Value()) {
				return fmt.Sprintf("*%s:* `online✅`\n", entity.Name)
			} else {
				return fmt.Sprintf("*%s:* `offline❌`\n", entity.Name)
			}
		}
	}
	if err != nil {
		glgf.Error(err)
		return fmt.Sprintf("*%s:* `Internal database error`\n", entity.Name)
	}
	return ""
}

package services

import (
	"context"
	"net/http"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/gin-gonic/gin"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SentrySingletonRefresh() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func sentrySingletonUserFromContext(ctx context.Context) *models.User {
	id := ctx.Value("SentrySingletonUserIDKey").(primitive.ObjectID)
	var user models.User
	err := storage.CUser.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		glgf.Error(err)
		return nil
	}
	return &user
}

func SentrySingletonAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Request-Token")
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		idHex, err := storage.QuickRCSearch[string](ctx, "SENTRY_SINGLETON_TOKEN"+token, func() (string, error) {
			var user models.User
			err := storage.CUser.FindOne(ctx, bson.M{"CustomSentries": token}).Decode(&user)
			if err != nil {
				return "", err
			}
			return user.ID.Hex(), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		id, err := primitive.ObjectIDFromHex(*idHex)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("SentrySingletonUserIDKey", id)
	}
}

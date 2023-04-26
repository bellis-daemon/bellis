package midwares

import (
	"context"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserIncomingCtx To get user model from incoming context in grpc
type UserIncomingCtx struct{}

// NeedAuthChecker All servers registered with this interceptor MUST implement this interface
type NeedAuthChecker interface {
	NeedAuth() bool
}

func GetUserFromCtx(ctx context.Context) *models.User {
	return ctx.Value(UserIncomingCtx{}).(*models.User)
}

var AuthChecker grpc.UnaryServerInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if !info.Server.(NeedAuthChecker).NeedAuth() {
		return handler(ctx, req)
	}
	user := check(ctx)
	if user == nil {
		return resp, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	return handler(context.WithValue(ctx, UserIncomingCtx{}, user), req)
}

func check(ctx context.Context) *models.User {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	var (
		user         models.User
		requestToken string
	)
	if value, ok := md["request_token"]; ok {
		requestToken = value[0]
	}
	if email, err := storage.Redis().Get(ctx, requestToken).Result(); err != nil {
		return nil
	} else {
		err = storage.CUser.FindOne(ctx, bson.M{"Email": email}).Decode(&user)
		if err != nil {
			glgf.Error(err)
			return nil
		}
		return &user
	}
}

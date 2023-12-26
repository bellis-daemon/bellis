package midwares

import (
	"context"
	"github.com/bellis-daemon/bellis/common/geo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
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

func AuthChecker() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !info.Server.(NeedAuthChecker).NeedAuth() {
			return handler(ctx, req)
		}
		user := check(ctx)
		if user == nil {
			return resp, status.Error(codes.Unauthenticated, "Unauthenticated")
		}
		return handler(context.WithValue(ctx, UserIncomingCtx{}, user), req)
	}
}

func AuthCheckerStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !srv.(NeedAuthChecker).NeedAuth() {
			return handler(srv, ss)
		}
		user := check(ss.Context())
		if user == nil {
			return status.Error(codes.Unauthenticated, "Unauthenticated")
		}
		wrapped := WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ss.Context(), UserIncomingCtx{}, user)
		return handler(srv, wrapped)
	}
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
	if value := md.Get("Request-Token"); len(value) > 0 {
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
		// success login check
		onAuthed(ctx, &user)
		return &user
	}
}

func onAuthed(ctx context.Context, user *models.User) {
	go func() {
		setted, err := storage.Redis().SetNX(ctx, "ONLINE"+user.Email, true, 10*time.Minute).Result()
		if err != nil {
			glgf.Error(err)
			return
		}
		if setted == true {
			ip := ipFromContext(ctx)
			loc, err := geo.FromLocal(ip)
			if err != nil {
				glgf.Error(err)
				return
			}
			_, err = storage.CUserLoginLog.InsertOne(ctx, &models.UserLoginLog{
				ID:        primitive.NewObjectID(),
				UserID:    user.ID,
				LoginTime: time.Now(),
				Location:  loc.String(),
				Device:    deviceFromContext(ctx),
			})
			if err != nil {
				glgf.Error(err)
				return
			}
		}
	}()
}

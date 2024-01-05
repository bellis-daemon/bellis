package midwares

import (
	"context"
	"time"

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
)

// userIdFromContextKey To get user model from incoming context in grpc
type userIdFromContextKey struct{}

// NeedAuthChecker All servers registered with this interceptor MUST implement this interface
type NeedAuthChecker interface {
	NeedAuth() bool
}

func GetUserFromCtx(ctx context.Context) *models.User {
	id := ctx.Value(userIdFromContextKey{}).(primitive.ObjectID)
	var user models.User
	err := storage.CUser.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		glgf.Error(err)
		return nil
	}
	return &user
}

func AuthChecker() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if !info.Server.(NeedAuthChecker).NeedAuth() {
			return handler(ctx, req)
		}
		userId := check(ctx)
		if userId == primitive.NilObjectID {
			return resp, status.Error(codes.Unauthenticated, "Unauthenticated")
		}
		return handler(context.WithValue(ctx, userIdFromContextKey{}, userId), req)
	}
}

func AuthCheckerStream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !srv.(NeedAuthChecker).NeedAuth() {
			return handler(srv, ss)
		}
		userId := check(ss.Context())
		if userId == primitive.NilObjectID {
			return status.Error(codes.Unauthenticated, "Unauthenticated")
		}
		wrapped := WrapServerStream(ss)
		wrapped.WrappedContext = context.WithValue(ss.Context(), userIdFromContextKey{}, userId)
		return handler(srv, wrapped)
	}
}

func check(ctx context.Context) primitive.ObjectID {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return primitive.NilObjectID
	}
	var (
		requestToken string
	)
	if value := md.Get("Request-Token"); len(value) > 0 {
		requestToken = value[0]
	}
	if idHex, err := storage.Redis().Get(ctx, "LOGIN"+requestToken).Result(); err != nil {
		return primitive.NilObjectID
	} else {
		id, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			glgf.Error(err)
			return primitive.NilObjectID
		}
		// success login check
		ok := checkAuthedUser(ctx, id, requestToken)
		if !ok {
			return primitive.NilObjectID
		}
		return id
	}
}

func checkAuthedUser(ctx context.Context, userId primitive.ObjectID, requestToken string) bool {
	setted, err := storage.Redis().SetNX(ctx, "ONLINE"+userId.Hex()+requestToken, true, 10*time.Minute).Result()
	if err != nil {
		glgf.Error(err)
		return false
	}
	storage.Redis().Expire(ctx, "ONLINE"+userId.Hex()+requestToken, 10*time.Minute)
	if setted {
		count, err := storage.CUser.CountDocuments(ctx, bson.M{"_id": userId})
		if err != nil {
			glgf.Error(err)
			return false
		}
		if count == 0 {
			storage.Redis().Del(ctx, "ONLINE"+userId.Hex()+requestToken)
			storage.Redis().Del(ctx, "LOGIN"+requestToken)
			return false
		}
		storage.Redis().Expire(ctx, "LOGIN"+requestToken, 30*24*time.Hour)
		ip := ipFromContext(ctx)
		var locString string
		loc, err := geo.FromLocal(ip)
		if err != nil {
			glgf.Error(err)
			locString = loc.String()
		} else {
			locString = "Unknown Location"
		}
		_, err = storage.CUserLoginLog.InsertOne(ctx, &models.UserLoginLog{
			ID:         primitive.NewObjectID(),
			UserID:     userId,
			LoginTime:  time.Now(),
			Location:   locString,
			Device:     deviceFromContext(ctx),
			DeviceType: deviceTypeFromContext(ctx),
		})
		if err != nil {
			glgf.Error(err)
		}
	}
	return true
}

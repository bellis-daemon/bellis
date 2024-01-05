package auth

import (
	"context"
	"time"

	"github.com/bellis-daemon/bellis/common/cache"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/producer"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type handler struct{}

func (handler) NeedAuth() bool {
	return false
}

func (handler) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	var user models.User
	err := storage.CUser.FindOne(ctx, bson.M{"Email": request.Email}).Decode(&user)
	if err != nil {
		glgf.Error(err)
		return &LoginResponse{}, status.Error(codes.InvalidArgument, "Cant find this user")
	}
	if !user.CheckPassword(request.Password) {
		return &LoginResponse{}, status.Error(codes.InvalidArgument, "Wrong password")
	}
	token := cryptoo.RandToken()
	err = storage.Redis().Set(ctx, "LOGIN"+token, user.ID.Hex(), 30*24*time.Hour).Err()
	if err != nil {
		glgf.Error(err)
		return &LoginResponse{}, status.Error(codes.Internal, "Redis error")
	}
	return &LoginResponse{
		Token: token,
	}, nil
}

func (handler) Register(ctx context.Context, request *RegisterRequest) (*emptypb.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count != 0 {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "User already exist")
		}
	}

	ok, err := cache.CaptchaCheck(request.Email, request.Captcha)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "Cant check captcha")
	}
	if !ok {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "Wrong captcha")
	}
	user := models.NewUser()
	user.Email = request.Email
	err = storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		_, err := storage.CUser.InsertOne(ctx, user)
		if err != nil {
			return err
		}
		return user.SetPassword(ctx, request.Password)
	})
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "DB Error %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h handler) GetRegisterCaptcha(ctx context.Context, request *RegisterCaptchaRequest) (*emptypb.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count != 0 {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "Email already exist")
		}
	}
	err := producer.EnvoyCaptchaToEmail(ctx, request.Email)
	if err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "Cant send captcha to email")
	}
	return &emptypb.Empty{}, nil
}

func (handler) GetForgetCaptcha(ctx context.Context, request *ForgetCaptchaRequest) (*emptypb.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count == 0 {
			return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "User does not exist")
		}
	}
	err := producer.EnvoyCaptchaToEmail(ctx, request.Email)
	if err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, "Cant send captcha to email")
	}
	return &emptypb.Empty{}, nil
}

func (handler) ForgetChangePassword(ctx context.Context, request *ForgetChangePasswordRequest) (*emptypb.Empty, error) {
	ok, err := cache.CaptchaCheck(request.Email, request.Captcha)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "Cant check captcha")
	}
	if !ok {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "Wrong captcha")
	}
	var user models.User
	err = storage.CUser.FindOne(ctx, bson.M{
		"Email": request.Email,
	}).Decode(&user)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "DB error")
	}
	err = user.SetPassword(ctx, request.Password)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, "DB error")
	}
	return &emptypb.Empty{}, nil
}

func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterAuthServiceServer(server, &handler{})
		return "Auth"
	})
}

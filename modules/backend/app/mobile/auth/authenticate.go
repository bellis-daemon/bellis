package auth

import (
	"context"
	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/producer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
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
	err = storage.Redis().Set(ctx, token, user.Email, 30*24*time.Hour).Err()
	if err != nil {
		glgf.Error(err)
		return &LoginResponse{}, status.Error(codes.Internal, "Redis error")
	}
	return &LoginResponse{
		Token: token,
	}, nil
}

func (handler) Register(ctx context.Context, request *RegisterRequest) (*empty.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count != 0 {
			return &empty.Empty{}, status.Error(codes.InvalidArgument, "User already exist")
		}
	}
	//// todo: implement email captcha
	//result, err := storage.Redis().Get(ctx, "RCAPTCHA"+request.Email).Result()
	//if err != nil {
	//	glgf.Error(err)
	//	return &empty.Empty{}, status.Error(codes.Internal, "Redis error")
	//}
	//if request.Captcha != result {
	//	return &empty.Empty{}, status.Error(codes.InvalidArgument, "Wrong captcha")
	//}
	user := models.NewUser()
	user.Email = request.Email
	err := storage.MongoUseSession(ctx, func(sessionContext mongo.SessionContext) error {
		_, err := storage.CUser.InsertOne(ctx, user)
		if err != nil {
			return err
		}
		return user.SetPassword(ctx, request.Password)
	})
	if err != nil {
		return &empty.Empty{}, status.Errorf(codes.Internal, "DB Error %s", err.Error())
	}
	return &empty.Empty{}, nil
}

func (h handler) GetRegisterCaptcha(ctx context.Context, request *RegisterCaptchaRequest) (*empty.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count != 0 {
			return &empty.Empty{}, status.Error(codes.InvalidArgument, "Email already exist")
		}
	}
	captcha := cast.ToString(rand.Int63() % 10000)
	err := producer.EnvoyCaptchaToEmail(ctx, request.Email, captcha)
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "Cant send captcha to email")
	}
	err = storage.Redis().Set(ctx, "RCAPTCHA"+request.Email, captcha, 10*time.Minute).Err()
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "Redis error")
	}
	return &empty.Empty{}, nil
}

func (handler) GetForgetCaptcha(ctx context.Context, request *ForgetCaptchaRequest) (*empty.Empty, error) {
	if count, err := storage.CUser.CountDocuments(ctx, bson.M{
		"Email": request.Email,
	}); err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "DB error")
	} else {
		if count == 0 {
			return &empty.Empty{}, status.Error(codes.InvalidArgument, "User does not exist")
		}
	}
	captcha := cast.ToString(rand.Int63() % 10000)
	err := producer.EnvoyCaptchaToEmail(ctx, request.Email, captcha)
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "Cant send captcha to email")
	}
	err = storage.Redis().Set(ctx, "FCAPTCHA"+request.Email, captcha, 10*time.Minute).Err()
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "Redis error")
	}
	return &empty.Empty{}, nil
}

func (handler) ForgetChangePassword(ctx context.Context, request *ForgetChangePasswordRequest) (*empty.Empty, error) {
	result, err := storage.Redis().Get(ctx, "FCAPTCHA"+request.Email).Result()
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, "Redis error")
	}
	if request.Captcha != result {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "Wrong captcha")
	}
	var user models.User
	err = storage.CUser.FindOne(ctx, bson.M{
		"Email": request.Email,
	}).Decode(&user)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Internal, "DB error")
	}
	err = user.SetPassword(ctx, request.Password)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Internal, "DB error")
	}
	return &empty.Empty{}, nil
}

func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterAuthServiceServer(server, &handler{})
		return "Auth"
	})
}

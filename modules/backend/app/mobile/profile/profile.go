package profile

import (
	"context"
	"github.com/bellis-daemon/bellis/common/midwares"
	"github.com/bellis-daemon/bellis/modules/backend/app/server"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type handler struct{}

func (h handler) NeedAuth() bool {
	return true
}

func (h handler) GetUserProfile(ctx context.Context, empty *emptypb.Empty) (*UserProfile, error) {
	user := midwares.GetUserFromCtx(ctx)
	return &UserProfile{
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsVip:     user.IsVip,
	}, nil
}

func init() {
	server.Register(func(server *grpc.Server) string {
		RegisterProfileServiceServer(server, &handler{})
		return "Profile"
	})
}

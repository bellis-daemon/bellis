package profile

import (
	"context"
	"errors"
	"github.com/bellis-daemon/bellis/common/generic"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile/public"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/minoic/glgf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// implements TLSServiceServer
type tlsHandler struct{}

func (this *tlsHandler) NeedAuth() bool {
	return true
}

func (this *tlsHandler) GetUserTLS(ctx context.Context, empty *emptypb.Empty) (*TLSList, error) {
	ret := &TLSList{}
	user := midwares.GetUserFromCtx(ctx)
	var tlsList []models.TLS
	find, err := storage.CTLS.Find(ctx, bson.M{"UserID": user.ID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ret, nil
		}
		glgf.Error(err)
		return ret, status.Errorf(codes.Internal, "database error while finding tls configs")
	}
	err = find.Decode(&tlsList)
	if err != nil {
		glgf.Error(err)
		return ret, status.Errorf(codes.Internal, "error decoding database result")
	}
	ret.List = generic.SliceConvert(tlsList, func(s models.TLS) *TLS {
		return &TLS{
			Id:            s.ID.Hex(),
			UserId:        s.UserID.Hex(),
			Name:          s.Name,
			TlsCA:         "hided",
			TlsCert:       "hided",
			TlsKey:        "hided",
			TlsKeyPwd:     "hided",
			TlsMinVersion: s.TLSMinVersion,
			Insecure:      s.Insecure,
		}
	})
	return ret, nil
}

func (this *tlsHandler) CreateTLS(ctx context.Context, tls *TLS) (*TLS, error) {
	//TODO implement me
	panic("implement me")
}

func (this *tlsHandler) UpdateTLS(ctx context.Context, tls *TLS) (*TLS, error) {
	//TODO implement me
	panic("implement me")
}

func (this *tlsHandler) DeleteTLS(ctx context.Context, id *public.PrimitiveID) (*emptypb.Empty, error) {
	oid, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = storage.CTLS.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "database error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterTLSServiceServer(server, &tlsHandler{})
		return "TLS"
	})
}

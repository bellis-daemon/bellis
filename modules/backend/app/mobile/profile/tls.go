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

func (this *tlsHandler) CreateTLS(ctx context.Context, cert *TLS) (*TLS, error) {
	user := midwares.GetUserFromCtx(ctx)
	tls := &models.TLS{
		ID:            primitive.NewObjectID(),
		UserID:        user.ID,
		Name:          cert.Name,
		TLSCA:         cert.TlsCA,
		TLSCert:       cert.TlsCert,
		TLSKey:        cert.TlsKey,
		TLSKeyPwd:     cert.TlsKeyPwd,
		TLSMinVersion: cert.TlsMinVersion,
		Insecure:      cert.Insecure,
	}
	_, err := tls.TLSConfig()
	if err != nil {
		return cert, status.Errorf(codes.InvalidArgument, "cant parse tls cert: %s", err.Error())
	}
	result, err := storage.CTLS.InsertOne(ctx, tls)
	if err != nil {
		return cert, status.Errorf(codes.Internal, "database error: %s", err.Error())
	}
	cert.Id = result.InsertedID.(primitive.ObjectID).Hex()
	cert.UserId = user.ID.Hex()
	return cert, nil
}

func (this *tlsHandler) UpdateTLS(ctx context.Context, cert *TLS) (*TLS, error) {
	id, err := primitive.ObjectIDFromHex(cert.Id)
	if err != nil {
		return cert, status.Errorf(codes.InvalidArgument, "invalid cert id: %s", err.Error())
	}
	user := midwares.GetUserFromCtx(ctx)
	var tls models.TLS
	err = storage.CTLS.FindOne(ctx, bson.M{"_id": id}).Decode(&tls)
	if err != nil {
		return cert, status.Errorf(codes.Internal, "cant find cert in database: %s", err.Error())
	}
	if user.ID.Hex() != tls.UserID.Hex() {
		return cert, status.Errorf(codes.PermissionDenied, "you have no permission to update this cert")
	}
	tls.TLSCA = cert.TlsCA
	tls.TLSCert = cert.TlsCert
	tls.TLSKey = cert.TlsKey
	tls.TLSKeyPwd = cert.TlsKeyPwd
	tls.TLSMinVersion = cert.TlsMinVersion
	tls.Insecure = cert.Insecure
	tls.Name = cert.Name
	_, err = storage.CTLS.ReplaceOne(ctx, bson.M{"_id": id}, tls)
	if err != nil {
		return cert, status.Errorf(codes.Internal, "database update method error: %s", err.Error())
	}
	return cert, nil
}

func (this *tlsHandler) DeleteTLS(ctx context.Context, id *public.PrimitiveID) (*emptypb.Empty, error) {
	oid, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = storage.CTLS.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "database error: %s", err.Error())
	}
	return &emptypb.Empty{}, nil
}

func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterTLSServiceServer(server, &tlsHandler{})
		return "TLS"
	})
}

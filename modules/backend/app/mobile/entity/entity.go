package entity

import (
	"context"
	"fmt"
	"github.com/bellis-daemon/bellis/common"
	"github.com/bellis-daemon/bellis/common/midwares"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/server"
	"github.com/bellis-daemon/bellis/modules/backend/producer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

// implements EntityServiceServer
type handler struct{}

func (h handler) DeleteEntity(ctx context.Context, id *EntityID) (*empty.Empty, error) {
	oid, err := primitive.ObjectIDFromHex(id.ID)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = storage.CEntity.DeleteOne(ctx, bson.M{
		"_id": oid,
	})
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, err.Error())
	}
	go func() {
		_ = producer.NoticeEntityDelete(ctx, id.ID)
	}()
	return &empty.Empty{}, nil
}

func (h handler) NewEntity(ctx context.Context, entity *Entity) (*EntityID, error) {
	e := &models.Application{
		ID:          primitive.NewObjectID(),
		Name:        entity.Name,
		Description: entity.Description,
		UserID:      midwares.GetUserFromCtx(ctx).ID,
		CreatedAt:   time.Now(),
		SchemeID:    int(entity.SchemeID),
		Active:      true,
		Options:     entity.Options.AsMap(),
	}
	_, err := storage.CEntity.InsertOne(ctx, e)
	if err != nil {
		glgf.Error(err)
		return &EntityID{}, status.Error(codes.Internal, err.Error())
	}
	go func() {
		_ = producer.NoticeEntityUpdate(ctx, e.ID.Hex(), e)
	}()
	return &EntityID{
		ID: e.ID.Hex(),
	}, nil
}

func (h handler) UpdateEntity(ctx context.Context, entity *Entity) (*empty.Empty, error) {
	oid, err := primitive.ObjectIDFromHex(entity.ID)
	if err != nil {
		glgf.Warn(err)
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "invalid entity id")
	}
	_, err = storage.CEntity.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{
			"_id":         oid,
			"Name":        entity.Name,
			"Description": entity.Description,
			"Active":      entity.Active,
			"Options":     entity.Options.AsMap(),
		},
	})
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, err.Error())
	}
	go func() {
		var entity models.Application
		err := storage.CEntity.FindOne(ctx, bson.M{"_id": oid}).Decode(&entity)
		if err != nil {
			glgf.Error(err)
			return
		}
		_ = producer.NoticeEntityUpdate(ctx, entity.ID.Hex(), &entity)
	}()
	return &empty.Empty{}, nil
}

func (h handler) GetEntity(ctx context.Context, id *EntityID) (*Entity, error) {
	oid, err := primitive.ObjectIDFromHex(id.ID)
	if err != nil {
		glgf.Warn(err)
		return &Entity{}, status.Error(codes.InvalidArgument, "invalid entity id")
	}
	var entity models.Application
	err = storage.CEntity.FindOne(ctx, bson.M{"_id": oid}).Decode(&entity)
	if err != nil {
		glgf.Error(err)
		return &Entity{}, status.Error(codes.Internal, err.Error())
	}
	options, err := structpb.NewStruct(entity.Options)
	if err != nil {
		glgf.Error(err)
		return &Entity{}, status.Error(codes.Internal, err.Error())
	}
	return &Entity{
		ID:          entity.ID.Hex(),
		Name:        entity.Name,
		Description: entity.Description,
		UserID:      entity.UserID.Hex(),
		CreatedAt:   entity.CreatedAt.Format(time.RFC3339),
		SchemeID:    int32(entity.SchemeID),
		Active:      entity.Active,
		Options:     options,
	}, nil
}

func (h handler) GetAllEntities(ctx context.Context, e *empty.Empty) (*AllEntities, error) {
	user := midwares.GetUserFromCtx(ctx)
	var entities []models.Application
	find, err := storage.CEntity.Find(ctx, bson.M{"UserID": user.ID})
	if err != nil {
		glgf.Error(err)
		return &AllEntities{}, status.Error(codes.Internal, err.Error())
	}
	err = find.All(ctx, &entities)
	if err != nil {
		glgf.Error(err)
		return &AllEntities{}, status.Error(codes.Internal, err.Error())
	}
	res := &AllEntities{}
	for i := range entities {
		options, err := structpb.NewStruct(entities[i].Options)
		if err != nil {
			glgf.Error(err)
			return &AllEntities{}, status.Error(codes.Internal, err.Error())
		}
		res.Entities = append(res.Entities, &Entity{
			ID:          entities[i].ID.Hex(),
			Name:        entities[i].Name,
			Description: entities[i].Description,
			UserID:      entities[i].UserID.Hex(),
			CreatedAt:   entities[i].CreatedAt.Format(time.RFC3339),
			SchemeID:    int32(entities[i].SchemeID),
			Active:      entities[i].Active,
			Options:     options,
		})
	}
	return res, nil
}

func (h handler) GetStatus(ctx context.Context, id *EntityID) (*EntityStatus, error) {
	entityStatus := &EntityStatus{
		ID:         id.ID,
		LiveSeries: []bool{},
	}
	query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(
		`from(bucket: "backend")
  |> range(start: -10m)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")`, common.Measurements[int(id.GetSchemeId())], id.GetID()))
	if err != nil {
		glgf.Error(err)
		return &EntityStatus{}, status.Error(codes.Internal, err.Error())
	}
	fields := map[string]interface{}{}
	for query.Next() {
		switch query.Record().Field() {
		case "c_live":
			entityStatus.SentryTime = query.Record().Time().Format(time.RFC3339)
			entityStatus.Live = cast.ToBool(query.Record().Value())
		case "c_err":
			entityStatus.ErrMessage = cast.ToString(query.Record().Value())
		default:
			fields[query.Record().Field()] = query.Record().Value()
		}
	}
	entityStatus.Fields, err = structpb.NewStruct(fields)
	if err != nil {
		glgf.Error(err)
		return &EntityStatus{}, status.Error(codes.Internal, err.Error())
	}
	query, err = storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(
		`from(bucket: "backend")
  |> range(start: -24h)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_live")
  |> filter(fn: (r) => r["id"] == "%s")
  |> aggregateWindow(every: 5m, fn: first, createEmpty: true)
  |> fill(column: "_value", value: true)
  |> yield(name: "first")`, common.Measurements[int(id.GetSchemeId())], id.GetID()))
	if err != nil {
		glgf.Error(err)
		return &EntityStatus{}, status.Error(codes.Internal, err.Error())
	}
	for query.Next() {
		entityStatus.LiveSeries = append(entityStatus.LiveSeries, cast.ToBool(query.Record().Value()))
	}
	return entityStatus, nil
}

func (h handler) GetAllStatus(ctx context.Context, e *empty.Empty) (*AllEntityStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (h handler) GetSeries(ctx context.Context, id *EntityID) (*EntitySeries, error) {
	ret := &EntitySeries{
		ID: id.GetID(),
	}
	series := map[string]interface{}{}
	query, err := storage.QueryInfluxDB.Query(ctx, fmt.Sprintf(`import "types"
import "strings"
from(bucket: "backend")
  |> range(start: -10m)
  |> sort(columns: ["_time"], desc: true)
  |> limit(n:60)
  |> sort(columns: ["_time"], desc: false)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => types.isNumeric(v: r["_value"]))
  |> filter(fn: (r) => not strings.hasPrefix(v: r["_field"], prefix: "c_"))`, common.Measurements[int(id.GetSchemeId())], id.GetID()))
	if err != nil {
		glgf.Error(err)
		return ret, status.Error(codes.Internal, err.Error())
	}
	for query.Next() {
		if series[query.Record().Field()] == nil {
			series[query.Record().Field()] = []interface{}{}
		}
		series[query.Record().Field()] = append(series[query.Record().Field()].([]interface{}), query.Record().Value())
	}
	ret.Series, err = structpb.NewStruct(series)
	if err != nil {
		glgf.Error(err)
		return ret, status.Error(codes.Internal, err.Error())
	}
	return ret, nil
}

func (h handler) NeedAuth() bool {
	return true
}

func init() {
	server.Register(func(server *grpc.Server) string {
		RegisterEntityServiceServer(server, &handler{})
		return "Entity"
	})
}

package entity

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/generic"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/assertion"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// implements EntityServiceServer
type handler struct{}

func (h handler) GetStreamAllStatus(e *emptypb.Empty, server EntityService_GetStreamAllStatusServer) error {
	user := midwares.GetUserFromCtx(server.Context())
	var entities []models.Application
	find, err := storage.CEntity.Find(server.Context(), bson.M{"UserID": user.ID})
	if err != nil {
		glgf.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	err = find.All(server.Context(), &entities)
	if err != nil {
		glgf.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			all := &AllEntityStatus{}
			for i := range entities {
				s, err := h.GetStatus(server.Context(), &EntityID{ID: entities[i].ID.Hex()})
				if err != nil {
					return status.Error(codes.Internal, err.Error())
				}
				all.Status = append(all.Status, s)
			}
			err := server.Send(all)
			if err != nil {
				return err
			}
		case <-server.Context().Done():
			return nil
		}
	}
}

func (h handler) GetOfflineLog(ctx context.Context, request *OfflineLogRequest) (*OfflineLogPage, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), request.EntityID),
	)
	if err != nil {
		return &OfflineLogPage{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	eid, err := primitive.ObjectIDFromHex(request.EntityID)
	if err != nil {
		return &OfflineLogPage{}, status.Error(codes.InvalidArgument, err.Error())
	}
	options := request.Pagination.ToOptions().SetSort(bson.M{"$natural": -1})
	glgf.Debug(options)
	find, err := storage.COfflineLog.Find(ctx, bson.M{"EntityID": eid}, options)
	if err != nil {
		return &OfflineLogPage{}, status.Error(codes.Internal, err.Error())
	}
	var logs []models.OfflineLog
	err = find.All(ctx, &logs)
	if err != nil {
		return &OfflineLogPage{}, status.Error(codes.Internal, err.Error())
	}
	return &OfflineLogPage{
		Length: int32(len(logs)),
		OfflineLogs: generic.SliceConvert(logs, func(log models.OfflineLog) *OfflineLog {
			return &OfflineLog{
				OfflineTime: log.OfflineTime.Local().Format(time.DateTime),
				EnvoyType:   log.EnvoyType,
				Duration:    cryptoo.FormatDuration(log.OnlineTime.Sub(log.OfflineTime)),
				SentryLogs: generic.SliceConvert(log.SentryLogs, func(log models.SentryLog) *SentryLog {
					return &SentryLog{
						SentryName:   log.SentryName,
						SentryTime:   log.SentryTime.Local().Format(time.DateTime),
						ErrorMessage: log.ErrorMessage,
					}
				}),
			}
		}),
	}, nil
}

func (h handler) DeleteEntity(ctx context.Context, id *EntityID) (*empty.Empty, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), id.ID),
	)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
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
	go afterDeleteEntity(id.GetID())
	return &empty.Empty{}, nil
}

func (h handler) NewEntity(ctx context.Context, entity *Entity) (*EntityID, error) {
	e := &models.Application{
		ID:          primitive.NewObjectID(),
		Name:        entity.Name,
		Description: entity.Description,
		UserID:      midwares.GetUserFromCtx(ctx).ID,
		CreatedAt:   time.Now(),
		Scheme:      entity.Scheme,
		Active:      true,
		Options:     entity.Options.AsMap(),
	}
	loadPublicOptions(entity, e)
	glgf.Debugf("creating entity: %+v => %+v", entity, e)
	_, err := storage.CEntity.InsertOne(ctx, e)
	if err != nil {
		glgf.Error(err)
		return &EntityID{}, status.Error(codes.Internal, err.Error())
	}
	go afterCreateEntity(e)
	return &EntityID{
		ID: e.ID.Hex(),
	}, nil
}

func (h handler) UpdateEntity(ctx context.Context, entity *Entity) (*empty.Empty, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), entity.ID),
	)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	oid, err := primitive.ObjectIDFromHex(entity.ID)
	if err != nil {
		glgf.Warn(err)
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "invalid entity id")
	}
	e := &models.Application{}
	err = storage.CEntity.FindOne(ctx, bson.M{"_id": oid}).Decode(e)
	if err != nil {
		glgf.Warn(err)
		return &empty.Empty{}, status.Error(codes.InvalidArgument, "cant find entity by id")
	}
	e.Name = entity.GetName()
	e.Description = entity.GetDescription()
	e.Active = entity.GetActive()
	e.Options = entity.GetOptions().AsMap()
	loadPublicOptions(entity, e)
	_, err = storage.CEntity.ReplaceOne(ctx, bson.M{"_id": oid}, e)
	if err != nil {
		glgf.Error(err)
		return &empty.Empty{}, status.Error(codes.Internal, err.Error())
	}
	go afterUpdateEntity(e)
	return &empty.Empty{}, nil
}

func (h handler) GetEntity(ctx context.Context, id *EntityID) (*Entity, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), id.ID),
	)
	if err != nil {
		return &Entity{}, status.Error(codes.FailedPrecondition, err.Error())
	}
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
		Scheme:      entity.Scheme,
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
			Scheme:      entities[i].Scheme,
			Active:      entities[i].Active,
			Options:     options,
		})
	}
	return res, nil
}

func (h handler) GetStatus(ctx context.Context, id *EntityID) (*EntityStatus, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), id.ID),
	)
	if err != nil {
		return &EntityStatus{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	entityStatus := &EntityStatus{
		ID:         id.GetID(),
		LiveSeries: []bool{},
	}
	query, err := storage.QueryInfluxDB.Query(ctx,
		fmt.Sprintf(`
from(bucket: "backend")
  |> range(start: -10m)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")`,
			id.GetScheme(),
			id.GetID()))
	if err != nil {
		glgf.Error(err)
		return &EntityStatus{}, status.Error(codes.Internal, err.Error())
	}
	fields := map[string]interface{}{}
	for query.Next() {
		switch query.Record().Field() {
		case "c_live":
			entityStatus.SentryTime = query.Record().Time().Local().Format(time.TimeOnly)
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
	entityStatus.UpTime = getEntityUptime(ctx, id.GetID())
	query, err = storage.QueryInfluxDB.Query(ctx,
		fmt.Sprintf(`
from(bucket: "backend") 
  |> range(start: -24h)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_live")
  |> filter(fn: (r) => r["id"] == "%s")
  |> aggregateWindow(every: 5m, fn: first, createEmpty: true)
  |> fill(column: "_value", value: true)
  |> yield(name: "first")`,
			id.GetScheme(),
			id.GetID()))
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
	ret := &AllEntityStatus{}
	user := midwares.GetUserFromCtx(ctx)
	var entities []models.Application
	find, err := storage.CEntity.Find(ctx, bson.M{"UserID": user.ID})
	if err != nil {
		glgf.Error(err)
		return ret, status.Error(codes.Internal, err.Error())
	}
	err = find.All(ctx, &entities)
	if err != nil {
		glgf.Error(err)
		return ret, status.Error(codes.Internal, err.Error())
	}
	for i := range entities {
		s, err := h.GetStatus(ctx, &EntityID{ID: entities[i].ID.Hex()})
		if err != nil {
			glgf.Error(err)
			return ret, status.Error(codes.Internal, err.Error())
		}
		ret.Status = append(ret.Status, s)
	}
	return ret, nil
}

func (h handler) GetSeries(ctx context.Context, id *EntityID) (*EntitySeries, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), id.ID),
	)
	if err != nil {
		return &EntitySeries{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	ret := &EntitySeries{
		ID: id.GetID(),
	}
	series := map[string]interface{}{}
	query, err := storage.QueryInfluxDB.Query(ctx,
		fmt.Sprintf(`
import "types"
import "strings"
from(bucket: "backend")
  |> range(start: -10m)
  |> sort(columns: ["_time"], desc: true)
  |> limit(n:60)
  |> sort(columns: ["_time"], desc: false)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => types.isNumeric(v: r["_value"]))
  |> filter(fn: (r) => not strings.hasPrefix(v: r["_field"], prefix: "c_"))`,
			id.GetScheme(),
			id.GetID()))
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
	mobile.Register(func(server *grpc.Server) string {
		RegisterEntityServiceServer(server, &handler{})
		return "Entity"
	})
}

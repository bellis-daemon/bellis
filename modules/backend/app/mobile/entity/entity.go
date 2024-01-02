package entity

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"

	"github.com/bellis-daemon/bellis/common/cryptoo"
	"github.com/bellis-daemon/bellis/common/generic"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/storage"
	"github.com/bellis-daemon/bellis/modules/backend/app/mobile"
	"github.com/bellis-daemon/bellis/modules/backend/assertion"
	"github.com/bellis-daemon/bellis/modules/backend/midwares"
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

// GetStreamAllStatus streams all status for the user's entities using a periodic ticker.
// It retrieves the user from the context, fetches the user's entities, and then streams the status of each entity periodically using a ticker.
// The function uses goroutines to concurrently fetch the status of each entity and sends the aggregated status to the client using server streaming.
// It also handles the deadline and cancellation of the stream.
func (h handler) GetStreamAllStatus(e *emptypb.Empty, server EntityService_GetStreamAllStatusServer) error {
	ddl, ok := server.Context().Deadline()
	glgf.Success("starting streaming all status with deadline", ddl, ok)
	defer glgf.Warn("stopping stream")
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
	ticker := time.NewTicker(4900 * time.Millisecond) // 4.9s
	trigger := make(chan struct{}, 1)
	go func() {
		trigger <- struct{}{}
		for {
			select {
			case <-ticker.C:
				trigger <- struct{}{}
			case <-server.Context().Done():
				return
			}
		}
	}()
	var wg sync.WaitGroup
	var once sync.Once
	defer ticker.Stop()
	for {
		select {
		case <-trigger:
			go func() {
				start := time.Now()
				all := &AllEntityStatus{}
				wg.Add(len(entities))
				for i := range entities {
					entity := &entities[i]
					go func() {
						defer wg.Done()
						s, err := h.GetStatus(server.Context(), &EntityID{ID: entity.ID.Hex(), Scheme: &entity.Scheme})
						if err != nil {
							glgf.Error(err)
							return
						}
						all.Status = append(all.Status, s)
					}()
				}
				wg.Wait()
				glgf.Debugf("done status get for user %s in %d(ms)", user.Email, time.Since(start).Milliseconds())
				once.Do(func() {
					server.Send(all)
				})
				err := server.Send(all)
				if err != nil {
					glgf.Error(err)
				}
			}()
		case <-server.Context().Done():
			return nil
		}
	}
}

// GetOfflineLog retrieves offline logs for a specific entity based on the provided request.
// It first checks the ownership of the entity, then fetches the offline logs based on the entity ID and pagination options.
// The retrieved logs are formatted and returned as an OfflineLogPage along with any encountered errors.
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

// DeleteEntity deletes the entity based on the provided ID after checking ownership.
// It first checks the ownership of the entity, then proceeds to delete the entity from the storage.
// After the deletion, it triggers a post-deletion process asynchronously and returns an empty response or an error.
func (h handler) DeleteEntity(ctx context.Context, id *EntityID) (*emptypb.Empty, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), id.ID),
	)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	oid, err := primitive.ObjectIDFromHex(id.ID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}
	_, err = storage.CEntity.DeleteOne(ctx, bson.M{
		"_id": oid,
	})
	if err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	go afterDeleteEntity(id.GetID())
	return &emptypb.Empty{}, nil
}

// NewEntity creates a new entity based on the provided Entity object.
// It initializes a new Application model, populates it with the provided data, and inserts it into the storage.
// After the creation, it triggers a post-creation process asynchronously and returns the ID of the newly created entity.
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

// UpdateEntity updates the entity based on the provided Entity object after checking ownership.
// It first checks the ownership of the entity, then proceeds to update the entity in the storage.
// After the update, it triggers a post-update process asynchronously and returns an empty response or an error.
func (h handler) UpdateEntity(ctx context.Context, entity *Entity) (*emptypb.Empty, error) {
	err := assertion.Assert(
		checkEntityOwnershipById(ctx, midwares.GetUserFromCtx(ctx), entity.ID),
	)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	oid, err := primitive.ObjectIDFromHex(entity.ID)
	if err != nil {
		glgf.Warn(err)
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "invalid entity id")
	}
	e := &models.Application{}
	err = storage.CEntity.FindOne(ctx, bson.M{"_id": oid}).Decode(e)
	if err != nil {
		glgf.Warn(err)
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "cant find entity by id")
	}
	e.Name = entity.GetName()
	e.Description = entity.GetDescription()
	e.Active = entity.GetActive()
	e.Options = entity.GetOptions().AsMap()
	loadPublicOptions(entity, e)
	_, err = storage.CEntity.ReplaceOne(ctx, bson.M{"_id": oid}, e)
	if err != nil {
		glgf.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	go afterUpdateEntity(e)
	return &emptypb.Empty{}, nil
}

// GetEntity retrieves the entity details based on the provided EntityID after checking ownership.
// It first checks the ownership of the entity, then fetches the entity details from the storage and returns it as an Entity object.
// The options are converted to a Struct message before returning.
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

// GetAllEntities retrieves all entities belonging to the user from the storage.
// It fetches the entities based on the user's ID and returns them as a list of Entity objects within an AllEntities response.
// The options for each entity are converted to structpb.Struct format before being included in the response.
func (h handler) GetAllEntities(ctx context.Context, e *emptypb.Empty) (*AllEntities, error) {
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

// GetStatus retrieves the status of a specific entity based on the provided EntityID after checking ownership.
// It fetches various status metrics including live status, uptime, error message, response time, and live series from the storage and external sources.
// The retrieved status is returned as an EntityStatus object.
// The function uses goroutines to concurrently fetch different aspects of the entity's status and aggregates the results before returning.
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
	var wg sync.WaitGroup
	errC := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		query, err := storage.QueryInfluxDB.Query(ctx,
			fmt.Sprintf(`
from(bucket: "backend")
  |> range(start: -15s)
  |> last()
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")`,
				id.GetScheme(),
				id.GetID()))
		if err != nil {
			glgf.Error(err)
			errC <- err
			return
		}
		fields := map[string]interface{}{}
		for query.Next() {
			switch query.Record().Field() {
			case "c_live":
				entityStatus.SentryTime = query.Record().Time().Local().Format(time.TimeOnly)
				entityStatus.Live = cast.ToBool(query.Record().Value())
			case "c_err":
				entityStatus.ErrMessage = cast.ToString(query.Record().Value())
			case "c_response_time":
				entityStatus.ResponseTime = cast.ToInt64(query.Record().Value())
			default:
				fields[query.Record().Field()] = query.Record().Value()
			}
		}
		entityStatus.Fields, err = structpb.NewStruct(fields)
		if err != nil {
			glgf.Error(err)
			errC <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		entityStatus.UpTime = getEntityUptime(ctx, id.GetID())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		series, err := storage.QuickRCSearch[[]bool](ctx, "EntityLiveSeries"+id.ID, func() ([]bool, error) {
			var ret []bool
			var query *api.QueryTableResult
			if id.Scheme != nil {
				query, err = storage.QueryInfluxDB.Query(ctx,
					fmt.Sprintf(`
from(bucket: "backend") 
  |> range(start: -24h)
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["_field"] == "c_live")
  |> filter(fn: (r) => r["id"] == "%s")
  |> fill(column: "_value", value: true)
  |> aggregateWindow(every: 5m, fn: first, createEmpty: true)
  |> yield(name: "first")`,
						id.GetScheme(),
						id.GetID()))
				if err != nil {
					glgf.Error(err)
					return nil, err
				}
			} else {
				query, err = storage.QueryInfluxDB.Query(ctx,
					fmt.Sprintf(`
from(bucket: "backend") 
  |> range(start: -24h)
  |> filter(fn: (r) => r["_field"] == "c_live")
  |> filter(fn: (r) => r["id"] == "%s")
  |> fill(column: "_value", value: true)
  |> aggregateWindow(every: 5m, fn: first, createEmpty: true)
  |> yield(name: "first")`,
						id.GetID()))
				if err != nil {
					glgf.Error(err)
					return nil, err
				}
			}
			for query.Next() {
				ret = append(ret, cast.ToBool(query.Record().Value()))
			}
			return ret, nil
		}, 10*time.Minute)
		if err != nil {
			glgf.Error(err)
			errC <- err
			return
		}
		entityStatus.LiveSeries = *series
	}()

	go func() {
		wg.Wait()
		close(errC)
	}()
	for e := range errC {
		err = errors.Wrap(err, e.Error())
	}
	if err != nil {
		return &EntityStatus{}, status.Error(codes.Internal, err.Error())
	}
	return entityStatus, nil
}

// GetAllStatus retrieves the status of all entities belonging to the user from the storage.
// It fetches the status of each entity using the GetStatus function and aggregates the results into an AllEntityStatus response.
func (h handler) GetAllStatus(ctx context.Context, e *emptypb.Empty) (*AllEntityStatus, error) {
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
		s, err := h.GetStatus(ctx, &EntityID{ID: entities[i].ID.Hex(), Scheme: &entities[i].Scheme})
		if err != nil {
			glgf.Error(err)
			return ret, status.Error(codes.Internal, err.Error())
		}
		ret.Status = append(ret.Status, s)
	}
	return ret, nil
}

// GetSeries retrieves time series data for a specific entity based on the provided EntityID after checking ownership.
// It fetches the time series data from the storage and returns it as an EntitySeries object, with the data organized into a structpb.Struct format.
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
  |> filter(fn: (r) => r["_measurement"] == "%s")
  |> filter(fn: (r) => r["id"] == "%s")
  |> filter(fn: (r) => types.isNumeric(v: r["_value"]))
  |> filter(fn: (r) => not strings.hasPrefix(v: r["_field"], prefix: "c_"))
  |> sort(columns: ["_time"], desc: true)
  |> limit(n:60)
  |> sort(columns: ["_time"], desc: false)`,
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

// init registers the EntityServiceServer with the provided gRPC server using the handler implementation.
// It is called during package initialization and returns the service name "Entity" for mobile registration.
func init() {
	mobile.Register(func(server *grpc.Server) string {
		RegisterEntityServiceServer(server, &handler{})
		return "Entity"
	})
}

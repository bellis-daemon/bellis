package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CUser          *mongo.Collection
	CEntity        *mongo.Collection
	CEnvoyGotify   *mongo.Collection
	CEnvoyEmail    *mongo.Collection
	CEnvoyWebhook  *mongo.Collection
	CEnvoyTelegram *mongo.Collection
	COfflineLog    *mongo.Collection
	CEnvoyLog      *mongo.Collection
	CUserLoginLog  *mongo.Collection
	CTLS           *mongo.Collection
	client         *mongo.Client
	bellisBackend  *mongo.Database
)

func ConnectMongo() {
	ctx := context.Background()
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(Config().MongoDBURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	var err error
	client, err = mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	bellisBackend = client.Database("BellisBackend")
	CUser = bellisBackend.Collection("User")
	CEntity = bellisBackend.Collection("Entity")
	COfflineLog = bellisBackend.Collection("OfflineLog")
	CUserLoginLog = bellisBackend.Collection("UserLoginLog")
	CEnvoyEmail = bellisBackend.Collection("EnvoyEmail")
	CEnvoyGotify = bellisBackend.Collection("EnvoyGotify")
	CEnvoyWebhook = bellisBackend.Collection("EnvoyWebhook")
	CEnvoyTelegram = bellisBackend.Collection("EnvoyTelegram")
	CEnvoyLog = bellisBackend.Collection("EnvoyLog")
	CTLS = bellisBackend.Collection("TLS")
}

func MongoUseSession(ctx context.Context, f func(sessionContext mongo.SessionContext) error) error {
	return client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}
		err = f(sessionContext)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		} else {
			sessionContext.CommitTransaction(sessionContext)
		}
		return nil
	})
}

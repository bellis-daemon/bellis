package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CUser         *mongo.Collection
	CEntity       *mongo.Collection
	client        *mongo.Client
	bellisBackend *mongo.Database
)

func ConnectMongo() {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI("mongodb://mongo1,mongo2,mongo3/?replicaSet=rs0")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	bellisBackend = client.Database("BellisBackend")
	CUser = bellisBackend.Collection("User")
	CEntity = bellisBackend.Collection("Entity")
}

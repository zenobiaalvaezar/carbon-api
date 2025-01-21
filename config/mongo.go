package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoCollection *mongo.Collection
var MongoClient *mongo.Client

func ConnectMongo(ctx context.Context) {
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_DATABASE")
	collectionName := os.Getenv("MONGO_COLLECTION")

	MongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+dbHost+":"+dbPort))
	if err != nil {
		panic(err)
	}

	collections, err := MongoClient.Database(dbName).ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		panic(err)
	}

	if len(collections) == 0 {
		err = MongoClient.Database(dbName).CreateCollection(ctx, collectionName)
		if err != nil {
			panic(err)
		}
	}

	MongoCollection = MongoClient.Database(dbName).Collection(collectionName)
}

func CloseMongo(ctx context.Context) {
	err := MongoClient.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

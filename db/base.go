package db

import (
	"context"
	"time"

	"github.com/Zhousiru/inker/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var client *mongo.Client

func init() {
	err := connDB()
	if err != nil {
		panic(err)
	}

	populateIndex("users", "username", true, 1)
	populateIndex("articles", "name", true, 1)
	populateIndex("articles", "editTime", false, -1)
	populateIndex("upload.files", "filename", true, 1)
}

func populateIndex(collName string, key string, unique bool, sort int) error {
	_, err := getColl(collName).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: key, Value: bsonx.Int32(int32(sort))}},
			Options: options.Index().SetUnique(unique),
		},
	)
	return err
}

func connDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if client != nil {
		err := client.Ping(ctx, readpref.Primary())
		if err == nil {
			return nil
		}
	}

	opts := new(options.ClientOptions)
	opts.SetMaxPoolSize(16)

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Conf.MongoURI), opts)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	return nil
}

func getDB() *mongo.Database {
	return client.Database(config.Conf.DBName)
}

func getColl(name string) *mongo.Collection {
	return getDB().Collection(name)
}

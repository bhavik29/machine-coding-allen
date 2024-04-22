package db

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	Client *mongo.Client
}

var instance *db
var once sync.Once

// Single Class
func GetInstance() *db {
	once.Do(func() {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI("mongodb+srv://thakkarb97:<password>@cluster0.kprnzr3.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)

		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}

		instance = &db{
			Client: client,
		}
	})

	return instance
}

func (db *db) Disconnect() {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}

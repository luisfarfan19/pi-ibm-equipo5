package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"wantson-service/pkg/utils"
)

var MongoClient *mongo.Client

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(utils.MongoDbUrl))
	if err != nil {
		return err
	}

	MongoClient = client
	log.Println("✅ Conexión a MongoDB establecida")
	return nil
}

package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() *mongo.Client {

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Println("MONGODB_URI is not set, analytics persistence is disabled")
		return nil
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error connecting to MongoDB, analytics persistence is disabled: %v\n", err)
		return nil
	}

	return client
}

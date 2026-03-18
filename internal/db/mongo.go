package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() {
	// Placeholder for any future database initialization logic

	uri := os.Getenv("MONGODB_URI")
	docs := "www.mongodb.com/docs/drivers/go/current/"
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. " +
			"See: " + docs +
			"usage-examples/#environment-variable")
	}

	// MongoDB operations are network calls, and Go wants you to always control them safely.
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(nil)

	fmt.Println("Database successfully connected")

}

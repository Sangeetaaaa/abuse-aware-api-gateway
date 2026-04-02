package workers

import (
	"context"
	"fmt"
	"time"
	"todo-golang/internal/queue"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartEventWorker(db *mongo.Database) {
	if db == nil {
		return
	}

	go func() {
		collection := db.Collection("events")

		for event := range queue.EventChannel {
			fmt.Printf("Inserting to the database")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			_, err := collection.InsertOne(ctx, event)
			cancel()
			if err != nil {
				fmt.Printf("Failed to save abuse event: %v\n", err)
			}
		}
	}()
}

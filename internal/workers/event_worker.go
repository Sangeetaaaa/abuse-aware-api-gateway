package workers

import (
	"context"
	"fmt"
	"time"
	"todo-golang/internal/queue"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartEventWorker(db *mongo.Database) {
	// This worker will process log events from the LogChannel
	go func() {
		collection := db.Collection("events")

		for event := range queue.EventChannel {
			// Here you would typically save the event to a database or log it
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

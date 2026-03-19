package workers

import (
	"context"
	"fmt"
	"time"
	"todo-golang/internal/queue"

	"go.mongodb.org/mongo-driver/mongo"
)

func StartAbuseWorker(db *mongo.Database) {
	// This worker will process abuse events from the AbuseChannel
	go func() {
		collection := db.Collection("abuse_events")

		for event := range queue.AbuseChannel {
			// Here you would typically save the event to a database or log it
			fmt.Printf("Inserting to the database")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			fmt.Println(db)
			fmt.Println(collection)
			_, err := collection.InsertOne(ctx, event)
			cancel()
			if err != nil {
				fmt.Printf("Failed to save abuse event: %v\n", err)
			}
		}
	}()
}

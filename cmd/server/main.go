package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"todo-golang/internal/db"
	"todo-golang/internal/geo"
	"todo-golang/internal/handlers"
	"todo-golang/internal/middleware"
	"todo-golang/internal/visitor"
	"todo-golang/internal/workers"

	"github.com/gin-gonic/gin"
)

func main() {

	database := db.InitDB()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.Disconnect(ctx); err != nil {
			fmt.Println("Error disconnecting from MongoDB:", err)
		}
	}()

	workers.StartAbuseWorker(database.Database("todoapp"))

	err := geo.InitGeoDB()
	if err != nil {
		fmt.Println("Error initializing GeoDB:", err)
	}
	visitor.CleanupVisitors()

	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(1, 5)) // Limit to 1 request per second with a burst of 5

	router.StaticFile("/", "./index.html")
	router.StaticFile("/dashboard", "./dashboard.html")

	router.GET("/todos", handlers.GetTodos)
	router.PUT("/add", handlers.AddTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server is running on port", port)
	router.Run(":" + port)
}

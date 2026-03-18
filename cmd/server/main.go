package main

import (
	"fmt"
	"os"
	"todo-golang/internal/geo"
	"todo-golang/internal/handlers"
	"todo-golang/internal/middleware"
	"todo-golang/internal/visitor"

	"github.com/gin-gonic/gin"
)

func main() {

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

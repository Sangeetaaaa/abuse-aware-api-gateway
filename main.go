package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	type Todo struct {
		Title     string `json:"title" bson:"title"`
		Completed bool   `json:"completed" bson:"completed"`
	}

	todoList := []Todo{
		{Title: "Buy groceries", Completed: false},
		{Title: "Walk the dog", Completed: true},
	}

	r := gin.Default()

	// Serve the frontend
	r.StaticFile("/", "./index.html")

	// API routes
	r.GET("/todos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todoList)
	})

	// Keep original route working too
	r.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todoList)
	})

	r.PUT("/add", func(ctx *gin.Context) {
		var body Todo
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": err})
			return
		}
		todoList = append(todoList, body)
		ctx.JSON(http.StatusOK, gin.H{"message": "Todo added successfully", "todo": todoList})
	})

	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exits := visitors[ip]
	if !exits {
		limiter = rate.NewLimiter(1, 5)
	}

	return limiter
}

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			ctx.JSON(429, gin.H{"message": "Too many requests"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func main() {

	type Todo struct {
		Title     string `json:"title" bson:"title"`
		Completed bool   `json:"completed" bson:"completed"`
	}

	todoList := []Todo{
		{Title: "Buy groceries", Completed: false},
		{Title: "Walk the dog", Completed: true},
	}

	router := gin.Default()
	router.Use(RateLimitMiddleware(1, 5)) // Limit to 1 request per second with a burst of 5

	// Serve the frontend
	router.StaticFile("/", "./index.html")

	// API routes
	router.GET("/todos", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todoList)
	})

	// Keep original route working too
	router.GET("/api", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, todoList)
	})

	router.PUT("/add", func(ctx *gin.Context) {
		var body Todo
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": err})
			return
		}
		todoList = append(todoList, body)
		ctx.JSON(http.StatusOK, gin.H{"message": "Todo added successfully", "todo": todoList})
	})

	if err := router.Run(":8080"); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Visitor struct {
	limiter           *rate.Limiter
	lastSeen          time.Time
	reputationScore   int
	isBlocked         bool
	blocketUntil      time.Time
	lastRateLimitTime time.Time
}

var visitors = make(map[string]*Visitor)
var mu sync.Mutex

func getVisitor(ip string) *Visitor {
	mu.Lock()
	defer mu.Unlock()

	visitorStruct, exits := visitors[ip]
	if !exits {
		visitorStruct = &Visitor{
			limiter:         rate.NewLimiter(1, 5),
			reputationScore: 0, // Initialize IP score to 0
		}
		visitors[ip] = visitorStruct
	}

	visitorStruct.lastSeen = time.Now()

	return visitorStruct
}

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()
		visitor := getVisitor(ip)

		if visitor.isBlocked {
			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if !visitor.limiter.Allow() {
			visitor.reputationScore++
			if visitor.reputationScore == 10 {
				fmt.Printf("Blocking IP %s due to suspicious activity\n", ip)
				visitor.isBlocked = true
				visitor.blocketUntil = time.Now().Add(1 * time.Minute)
			}
			visitor.lastRateLimitTime = time.Now()
			ctx.JSON(429, gin.H{"message": "Too many requests"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// loop through visitors map
//     if now - lastSeen > 3 minutes
//         delete visitor

func cleanupVisitors() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		for {
			<-ticker.C
			mu.Lock()
			fmt.Println("Cleaning up expired visitors list")
			for ip, visitor := range visitors {
				// if time.Now().After(visitor.lastSeen.Add(3 * time.Minute)) {
				// 	fmt.Println("Deleting expired visitor:", ip)
				// 	delete(visitors, ip)
				// }

				if visitor.isBlocked {
					if time.Now().After(visitor.blocketUntil) {
						visitor.isBlocked = false
						visitor.reputationScore--
					}
				} else {
					if time.Now().After(visitor.lastRateLimitTime.Add(10 * time.Minute)) {
						fmt.Printf("Resetting reputation score for IP %s\n", ip)
						if visitor.reputationScore > 0 {
							visitor.reputationScore--
						}
					}
				}
			}
			mu.Unlock()
		}
	}()
}

// IP Reputation Service Integration

func main() {

	cleanupVisitors()

	type Todo struct {
		Title     string `json:"title" bson:"title"`
		Completed bool   `json:"completed" bson:"completed"`
	}

	todoList := []Todo{}

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

	if err := router.Run(":8000"); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

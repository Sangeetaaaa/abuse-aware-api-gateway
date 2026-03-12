package main

import (
	"fmt"
	"net/http"
	"strings"
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
	blockUntil        time.Time
	lastRateLimitTime time.Time
	notFoundCount     int
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

func isBot(ua string, pattern []string) bool {
	ua = strings.ToLower(ua)
	for _, bot := range pattern {
		if strings.Contains(ua, strings.ToLower(bot)) {
			return true
		}
	}
	return false
}

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()
		visitor := getVisitor(ip)
		suspciousUserAgents := []string{
			"curl",
			"wget",
			"bot",
			"spider",
			"crawler",
			"python-requests",
			"axios",
			"httpclient",
			"java",
			"go-http-client",
			"PostmanRuntime/7.39.1",
		}

		if ctx.Request.UserAgent() == "" || isBot(ctx.Request.UserAgent(), suspciousUserAgents) {
			fmt.Println(ctx.Request.UserAgent(), "user agency found")
			visitor.reputationScore = visitor.reputationScore + 3
		}

		if visitor.reputationScore > 10 && !visitor.isBlocked {
			fmt.Printf("Blocking IP %s due to suspicious activity\n", ip)
			visitor.isBlocked = true
			visitor.blockUntil = time.Now().Add(1 * time.Minute)
			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if visitor.isBlocked {
			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if !visitor.limiter.Allow() {
			visitor.reputationScore++
			visitor.lastRateLimitTime = time.Now()
			ctx.JSON(429, gin.H{"message": "Too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()

		status := ctx.Writer.Status()
		if status == http.StatusNotFound {
			visitor.notFoundCount++
			if visitor.notFoundCount > 5 {
				fmt.Printf("Blocking IP %s due to excessive 404 errors\n", ip)
				visitor.isBlocked = true
				visitor.blockUntil = time.Now().Add(1 * time.Minute)
				visitor.reputationScore = visitor.reputationScore + 10
				ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
				ctx.Abort()
				return
			}
		}
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
					if time.Now().After(visitor.blockUntil) {
						mu.Lock()
						visitor.isBlocked = false
						if visitor.reputationScore > 0 {
							visitor.reputationScore--
						}
						mu.Unlock()
					}
				} else {
					if time.Now().After(visitor.lastRateLimitTime.Add(10 * time.Minute)) {
						fmt.Printf("Resetting reputation score for IP %s\n", ip)
						mu.Lock()
						if visitor.reputationScore > 0 {
							visitor.reputationScore--
						}
						mu.Unlock()
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

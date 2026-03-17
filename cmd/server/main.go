package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"todo-golang/internal/geo"
	"todo-golang/internal/handlers"
	"todo-golang/internal/middleware"
	"todo-golang/internal/visitor"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Visitor struct {
	limiter            *rate.Limiter
	lastSeen           time.Time
	reputationScore    int
	isBlocked          bool
	isPermanentBlocked bool
	blockUntil         time.Time
	lastRateLimitTime  time.Time
	notFoundCount      int
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

	router.Run(":3000")
}

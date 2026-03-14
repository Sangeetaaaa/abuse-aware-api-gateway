package middleware

import (
	"fmt"
	"net/http"
	"time"
	"todo-golang/internal/geo"
	"todo-golang/internal/visitor"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()
		visitor := visitor.GetVisitor(ip)
		geoLocation := geo.GetGeoLocation(ip)

		if geoLocation.Country.ISOCode == "NL" {
			visitor.IsBlocked = true
			visitor.IsPermanentBlocked = true
		}
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
			visitor.ReputationScore = visitor.ReputationScore + 3
		}

		if visitor.ReputationScore > 10 && !visitor.IsBlocked {
			fmt.Printf("Blocking IP %s due to suspicious activity\n", ip)
			visitor.IsBlocked = true
			visitor.BlockUntil = time.Now().Add(1 * time.Minute)
			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if visitor.IsPermanentBlocked {
			ctx.JSON(403, gin.H{"message": "Site is not available in your region"})
			ctx.Abort()
			return
		}

		if visitor.IsBlocked {
			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if !visitor.Limiter.Allow() {
			visitor.ReputationScore++
			visitor.LastRateLimitTime = time.Now()
			ctx.JSON(429, gin.H{"message": "Too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()

		status := ctx.Writer.Status()
		if status == http.StatusNotFound {
			visitor.NotFoundCount++
			if visitor.NotFoundCount > 5 {
				fmt.Printf("Blocking IP %s due to excessive 404 errors\n", ip)
				visitor.IsBlocked = true
				visitor.BlockUntil = time.Now().Add(1 * time.Minute)
				visitor.ReputationScore = visitor.ReputationScore + 10
				ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
				ctx.Abort()
				return
				// ok
			}
		}
	}
}

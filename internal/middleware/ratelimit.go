package middleware

import (
	"log"
	"net"
	"net/http"
	"time"
	"todo-golang/internal/geo"
	"todo-golang/internal/models"
	"todo-golang/internal/queue"
	repo "todo-golang/internal/repositories"
	"todo-golang/internal/utils"
	"todo-golang/internal/visitor"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(limit rate.Limit, burst int) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ip := ctx.ClientIP()
		isLocalClient := isLocalIP(ip)
		visitor := visitor.GetVisitor(ip)
		geoLocation := geo.GetGeoLocation(ip)

		if !isLocalClient && geoLocation.Country.ISOCode == "NL" {
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

		if !isLocalClient && (ctx.Request.UserAgent() == "" || isBot(ctx.Request.UserAgent(), suspciousUserAgents)) {
			visitor.ReputationScore = visitor.ReputationScore + 3
		}

		if visitor.ReputationScore > 10 && !visitor.IsBlocked {
			enqueueEvent(models.Event{
				Timestamp: time.Now(),
				IP:        ip,
				Endpoint:  ctx.Request.RequestURI,
				Method:    ctx.Request.Method,
				Country:   geoLocation.Country.ISOCode,
				Action:    utils.ActionBlocked,
				Score:     visitor.ReputationScore,
			})

			log.Printf("Blocking IP %s due to suspicious activity\n", ip)
			visitor.IsBlocked = true
			visitor.BlockUntil = time.Now().Add(1 * time.Minute)

			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if visitor.IsPermanentBlocked {

			enqueueEvent(models.Event{
				Timestamp: time.Now(),
				IP:        ip,
				Endpoint:  ctx.Request.RequestURI,
				Method:    ctx.Request.Method,
				Country:   geoLocation.Country.ISOCode,
				Action:    utils.ActionBlocked,
				Score:     visitor.ReputationScore,
			})

			ctx.JSON(403, gin.H{"message": "Site is not available in your region"})
			ctx.Abort()
			return
		}

		if visitor.IsBlocked {

			enqueueEvent(models.Event{
				Timestamp: time.Now(),
				IP:        ip,
				Endpoint:  ctx.Request.RequestURI,
				Method:    ctx.Request.Method,
				Country:   geoLocation.Country.ISOCode,
				Action:    utils.ActionBlocked,
				Score:     visitor.ReputationScore,
			})

			ctx.JSON(403, gin.H{"message": "Forbidden: Your IP has been blocked due to suspicious activity"})
			ctx.Abort()
			return
		}

		if !visitor.Limiter.Allow() {
			visitor.ReputationScore++
			visitor.LastRateLimitTime = time.Now()

			enqueueEvent(models.Event{
				Timestamp: time.Now(),
				IP:        ip,
				Endpoint:  ctx.Request.RequestURI,
				Method:    ctx.Request.Method,
				Country:   geoLocation.Country.ISOCode,
				Action:    utils.ActionRateLimited,
				Score:     visitor.ReputationScore,
			})

			ctx.JSON(429, gin.H{"message": "Too many requests"})
			ctx.Abort()
			return
		}

		ctx.Next()

		status := ctx.Writer.Status()
		if status == http.StatusNotFound {
			visitor.NotFoundCount++
			if visitor.NotFoundCount > 5 {
				log.Printf("Blocking IP %s due to excessive 404 errors\n", ip)
				visitor.IsBlocked = true
				visitor.BlockUntil = time.Now().Add(1 * time.Minute)
				visitor.ReputationScore = visitor.ReputationScore + 10

				enqueueEvent(models.Event{
					Timestamp: time.Now(),
					IP:        ip,
					Endpoint:  ctx.Request.RequestURI,
					Method:    ctx.Request.Method,
					Country:   geoLocation.Country.ISOCode,
					Action:    utils.ActionBlocked,
					Score:     visitor.ReputationScore,
				})
				return
				// ok
			}
		}

		enqueueEvent(models.Event{
			Timestamp: time.Now(),
			IP:        ip,
			Endpoint:  ctx.Request.RequestURI,
			Method:    ctx.Request.Method,
			Country:   geoLocation.Country.ISOCode,
			Action:    utils.ActionAllowed,
			Score:     visitor.ReputationScore,
		})
	}
}

func enqueueEvent(event models.Event) {
	repo.StoreEvent(event)

	select {
	case queue.EventChannel <- event:
	default:
		log.Printf("Abuse channel is full, dropping event: %+v\n", event)
	}
}

func isLocalIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return parsedIP.IsLoopback() || parsedIP.IsPrivate()
}

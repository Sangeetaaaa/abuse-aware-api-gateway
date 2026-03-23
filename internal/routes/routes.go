package routes

import (
	analyticsHandlers "todo-golang/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	analytics := r.Group("/analysis")
	analytics.GET("/dashboard", analyticsHandlers.GetDashboard)
}

package handlers

import (
	"net/http"
	analyticsService "todo-golang/internal/services"

	"github.com/gin-gonic/gin"
)

func GetDashboard(ctx *gin.Context) {
	params := ctx.Query("range")

	data := analyticsService.GetDashboard(params)
	ctx.JSON(http.StatusOK, data)
}

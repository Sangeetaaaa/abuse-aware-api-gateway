package handlers

import (
	analyticsService "todo-golang/internal/services"
	// "todo-golang/internal/utils"

	"github.com/gin-gonic/gin"
)

func GetDashboard(ctx *gin.Context) {
	params := ctx.Query("range")

	data := analyticsService.GetDashboard(params)
	ctx.JSON(200, data)

	// utils.SuccessResponse(ctx, data)
}

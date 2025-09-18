package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

func SetupRoutes(router *gin.Engine, userService services.UserService, itemService services.ItemService, dailyDigestService services.DailyDigestService) {
	handler := NewHandler(userService, itemService, dailyDigestService)
	handler.SetupRoutes(router)
}

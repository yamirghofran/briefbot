package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

func SetupRoutes(router *gin.Engine, userService services.UserService, itemService services.ItemService) {
	handler := NewHandler(userService, itemService)
	handler.SetupRoutes(router)
}

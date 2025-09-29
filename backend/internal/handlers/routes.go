package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

func SetupRoutes(router *gin.Engine, userService services.UserService, itemService services.ItemService, digestService services.DigestService, podcastService services.PodcastService) {
	handler := NewHandler(userService, itemService, digestService, podcastService)
	handler.SetupRoutes(router)
}

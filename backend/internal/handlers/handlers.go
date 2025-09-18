package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

type Handler struct {
	userService        services.UserService
	itemService        services.ItemService
	dailyDigestService services.DailyDigestService
}

func NewHandler(userService services.UserService, itemService services.ItemService, dailyDigestService services.DailyDigestService) *Handler {
	return &Handler{
		userService:        userService,
		itemService:        itemService,
		dailyDigestService: dailyDigestService,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	// User routes
	userGroup := router.Group("/users")
	{
		userGroup.GET("", h.ListUsers)
		userGroup.POST("", h.CreateUser)
		userGroup.GET("/:id", h.GetUser)
		userGroup.GET("/email/:email", h.GetUserByEmail)
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.DELETE("/:id", h.DeleteUser)
	}

	// Item routes
	itemGroup := router.Group("/items")
	{
		itemGroup.POST("", h.CreateItem)
		itemGroup.GET("/:id", h.GetItem)
		itemGroup.GET("/:id/status", h.GetItemProcessingStatus)
		itemGroup.GET("/status", h.GetItemsByProcessingStatus)
		itemGroup.GET("/user/:userID", h.GetItemsByUser)
		itemGroup.GET("/user/:userID/unread", h.GetUnreadItemsByUser)
		itemGroup.PUT("/:id", h.UpdateItem)
		itemGroup.PATCH("/:id/read", h.MarkItemAsRead)
		itemGroup.DELETE("/:id", h.DeleteItem)
	}

	// Daily digest routes
	digestGroup := router.Group("/daily-digest")
	{
		digestGroup.POST("/trigger", h.TriggerDailyDigest)
		digestGroup.POST("/trigger/user/:userID", h.TriggerDailyDigestForUser)
	}
}

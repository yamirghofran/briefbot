package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yamirghofran/briefbot/internal/services"
)

type Handler struct {
	userService    services.UserService
	itemService    services.ItemService
	digestService  services.DigestService
	podcastService services.PodcastService
}

func NewHandler(userService services.UserService, itemService services.ItemService, digestService services.DigestService, podcastService services.PodcastService) *Handler {
	return &Handler{
		userService:    userService,
		itemService:    itemService,
		digestService:  digestService,
		podcastService: podcastService,
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
		itemGroup.PATCH("/:id", h.PatchItem)
		itemGroup.PATCH("/:id/read", h.MarkItemAsRead)
		itemGroup.PATCH("/:id/toggle-read", h.ToggleItemReadStatus)
		itemGroup.DELETE("/:id", h.DeleteItem)
	}

	// Podcast routes
	podcastHandler := NewPodcastHandler(h.podcastService)
	podcastGroup := router.Group("/podcasts")
	{
		// Podcast creation
		podcastGroup.POST("", podcastHandler.CreatePodcast)
		podcastGroup.POST("/from-item", podcastHandler.CreatePodcastFromSingleItem)

		// Podcast retrieval
		podcastGroup.GET("/:id", podcastHandler.GetPodcast)
		podcastGroup.GET("/user/:userID", podcastHandler.GetPodcastsByUser)
		podcastGroup.GET("/status/:status", podcastHandler.GetPodcastsByStatus)
		podcastGroup.GET("/pending", podcastHandler.GetPendingPodcasts)

		// Podcast items management
		podcastGroup.GET("/:id/items", podcastHandler.GetPodcastItems)
		podcastGroup.POST("/:id/items", podcastHandler.AddItemToPodcast)
		podcastGroup.DELETE("/:id/items/:itemID", podcastHandler.RemoveItemFromPodcast)

		// Podcast audio
		podcastGroup.GET("/:id/audio", podcastHandler.GetPodcastAudio)
		podcastGroup.GET("/:id/upload-url", podcastHandler.GeneratePodcastUploadURL)

		// Podcast management
		podcastGroup.PUT("/:id", podcastHandler.UpdatePodcast)
		podcastGroup.DELETE("/:id", podcastHandler.DeletePodcast)
	}

	// Digest routes (unified - handles both regular and integrated digests)
	digestGroup := router.Group("/digest")
	{
		digestGroup.POST("/trigger", h.TriggerDailyDigest)
		digestGroup.POST("/trigger/user/:userID", h.TriggerDailyDigestForUser)
		digestGroup.POST("/trigger/integrated", h.TriggerIntegratedDigest)
		digestGroup.POST("/trigger/integrated/user/:userID", h.TriggerIntegratedDigestForUser)
	}
}

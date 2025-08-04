package v1

import (
	"ape-go-services/website-cms-service/internal/api/v1/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")

	// Item routes
	itemHandler := handlers.NewItemHandler()
	items := v1.Group("/items")
	{
		items.GET("", itemHandler.GetAllItems)
		items.GET("/stats", itemHandler.GetItemStats)
		items.GET("/:id", itemHandler.GetItemByID)
		items.POST("", itemHandler.CreateItem)
		items.PUT("/:id", itemHandler.UpdateItem)
		items.DELETE("/:id", itemHandler.DeleteItem)
	}
}

package handlers

import (
	"ape-go-services/boilerplate-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ItemHandler handles HTTP requests for items
type ItemHandler struct {
	service *services.ItemService
}

// NewItemHandler creates a new ItemHandler
func NewItemHandler() *ItemHandler {
	return &ItemHandler{
		service: services.NewItemService(),
	}
}

// GetAllItems handles GET /api/v1/items
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.service.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"count": len(items),
	})
}

// GetItemByID handles GET /api/v1/items/:id
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// CreateItem handles POST /api/v1/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	item, err := h.service.CreateItem(req.Name, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": item,
	})
}

// UpdateItem handles PUT /api/v1/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name     string `json:"name" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.UpdateItem(id, req.Name, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item updated successfully",
	})
}

// DeleteItem handles DELETE /api/v1/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
}

// GetItemStats handles GET /api/v1/items/stats
func (h *ItemHandler) GetItemStats(c *gin.Context) {
	count, err := h.service.GetItemCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_items": count,
	})
}

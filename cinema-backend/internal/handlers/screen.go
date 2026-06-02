package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScreenHandler struct {
	service *services.ScreenService
}

func NewScreenHandler(service *services.ScreenService) *ScreenHandler {
	return &ScreenHandler{service: service}
}

func (h *ScreenHandler) Create(c *gin.Context) {
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": screen})
}

func (h *ScreenHandler) GetAll(c *gin.Context) {
	screens, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": screens})
}

func (h *ScreenHandler) GetByTheater(c *gin.Context) {
	theaterID := c.Param("theaterId")
	screens, err := h.service.GetByTheater(theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": screens})
}

func (h *ScreenHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	screen, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "screen not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": screen})
}

func (h *ScreenHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	screenID, _ := uuid.Parse(id)
	screen.ID = screenID

	if err := h.service.Update(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": screen})
}

func (h *ScreenHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "screen deleted"})
}

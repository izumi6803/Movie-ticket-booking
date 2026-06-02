package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TheaterHandler struct {
	service *services.TheaterService
}

func NewTheaterHandler(service *services.TheaterService) *TheaterHandler {
	return &TheaterHandler{service: service}
}

func (h *TheaterHandler) Create(c *gin.Context) {
	var theater models.Theater
	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&theater); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": theater})
}

func (h *TheaterHandler) GetAll(c *gin.Context) {
	theaters, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": theaters})
}

func (h *TheaterHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	theater, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "theater not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": theater})
}

func (h *TheaterHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var theater models.Theater
	if err := c.ShouldBindJSON(&theater); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	theaterID, _ := uuid.Parse(id)
	theater.ID = theaterID

	if err := h.service.Update(&theater); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": theater})
}

func (h *TheaterHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "theater deleted"})
}

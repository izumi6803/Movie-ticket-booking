package handlers

import (
	"cinema-backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	service *services.SeatService
}

func NewSeatHandler(service *services.SeatService) *SeatHandler {
	return &SeatHandler{service: service}
}

func (h *SeatHandler) GetByScreen(c *gin.Context) {
	screenID := c.Param("screenId")
	showtimeID := c.Query("showtimeId")

	// Debug logging
	fmt.Printf("GetByScreen called - screenID: '%s', showtimeID: '%s'\n", screenID, showtimeID)

	var seats interface{}
	var err error

	// If showtimeId is provided, return seats with availability status
	if showtimeID != "" {
		seats, err = h.service.GetByScreenWithStatus(screenID, showtimeID)
	} else {
		seats, err = h.service.GetByScreen(screenID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": seats})
}

func (h *SeatHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	seat, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "seat not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": seat})
}

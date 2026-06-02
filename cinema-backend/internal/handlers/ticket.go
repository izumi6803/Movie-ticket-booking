package handlers

import (
	"cinema-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	service *services.TicketService
}

func NewTicketHandler(service *services.TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

type BookTicketRequest struct {
	ShowtimeID string                   `json:"showtimeId" binding:"required"`
	Seats      []services.SeatSelection `json:"seats" binding:"required"`
}

func (h *TicketHandler) Book(c *gin.Context) {
	var req BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	showtimeID, err := uuid.Parse(req.ShowtimeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid showtime ID"})
		return
	}

	ticket, err := h.service.Book(uid, showtimeID, req.Seats)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": ticket})
}

func (h *TicketHandler) GetMyTickets(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	tickets, err := h.service.GetMyTickets(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tickets})
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	ticket, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": ticket})
}

func (h *TicketHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Cancel(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ticket cancelled"})
}

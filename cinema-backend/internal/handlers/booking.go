package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func generateBookingCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

type BookingHandler struct {
	service *services.BookingService
}

func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) Create(c *gin.Context) {
	var request struct {
		ShowtimeID string `json:"showtimeId"`
		Seats      []struct {
			SeatID    string  `json:"seatId"`
			SeatLabel string  `json:"seatLabel"`
			Price     float64 `json:"price"`
		} `json:"seats"`
		Concessions []struct {
			ConcessionID string  `json:"concessionId"`
			Quantity     int     `json:"quantity"`
			UnitPrice    float64 `json:"unitPrice"`
			TotalPrice   float64 `json:"totalPrice"`
		} `json:"concessions"`
		TotalTicketPrice     float64 `json:"totalTicketPrice"`
		TotalConcessionPrice float64 `json:"totalConcessionPrice"`
		TotalAmount          float64 `json:"totalAmount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	showtimeID, _ := uuid.Parse(request.ShowtimeID)

	// Create booking
	booking := &models.Booking{
		UserID:               uid,
		ShowtimeID:           showtimeID,
		BookingCode:          generateBookingCode(),
		TotalTicketPrice:     request.TotalTicketPrice,
		TotalConcessionPrice: request.TotalConcessionPrice,
		TotalAmount:          request.TotalAmount,
		Status:               models.BookingPending,
		PaymentStatus:        models.PaymentPending,
	}

	// Create booking seats
	for _, seat := range request.Seats {
		seatID, _ := uuid.Parse(seat.SeatID)
		booking.BookingSeats = append(booking.BookingSeats, models.BookingSeat{
			SeatID:    seatID,
			SeatLabel: seat.SeatLabel,
			Price:     seat.Price,
		})
	}

	// Create order items
	for _, item := range request.Concessions {
		concessionID, _ := uuid.Parse(item.ConcessionID)
		booking.OrderItems = append(booking.OrderItems, models.OrderItem{
			ConcessionID: concessionID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.TotalPrice,
		})
	}

	if err := h.service.Create(booking, uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Reload booking to get generated fields
	createdBooking, err := h.service.GetByID(booking.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": createdBooking})
}

func (h *BookingHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	bookings, total, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    bookings,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	// Check if userID is from URL param (admin viewing user bookings) or from auth context
	var uid uuid.UUID
	var err error

	if userIDParam := c.Param("id"); userIDParam != "" {
		// Admin is viewing specific user's bookings
		uid, err = uuid.Parse(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}
	} else {
		// Current user viewing their own bookings
		userID, _ := c.Get("userID")
		uid = userID.(uuid.UUID)
	}

	bookings, err := h.service.GetMyBookings(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": bookings})
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	booking, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": booking})
}

func (h *BookingHandler) Confirm(c *gin.Context) {
	id := c.Param("id")

	// Admin confirm booking - set to CONFIRMED
	if err := h.service.ConfirmPayment(id, map[string]string{"source": "admin"}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "booking confirmed"})
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	var request struct {
		Reason string `json:"reason,omitempty"`
	}
	c.ShouldBindJSON(&request)

	if err := h.service.CancelBooking(id, request.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "booking cancelled"})
}

func (h *BookingHandler) Refund(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "refund is not available"})
}

func (h *BookingHandler) ClearMyBookings(c *gin.Context) {
	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.service.DeleteMyBookings(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "all bookings cleared"})
}

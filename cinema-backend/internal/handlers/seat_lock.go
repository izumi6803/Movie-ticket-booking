package handlers

import (
	"cinema-backend/internal/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeatLockHandler struct {
	service *services.SeatLockService
}

func NewSeatLockHandler(service *services.SeatLockService) *SeatLockHandler {
	return &SeatLockHandler{service: service}
}

// LockSeatsRequest represents the request to lock seats
type LockSeatsRequest struct {
	ShowtimeID string   `json:"showtimeId" binding:"required"`
	SeatIDs    []string `json:"seatIds" binding:"required,min=1"`
	SeatLabels []string `json:"seatLabels" binding:"required,min=1"`
}

// LockSeats handles locking seats for booking
func (h *SeatLockHandler) LockSeats(c *gin.Context) {
	var request LockSeatsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	fmt.Printf("LockSeats called by user: %s, showtime: %s, seats: %v\n", uid, request.ShowtimeID, request.SeatIDs)

	showtimeID, err := uuid.Parse(request.ShowtimeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid showtime id"})
		return
	}

	// Parse seat IDs
	seatIDs := make([]uuid.UUID, len(request.SeatIDs))
	for i, id := range request.SeatIDs {
		seatID, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid seat id"})
			return
		}
		seatIDs[i] = seatID
	}

	lock, err := h.service.LockSeats(uid, showtimeID, seatIDs, request.SeatLabels)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"lockId":    lock.ID,
			"expiresAt": lock.ExpiresAt.Format(time.RFC3339),
			"duration":  services.DefaultLockDuration.Seconds(),
		},
	})
}

// UnlockSeats handles releasing a seat lock
func (h *SeatLockHandler) UnlockSeats(c *gin.Context) {
	lockID := c.Param("id")

	lid, err := uuid.Parse(lockID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid lock id"})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	if err := h.service.UnlockSeats(lid, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "seats unlocked"})
}

// ExtendLock handles extending a seat lock
func (h *SeatLockHandler) ExtendLock(c *gin.Context) {
	lockID := c.Param("id")

	lid, err := uuid.Parse(lockID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid lock id"})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uuid.UUID)

	lock, err := h.service.ExtendLock(lid, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"lockId":    lock.ID,
			"expiresAt": lock.ExpiresAt.Format(time.RFC3339),
		},
	})
}

// GetLockStatus gets the status of a lock
func (h *SeatLockHandler) GetLockStatus(c *gin.Context) {
	lockID := c.Param("id")

	lid, err := uuid.Parse(lockID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid lock id"})
		return
	}

	lock, err := h.service.GetLockStatus(lid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "lock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"lockId":    lock.ID,
			"status":    lock.Status,
			"expiresAt": lock.ExpiresAt.Format(time.RFC3339),
			"isValid":   lock.IsValid(),
		},
	})
}

// GetActiveLocks gets all active locks for a showtime
func (h *SeatLockHandler) GetActiveLocks(c *gin.Context) {
	showtimeID := c.Query("showtimeId")
	if showtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "showtimeId is required"})
		return
	}

	stid, err := uuid.Parse(showtimeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid showtime id"})
		return
	}

	locks, err := h.service.GetActiveLocksByShowtime(stid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": locks})
}

// ClearAllLocks clears all active locks (admin only)
func (h *SeatLockHandler) ClearAllLocks(c *gin.Context) {
	if err := h.service.ClearAllLocks(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "all locks cleared"})
}

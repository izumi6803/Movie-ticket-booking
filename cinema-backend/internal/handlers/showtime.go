package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShowtimeHandler struct {
	service *services.ShowtimeService
}

func NewShowtimeHandler(service *services.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{service: service}
}

func (h *ShowtimeHandler) Create(c *gin.Context) {
	var showtime models.Showtime
	if err := c.ShouldBindJSON(&showtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&showtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": showtime})
}

func (h *ShowtimeHandler) GetAll(c *gin.Context) {
	showtimes, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": showtimes})
}

func (h *ShowtimeHandler) GetByMovie(c *gin.Context) {
	movieID := c.Param("movieId")
	showtimes, err := h.service.GetByMovie(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": showtimes})
}

func (h *ShowtimeHandler) GetByMovieAndTheater(c *gin.Context) {
	movieID := c.Param("movieId")
	theaterID := c.Param("theaterId")
	showtimes, err := h.service.GetByMovieAndTheater(movieID, theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": showtimes})
}

func (h *ShowtimeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	showtime, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "showtime not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": showtime})
}

func (h *ShowtimeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var showtime models.Showtime
	if err := c.ShouldBindJSON(&showtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	showtimeID, _ := uuid.Parse(id)
	showtime.ID = showtimeID

	if err := h.service.Update(&showtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": showtime})
}

func (h *ShowtimeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "showtime deleted"})
}

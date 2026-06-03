package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) Create(c *gin.Context) {
	var request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Duration    int    `json:"duration"`
		Genre       string `json:"genre"`
		Rating      string `json:"rating"`
		PosterURL   string `json:"posterUrl"`
		TrailerURL  string `json:"trailerUrl"`
		ReleaseDate string `json:"releaseDate"`
		Director    string `json:"director"`
		Cast        string `json:"cast"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Parse genre from comma-separated string
	var genres []string
	if request.Genre != "" {
		for _, g := range strings.Split(request.Genre, ",") {
			g = strings.TrimSpace(g)
			if g != "" {
				genres = append(genres, g)
			}
		}
	}

	// Parse cast from comma-separated string
	var casts []string
	if request.Cast != "" {
		for _, c := range strings.Split(request.Cast, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				casts = append(casts, c)
			}
		}
	}

	movie := &models.Movie{
		Title:       request.Title,
		Description: request.Description,
		Duration:    request.Duration,
		Genre:       genres,
		Rating:      models.MovieRating(request.Rating),
		Director:    request.Director,
		Status:      models.MovieStatus(request.Status),
		Cast:        casts,
	}

	if request.PosterURL != "" {
		movie.PosterURL = &request.PosterURL
	}
	if request.TrailerURL != "" {
		movie.TrailerURL = &request.TrailerURL
	}
	if request.ReleaseDate != "" {
		// Parse release date
		releaseDate, err := time.Parse("2006-01-02", request.ReleaseDate)
		if err == nil {
			movie.ReleaseDate = &releaseDate
		}
	}

	if err := h.service.Create(movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": movie})
}

func (h *MovieHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	genre := c.Query("genre")
	status := c.Query("status")

	movies, total, err := h.service.GetAll(page, limit, search, genre, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    movies,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *MovieHandler) GetNowShowing(c *gin.Context) {
	movies, err := h.service.GetNowShowing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

func (h *MovieHandler) GetComingSoon(c *gin.Context) {
	movies, err := h.service.GetComingSoon()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movies})
}

func (h *MovieHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	movie, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

func (h *MovieHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Duration    int    `json:"duration"`
		Genre       string `json:"genre"`
		Rating      string `json:"rating"`
		PosterURL   string `json:"posterUrl"`
		TrailerURL  string `json:"trailerUrl"`
		ReleaseDate string `json:"releaseDate"`
		Director    string `json:"director"`
		Cast        string `json:"cast"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Parse genre from comma-separated string
	var genres []string
	if request.Genre != "" {
		for _, g := range strings.Split(request.Genre, ",") {
			g = strings.TrimSpace(g)
			if g != "" {
				genres = append(genres, g)
			}
		}
	}

	// Parse cast from comma-separated string
	var casts []string
	if request.Cast != "" {
		for _, c := range strings.Split(request.Cast, ",") {
			c = strings.TrimSpace(c)
			if c != "" {
				casts = append(casts, c)
			}
		}
	}

	movie := &models.Movie{
		Title:       request.Title,
		Description: request.Description,
		Duration:    request.Duration,
		Genre:       genres,
		Rating:      models.MovieRating(request.Rating),
		Director:    request.Director,
		Status:      models.MovieStatus(request.Status),
		Cast:        casts,
	}

	if request.PosterURL != "" {
		movie.PosterURL = &request.PosterURL
	}
	if request.TrailerURL != "" {
		movie.TrailerURL = &request.TrailerURL
	}
	if request.ReleaseDate != "" {
		releaseDate, err := time.Parse("2006-01-02", request.ReleaseDate)
		if err == nil {
			movie.ReleaseDate = &releaseDate
		}
	}

	// Set ID from URL
	movieID, _ := uuid.Parse(id)
	movie.ID = movieID

	if err := h.service.Update(movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

func (h *MovieHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "movie deleted"})
}

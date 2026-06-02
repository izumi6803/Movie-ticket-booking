package handlers

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ConcessionHandler struct {
	service *services.ConcessionService
}

func NewConcessionHandler(service *services.ConcessionService) *ConcessionHandler {
	return &ConcessionHandler{service: service}
}

func (h *ConcessionHandler) Create(c *gin.Context) {
	var concession models.Concession
	if err := c.ShouldBindJSON(&concession); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&concession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": concession})
}

func (h *ConcessionHandler) GetAll(c *gin.Context) {
	concessions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": concessions})
}

func (h *ConcessionHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	concessions, err := h.service.GetByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": concessions})
}

func (h *ConcessionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	concession, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "concession not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": concession})
}

func (h *ConcessionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var concession models.Concession
	if err := c.ShouldBindJSON(&concession); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	concessionID, _ := uuid.Parse(id)
	concession.ID = concessionID

	if err := h.service.Update(&concession); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": concession})
}

func (h *ConcessionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "concession deleted"})
}

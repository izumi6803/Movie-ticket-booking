package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HallRepository struct {
	db *gorm.DB
}

func NewHallRepository(db *gorm.DB) *HallRepository {
	return &HallRepository{db: db}
}

func (r *HallRepository) Create(hall *models.CinemaHall) error {
	return r.db.Create(hall).Error
}

func (r *HallRepository) FindAll() ([]models.CinemaHall, error) {
	var halls []models.CinemaHall
	err := r.db.Find(&halls).Error
	return halls, err
}

func (r *HallRepository) FindByID(id uuid.UUID) (*models.CinemaHall, error) {
	var hall models.CinemaHall
	err := r.db.First(&hall, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &hall, nil
}

func (r *HallRepository) Update(hall *models.CinemaHall) error {
	return r.db.Save(hall).Error
}

func (r *HallRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.CinemaHall{}, "id = ?", id).Error
}

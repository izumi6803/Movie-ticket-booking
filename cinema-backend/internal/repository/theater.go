package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TheaterRepository struct {
	db *gorm.DB
}

func NewTheaterRepository(db *gorm.DB) *TheaterRepository {
	return &TheaterRepository{db: db}
}

func (r *TheaterRepository) Create(theater *models.Theater) error {
	return r.db.Create(theater).Error
}

func (r *TheaterRepository) FindAll() ([]models.Theater, error) {
	var theaters []models.Theater
	err := r.db.Find(&theaters).Error
	return theaters, err
}

func (r *TheaterRepository) FindByID(id uuid.UUID) (*models.Theater, error) {
	var theater models.Theater
	err := r.db.Preload("Screens").First(&theater, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &theater, nil
}

func (r *TheaterRepository) Update(theater *models.Theater) error {
	return r.db.Save(theater).Error
}

func (r *TheaterRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Theater{}, "id = ?", id).Error
}

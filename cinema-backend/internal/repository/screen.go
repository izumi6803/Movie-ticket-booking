package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScreenRepository struct {
	db *gorm.DB
}

func NewScreenRepository(db *gorm.DB) *ScreenRepository {
	return &ScreenRepository{db: db}
}

func (r *ScreenRepository) Create(screen *models.Screen) error {
	return r.db.Create(screen).Error
}

func (r *ScreenRepository) FindAll() ([]models.Screen, error) {
	var screens []models.Screen
	err := r.db.Preload("Theater").Find(&screens).Error
	return screens, err
}

func (r *ScreenRepository) FindByTheater(theaterID uuid.UUID) ([]models.Screen, error) {
	var screens []models.Screen
	err := r.db.Preload("Theater").Where("theater_id = ?", theaterID).Find(&screens).Error
	return screens, err
}

func (r *ScreenRepository) FindByID(id uuid.UUID) (*models.Screen, error) {
	var screen models.Screen
	err := r.db.Preload("Theater").Preload("Seats").First(&screen, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &screen, nil
}

func (r *ScreenRepository) Update(screen *models.Screen) error {
	return r.db.Save(screen).Error
}

func (r *ScreenRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Screen{}, "id = ?", id).Error
}

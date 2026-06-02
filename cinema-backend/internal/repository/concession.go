package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConcessionRepository struct {
	db *gorm.DB
}

func NewConcessionRepository(db *gorm.DB) *ConcessionRepository {
	return &ConcessionRepository{db: db}
}

func (r *ConcessionRepository) Create(concession *models.Concession) error {
	return r.db.Create(concession).Error
}

func (r *ConcessionRepository) FindAll() ([]models.Concession, error) {
	var concessions []models.Concession
	err := r.db.Where("is_active = ?", true).Find(&concessions).Error
	return concessions, err
}

func (r *ConcessionRepository) FindByCategory(category string) ([]models.Concession, error) {
	var concessions []models.Concession
	err := r.db.Where("category = ? AND is_active = ?", category, true).Find(&concessions).Error
	return concessions, err
}

func (r *ConcessionRepository) FindByID(id uuid.UUID) (*models.Concession, error) {
	var concession models.Concession
	err := r.db.First(&concession, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &concession, nil
}

func (r *ConcessionRepository) Update(concession *models.Concession) error {
	return r.db.Save(concession).Error
}

func (r *ConcessionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Concession{}, "id = ?", id).Error
}

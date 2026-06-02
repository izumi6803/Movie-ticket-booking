package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) Create(seat *models.Seat) error {
	return r.db.Create(seat).Error
}

func (r *SeatRepository) FindByScreen(screenID uuid.UUID) ([]models.Seat, error) {
	var seats []models.Seat
	err := r.db.Where("screen_id = ?", screenID).Order("row_label, seat_number").Find(&seats).Error
	return seats, err
}

func (r *SeatRepository) FindByID(id uuid.UUID) (*models.Seat, error) {
	var seat models.Seat
	err := r.db.First(&seat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (r *SeatRepository) Update(seat *models.Seat) error {
	return r.db.Save(seat).Error
}

func (r *SeatRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Seat{}, "id = ?", id).Error
}

func (r *SeatRepository) DeleteByScreen(screenID uuid.UUID) error {
	return r.db.Where("screen_id = ?", screenID).Delete(&models.Seat{}).Error
}

func (r *SeatRepository) BulkCreate(seats []models.Seat) error {
	return r.db.Create(&seats).Error
}

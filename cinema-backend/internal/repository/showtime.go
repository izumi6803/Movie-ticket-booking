package repository

import (
	"cinema-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShowtimeRepository struct {
	db *gorm.DB
}

func NewShowtimeRepository(db *gorm.DB) *ShowtimeRepository {
	return &ShowtimeRepository{db: db}
}

func (r *ShowtimeRepository) Create(showtime *models.Showtime) error {
	return r.db.Create(showtime).Error
}

func (r *ShowtimeRepository) FindAll() ([]models.Showtime, error) {
	var showtimes []models.Showtime
	err := r.db.Preload("Movie").Preload("Screen.Theater").
		Where("start_time > ?", time.Now()).
		Order("start_time asc").
		Find(&showtimes).Error
	return showtimes, err
}

func (r *ShowtimeRepository) DeleteExpired() error {
	return r.db.Where("start_time < ?", time.Now()).Delete(&models.Showtime{}).Error
}

func (r *ShowtimeRepository) FindByMovie(movieID uuid.UUID) ([]models.Showtime, error) {
	var showtimes []models.Showtime
	err := r.db.Preload("Movie").Preload("Screen.Theater").
		Where("movie_id = ? AND start_time > ?", movieID, time.Now()).
		Find(&showtimes).Error
	return showtimes, err
}

func (r *ShowtimeRepository) FindByMovieAndTheater(movieID, screenID uuid.UUID) ([]models.Showtime, error) {
	var showtimes []models.Showtime
	err := r.db.Preload("Movie").Preload("Screen.Theater").
		Where("movie_id = ? AND screen_id = ? AND start_time > ?", movieID, screenID, time.Now()).
		Find(&showtimes).Error
	return showtimes, err
}

func (r *ShowtimeRepository) FindByID(id uuid.UUID) (*models.Showtime, error) {
	var showtime models.Showtime
	err := r.db.Preload("Movie").Preload("Screen.Theater").First(&showtime, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &showtime, nil
}

func (r *ShowtimeRepository) Update(showtime *models.Showtime) error {
	return r.db.Save(showtime).Error
}

func (r *ShowtimeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Showtime{}, "id = ?", id).Error
}

func (r *ShowtimeRepository) UpdateAvailableSeats(id uuid.UUID, seats int) error {
	return r.db.Model(&models.Showtime{}).Where("id = ?", id).
		Update("available_seats", gorm.Expr("available_seats - ?", seats)).Error
}

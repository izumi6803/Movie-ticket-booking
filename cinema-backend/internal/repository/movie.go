package repository

import (
	"cinema-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) FindAll(page, limit int, search string, genre string, status string) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64

	query := r.db.Model(&models.Movie{})

	if search != "" {
		query = query.Where("title ILIKE ? OR director ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if genre != "" {
		query = query.Where("? = ANY(genre)", genre)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&movies).Error
	return movies, total, err
}

func (r *MovieRepository) FindByStatus(status models.MovieStatus) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Where("status = ?", status).Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) FindByID(id uuid.UUID) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.First(&movie, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) Update(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

func (r *MovieRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Movie{}, "id = ?", id).Error
}

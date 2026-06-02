package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type MovieService struct {
	repo *repository.MovieRepository
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) Create(movie *models.Movie) error {
	return s.repo.Create(movie)
}

func (s *MovieService) GetAll(page, limit int, search, genre, status string) ([]models.Movie, int64, error) {
	return s.repo.FindAll(page, limit, search, genre, status)
}

func (s *MovieService) GetNowShowing() ([]models.Movie, error) {
	return s.repo.FindByStatus(models.MovieNowShowing)
}

func (s *MovieService) GetComingSoon() ([]models.Movie, error) {
	return s.repo.FindByStatus(models.MovieComingSoon)
}

func (s *MovieService) GetByID(id string) (*models.Movie, error) {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(movieID)
}

func (s *MovieService) Update(movie *models.Movie) error {
	return s.repo.Update(movie)
}

func (s *MovieService) Delete(id string) error {
	movieID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(movieID)
}

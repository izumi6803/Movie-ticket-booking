package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type TheaterService struct {
	repo *repository.TheaterRepository
}

func NewTheaterService(repo *repository.TheaterRepository) *TheaterService {
	return &TheaterService{repo: repo}
}

func (s *TheaterService) Create(theater *models.Theater) error {
	return s.repo.Create(theater)
}

func (s *TheaterService) GetAll() ([]models.Theater, error) {
	return s.repo.FindAll()
}

func (s *TheaterService) GetByID(id string) (*models.Theater, error) {
	theaterID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(theaterID)
}

func (s *TheaterService) Update(theater *models.Theater) error {
	return s.repo.Update(theater)
}

func (s *TheaterService) Delete(id string) error {
	theaterID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(theaterID)
}

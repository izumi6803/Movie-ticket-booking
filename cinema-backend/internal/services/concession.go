package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type ConcessionService struct {
	repo *repository.ConcessionRepository
}

func NewConcessionService(repo *repository.ConcessionRepository) *ConcessionService {
	return &ConcessionService{repo: repo}
}

func (s *ConcessionService) Create(concession *models.Concession) error {
	return s.repo.Create(concession)
}

func (s *ConcessionService) GetAll() ([]models.Concession, error) {
	return s.repo.FindAll()
}

func (s *ConcessionService) GetByCategory(category string) ([]models.Concession, error) {
	return s.repo.FindByCategory(category)
}

func (s *ConcessionService) GetByID(id string) (*models.Concession, error) {
	cid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(cid)
}

func (s *ConcessionService) Update(concession *models.Concession) error {
	return s.repo.Update(concession)
}

func (s *ConcessionService) Delete(id string) error {
	cid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(cid)
}

package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ShowtimeService struct {
	repo *repository.ShowtimeRepository
}

func NewShowtimeService(repo *repository.ShowtimeRepository) *ShowtimeService {
	return &ShowtimeService{repo: repo}
}

func (s *ShowtimeService) Create(showtime *models.Showtime) error {
	// Validate showtime
	if showtime.StartTime.Before(time.Now()) {
		return errors.New("showtime start time must be in the future")
	}
	if showtime.EndTime.Before(showtime.StartTime) {
		return errors.New("showtime end time must be after start time")
	}

	return s.repo.Create(showtime)
}

func (s *ShowtimeService) DeleteExpired() error {
	return s.repo.DeleteExpired()
}

func (s *ShowtimeService) GetAll() ([]models.Showtime, error) {
	return s.repo.FindAll()
}

func (s *ShowtimeService) GetByMovie(movieID string) ([]models.Showtime, error) {
	id, err := uuid.Parse(movieID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByMovie(id)
}

func (s *ShowtimeService) GetByMovieAndTheater(movieID, screenID string) ([]models.Showtime, error) {
	mid, err := uuid.Parse(movieID)
	if err != nil {
		return nil, err
	}
	sid, err := uuid.Parse(screenID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByMovieAndTheater(mid, sid)
}

func (s *ShowtimeService) GetByID(id string) (*models.Showtime, error) {
	showtimeID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(showtimeID)
}

func (s *ShowtimeService) Update(showtime *models.Showtime) error {
	return s.repo.Update(showtime)
}

func (s *ShowtimeService) Delete(id string) error {
	showtimeID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(showtimeID)
}

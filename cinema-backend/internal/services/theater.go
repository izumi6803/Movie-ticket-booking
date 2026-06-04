package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type TheaterService struct {
	repo          *repository.TheaterRepository
	screenService *ScreenService
}

func NewTheaterService(repo *repository.TheaterRepository, screenService *ScreenService) *TheaterService {
	return &TheaterService{repo: repo, screenService: screenService}
}

func (s *TheaterService) Create(theater *models.Theater) error {
	// Create theater first
	if err := s.repo.Create(theater); err != nil {
		return err
	}

	// Create screens based on total_screens
	if theater.TotalScreens > 0 {
		for i := 1; i <= theater.TotalScreens; i++ {
			screen := &models.Screen{
				TheaterID:   theater.ID,
				Name:        "Screen " + string(rune('0'+i)),
				ScreenType:  "standard",
				TotalRows:   10,
				SeatsPerRow: 10,
				TotalSeats:  100,
				SoundSystem: "Dolby Atmos",
			}
			if err := s.screenService.Create(screen); err != nil {
				return err
			}
		}
	}

	return nil
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

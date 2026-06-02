package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type ScreenService struct {
	repo     *repository.ScreenRepository
	seatRepo *repository.SeatRepository
}

func NewScreenService(repo *repository.ScreenRepository, seatRepo *repository.SeatRepository) *ScreenService {
	return &ScreenService{repo: repo, seatRepo: seatRepo}
}

func (s *ScreenService) Create(screen *models.Screen) error {
	if err := s.repo.Create(screen); err != nil {
		return err
	}

	// Auto-generate seats for the screen
	seats := make([]models.Seat, 0, screen.TotalSeats)
	for row := 0; row < screen.TotalRows; row++ {
		rowLabel := string(rune('A' + row))
		for seatNum := 1; seatNum <= screen.SeatsPerRow; seatNum++ {
			seatType := models.SeatStandard
			priceMultiplier := 1.0

			// VIP seats in last rows
			if row >= screen.TotalRows-2 {
				seatType = models.SeatVIP
				priceMultiplier = 1.5
			}

			seats = append(seats, models.Seat{
				ScreenID:        screen.ID,
				RowLabel:        rowLabel,
				SeatNumber:      seatNum,
				SeatType:        seatType,
				PriceMultiplier: priceMultiplier,
			})
		}
	}

	return s.seatRepo.BulkCreate(seats)
}

func (s *ScreenService) GetAll() ([]models.Screen, error) {
	return s.repo.FindAll()
}

func (s *ScreenService) GetByTheater(theaterID string) ([]models.Screen, error) {
	tid, err := uuid.Parse(theaterID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByTheater(tid)
}

func (s *ScreenService) GetByID(id string) (*models.Screen, error) {
	sid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(sid)
}

func (s *ScreenService) Update(screen *models.Screen) error {
	return s.repo.Update(screen)
}

func (s *ScreenService) Delete(id string) error {
	sid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(sid)
}

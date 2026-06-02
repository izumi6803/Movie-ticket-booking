package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"

	"github.com/google/uuid"
)

type SeatWithStatus struct {
	models.Seat
	Status string `json:"status"` // available, occupied, reserved
}

type SeatService struct {
	repo        *repository.SeatRepository
	bookingRepo *repository.BookingRepository
	lockRepo    *repository.SeatLockRepository
}

func NewSeatService(repo *repository.SeatRepository, bookingRepo *repository.BookingRepository, lockRepo *repository.SeatLockRepository) *SeatService {
	return &SeatService{repo: repo, bookingRepo: bookingRepo, lockRepo: lockRepo}
}

func (s *SeatService) GetByScreen(screenID string) ([]models.Seat, error) {
	sid, err := uuid.Parse(screenID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByScreen(sid)
}

func (s *SeatService) GetByScreenWithStatus(screenID string, showtimeID string) ([]SeatWithStatus, error) {
	sid, err := uuid.Parse(screenID)
	if err != nil {
		return nil, err
	}

	// If no showtimeID provided, return seats without status
	if showtimeID == "" {
		seats, err := s.GetByScreen(screenID)
		if err != nil {
			return nil, err
		}
		var result []SeatWithStatus
		for _, seat := range seats {
			result = append(result, SeatWithStatus{Seat: seat, Status: "available"})
		}
		return result, nil
	}

	stid, err := uuid.Parse(showtimeID)
	if err != nil {
		return nil, err
	}

	// Get all seats for the screen
	seats, err := s.repo.FindByScreen(sid)
	if err != nil {
		return nil, err
	}

	// Get booked seats for the showtime
	bookedSeatIDs, err := s.bookingRepo.GetBookedSeatsByShowtime(stid)
	if err != nil {
		return nil, err
	}

	// Get locked seats for the showtime
	lockedSeatIDs, err := s.lockRepo.GetLockedSeatsByShowtime(stid)
	if err != nil {
		return nil, err
	}

	// Create maps for quick lookup
	bookedMap := make(map[uuid.UUID]bool)
	for _, id := range bookedSeatIDs {
		bookedMap[id] = true
	}

	lockedMap := make(map[uuid.UUID]bool)
	for _, id := range lockedSeatIDs {
		lockedMap[id] = true
	}

	// Add status to each seat
	var seatsWithStatus []SeatWithStatus
	for _, seat := range seats {
		status := "available"
		if bookedMap[seat.ID] {
			status = "occupied"
		} else if lockedMap[seat.ID] {
			status = "locked"
		}
		seatsWithStatus = append(seatsWithStatus, SeatWithStatus{
			Seat:   seat,
			Status: status,
		})
	}

	return seatsWithStatus, nil
}

func (s *SeatService) GetByID(id string) (*models.Seat, error) {
	sid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(sid)
}

func (s *SeatService) Update(seat *models.Seat) error {
	return s.repo.Update(seat)
}

func (s *SeatService) Delete(id string) error {
	sid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(sid)
}

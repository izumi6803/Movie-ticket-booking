package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type TicketService struct {
	ticketRepo   *repository.TicketRepository
	showtimeRepo *repository.ShowtimeRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository, showtimeRepo *repository.ShowtimeRepository) *TicketService {
	return &TicketService{
		ticketRepo:   ticketRepo,
		showtimeRepo: showtimeRepo,
	}
}

type SeatSelection struct {
	Row    string  `json:"row"`
	Number int     `json:"number"`
	Type   string  `json:"type"`
	Price  float64 `json:"price"`
}

func (s *TicketService) Book(userID uuid.UUID, showtimeID uuid.UUID, seats []SeatSelection) (*models.Ticket, error) {
	// Get showtime
	showtime, err := s.showtimeRepo.FindByID(showtimeID)
	if err != nil {
		return nil, errors.New("showtime not found")
	}

	if showtime.AvailableSeats < len(seats) {
		return nil, errors.New("not enough seats available")
	}

	// Calculate total price
	var totalPrice float64
	for _, seat := range seats {
		totalPrice += seat.Price
	}

	// Marshal seats to JSON
	seatsJSON, err := json.Marshal(seats)
	if err != nil {
		return nil, err
	}

	// Create ticket
	ticket := &models.Ticket{
		ShowtimeID: showtimeID,
		UserID:     userID,
		Seats:      string(seatsJSON),
		TotalPrice: totalPrice,
		Status:     models.TicketPending,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	// Generate QR code data
	qrData := fmt.Sprintf("TICKET:%s|USER:%s|SHOWTIME:%s|SEATS:%s|PRICE:%.2f",
		ticket.ID, userID, showtimeID, string(seatsJSON), totalPrice)
	qrCode := base64.StdEncoding.EncodeToString([]byte(qrData))
	ticket.QRCode = &qrCode

	// Update ticket with QR code
	if err := s.ticketRepo.UpdateQRCode(ticket.ID, qrCode); err != nil {
		return nil, err
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	// Update available seats
	if err := s.showtimeRepo.UpdateAvailableSeats(showtimeID, len(seats)); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) GetMyTickets(userID uuid.UUID) ([]models.Ticket, error) {
	return s.ticketRepo.FindByUser(userID)
}

func (s *TicketService) GetByID(id string) (*models.Ticket, error) {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.ticketRepo.FindByID(ticketID)
}

func (s *TicketService) Cancel(id string) error {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.ticketRepo.Cancel(ticketID)
}

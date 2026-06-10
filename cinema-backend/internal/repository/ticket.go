package repository

import (
	"cinema-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) FindByUser(userID uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Showtime.Movie").Preload("Showtime.Screen.Theater").
		Where("user_id = ?", userID).
		Order("booking_time DESC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) FindByID(id uuid.UUID) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Preload("Showtime.Movie").Preload("Showtime.Screen.Theater").
		First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepository) Cancel(id uuid.UUID) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", id).
		Update("status", models.TicketCancelled).Error
}

func (r *TicketRepository) UpdateQRCode(id uuid.UUID, qrCode string) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", id).
		Update("qr_code", qrCode).Error
}

func (r *TicketRepository) DeleteExpired(retentionDays int) error {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	return r.db.
		Joins("JOIN showtimes ON showtimes.id = tickets.showtime_id").
		Where("showtimes.end_time < ?", cutoff).
		Delete(&models.Ticket{}).Error
}

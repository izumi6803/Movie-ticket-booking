package repository

import (
	"cinema-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindAll(page, limit int) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Booking{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("User").Offset(offset).Limit(limit).
		Order("created_at DESC").Find(&bookings).Error
	return bookings, total, err
}

func (r *BookingRepository) FindByUser(userID uuid.UUID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Showtime.Movie").Preload("Showtime.Screen.Theater").Preload("BookingSeats").
		Where("user_id = ?", userID).
		Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("Showtime.Movie").Preload("Showtime.Screen.Theater").Preload("BookingSeats").First(&booking, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

func (r *BookingRepository) UpdateStatus(id uuid.UUID, status models.BookingStatus) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *BookingRepository) UpdatePaymentStatus(id uuid.UUID, status models.PaymentStatus) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).
		Update("payment_status", status).Error
}

func (r *BookingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Booking{}, "id = ?", id).Error
}

func (r *BookingRepository) DeleteByUser(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.Booking{}).Error
}

func (r *BookingRepository) GetBookedSeatsByShowtime(showtimeID uuid.UUID) ([]uuid.UUID, error) {
	var seatIDs []uuid.UUID
	err := r.db.Model(&models.BookingSeat{}).
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("bookings.showtime_id = ? AND bookings.status IN ?", showtimeID, []models.BookingStatus{models.BookingConfirmed, models.BookingPending}).
		Pluck("booking_seats.seat_id", &seatIDs).Error
	return seatIDs, err
}

// FindPendingExpired tìm các booking PENDING đã quá thởi gian
func (r *BookingRepository) FindPendingExpired(expirationTime time.Time, bookings *[]models.Booking) error {
	return r.db.Where("status = ? AND created_at < ?", models.BookingPending, expirationTime).
		Find(bookings).Error
}

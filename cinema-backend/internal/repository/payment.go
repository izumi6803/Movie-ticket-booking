package repository

import (
	"cinema-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create tạo payment record mới
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID tìm payment theo ID
func (r *PaymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Preload("Booking").Preload("User").First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByBookingID tìm payment theo booking ID
func (r *PaymentRepository) FindByBookingID(bookingID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("booking_id = ?", bookingID).Order("created_at desc").Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindByTransactionID tìm payment theo transaction ID
func (r *PaymentRepository) FindByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update cập nhật payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// UpdateStatus cập nhật status payment
func (r *PaymentRepository) UpdateStatus(id uuid.UUID, status models.PaymentStatus) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
}

// FindPendingPayments tìm các payment đang pending
func (r *PaymentRepository) FindPendingPayments() ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.Where("status = ?", models.PaymentPending).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindExpiredPayments tìm các payment đã quá hạn (created > 15 phút)
func (r *PaymentRepository) FindExpiredPayments() ([]models.Payment, error) {
	var payments []models.Payment
	expirationTime := time.Now().Add(-15 * time.Minute)
	if err := r.db.Where("status = ? AND created_at < ?", models.PaymentPending, expirationTime).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

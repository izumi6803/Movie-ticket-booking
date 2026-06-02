package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

// CreatePayment tạo payment record mới
func (s *PaymentService) CreatePayment(bookingID string, userID string, amount float64, method models.PaymentMethod) (*models.Payment, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	payment := &models.Payment{
		BookingID: bookingUUID,
		UserID:    userUUID,
		Amount:    amount,
		Method:    method,
		Status:    models.PaymentPending,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// UpdateOrderID cập nhật order ID cho payment
func (s *PaymentService) UpdateOrderID(paymentID string, orderID string) error {
	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		return err
	}

	payment, err := s.repo.FindByID(paymentUUID)
	if err != nil {
		return err
	}

	payment.OrderID = &orderID
	return s.repo.Update(payment)
}

// UpdatePaymentStatus cập nhật status payment
func (s *PaymentService) UpdatePaymentStatus(paymentID string, status models.PaymentStatus) error {
	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		return err
	}

	return s.repo.UpdateStatus(paymentUUID, status)
}

// UpdatePaymentStatusByOrderID cập nhật status payment bằng order ID
func (s *PaymentService) UpdatePaymentStatusByOrderID(orderID string, status models.PaymentStatus, transactionID string, bankCode string) error {
	// Tìm payment theo order ID
	payments, err := s.repo.FindByBookingID(uuid.Nil) // We need to find by order ID
	if err != nil {
		return err
	}

	var payment *models.Payment
	for _, p := range payments {
		if p.OrderID != nil && *p.OrderID == orderID {
			payment = &p
			break
		}
	}

	if payment == nil {
		return errors.New("payment not found")
	}

	payment.Status = status
	if transactionID != "" {
		payment.TransactionID = &transactionID
	}
	if bankCode != "" {
		payment.BankCode = &bankCode
	}

	return s.repo.Update(payment)
}

// GetPaymentByBookingID lấy payment theo booking ID
func (s *PaymentService) GetPaymentByBookingID(bookingID string) ([]models.Payment, error) {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByBookingID(bookingUUID)
}

// GetPaymentByID lấy payment theo ID
func (s *PaymentService) GetPaymentByID(paymentID string) (*models.Payment, error) {
	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(paymentUUID)
}

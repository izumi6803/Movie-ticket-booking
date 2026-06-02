package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentMethod định nghĩa các phương thức thanh toán
type PaymentMethod string

const (
	PaymentMethodVNPay   PaymentMethod = "vnpay"
	PaymentMethodCash    PaymentMethod = "cash"
	PaymentMethodMomo    PaymentMethod = "momo"
	PaymentMethodZaloPay PaymentMethod = "zalopay"
)

// Payment model lưu trữ lịch sử thanh toán
type Payment struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BookingID     uuid.UUID      `json:"bookingId" gorm:"not null;index"`
	UserID        uuid.UUID      `json:"userId" gorm:"not null"`
	Amount        float64        `json:"amount" gorm:"not null"`
	Method        PaymentMethod  `json:"method" gorm:"default:'vnpay'"`
	Status        PaymentStatus  `json:"status" gorm:"default:'pending'"`
	TransactionID *string        `json:"transactionId,omitempty"` // Mã giao dịch từ VNPay
	OrderID       *string        `json:"orderId,omitempty"`       // Mã đơn hàng VNPay
	ResponseCode  *string        `json:"responseCode,omitempty"`  // Mã phản hồi từ VNPay
	BankCode      *string        `json:"bankCode,omitempty"`      // Mã ngân hàng
	PayDate       *time.Time     `json:"payDate,omitempty"`       // Thởi gian thanh toán
	Notes         *string        `json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Booking Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

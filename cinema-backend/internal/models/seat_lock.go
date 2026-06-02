package models

import (
	"time"

	"github.com/google/uuid"
)

type SeatLockStatus string

const (
	SeatLockActive   SeatLockStatus = "active"
	SeatLockExpired  SeatLockStatus = "expired"
	SeatLockReleased SeatLockStatus = "released"
)

// SeatLock represents a temporary lock on seats during booking process
type SeatLock struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID      `json:"userId" gorm:"not null"`
	ShowtimeID uuid.UUID      `json:"showtimeId" gorm:"not null"`
	SeatIDs    []uuid.UUID    `json:"seatIds" gorm:"type:jsonb;not null;serializer:json"`
	SeatLabels []string       `json:"seatLabels" gorm:"type:jsonb;not null;serializer:json"`
	Status     SeatLockStatus `json:"status" gorm:"default:'active'"`
	ExpiresAt  time.Time      `json:"expiresAt" gorm:"not null"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

// IsExpired checks if the lock has expired
func (sl *SeatLock) IsExpired() bool {
	return time.Now().After(sl.ExpiresAt)
}

// IsValid checks if the lock is still active and not expired
func (sl *SeatLock) IsValid() bool {
	return sl.Status == SeatLockActive && !sl.IsExpired()
}

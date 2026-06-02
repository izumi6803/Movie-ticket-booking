package repository

import (
	"cinema-backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatLockRepository struct {
	db *gorm.DB
}

func NewSeatLockRepository(db *gorm.DB) *SeatLockRepository {
	return &SeatLockRepository{db: db}
}

// Create creates a new seat lock
func (r *SeatLockRepository) Create(lock *models.SeatLock) error {
	return r.db.Create(lock).Error
}

// FindByID finds a seat lock by ID
func (r *SeatLockRepository) FindByID(id uuid.UUID) (*models.SeatLock, error) {
	var lock models.SeatLock
	err := r.db.First(&lock, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lock, nil
}

// FindActiveByShowtime finds all active locks for a showtime
func (r *SeatLockRepository) FindActiveByShowtime(showtimeID uuid.UUID) ([]models.SeatLock, error) {
	var locks []models.SeatLock
	err := r.db.Where("showtime_id = ? AND status = ? AND expires_at > ?",
		showtimeID, models.SeatLockActive, time.Now()).
		Find(&locks).Error
	return locks, err
}

// FindActiveByUser finds active lock by user and showtime
func (r *SeatLockRepository) FindActiveByUser(userID, showtimeID uuid.UUID) (*models.SeatLock, error) {
	var lock models.SeatLock
	err := r.db.Where("user_id = ? AND showtime_id = ? AND status = ? AND expires_at > ?",
		userID, showtimeID, models.SeatLockActive, time.Now()).
		First(&lock).Error
	if err != nil {
		return nil, err
	}
	return &lock, nil
}

// FindBySeatAndShowtime checks if a specific seat is locked for a showtime
func (r *SeatLockRepository) FindBySeatAndShowtime(seatID, showtimeID uuid.UUID) (*models.SeatLock, error) {
	var lock models.SeatLock
	// Use PostgreSQL array containment operator with UUID
	err := r.db.Where("showtime_id = ? AND status = ? AND expires_at > ? AND ? = ANY(seat_ids)",
		showtimeID, models.SeatLockActive, time.Now(), seatID).
		First(&lock).Error
	if err != nil {
		return nil, err
	}
	return &lock, nil
}

// Update updates a lock
func (r *SeatLockRepository) Update(lock *models.SeatLock) error {
	return r.db.Save(lock).Error
}

// UpdateStatus updates the status of a lock
func (r *SeatLockRepository) UpdateStatus(id uuid.UUID, status models.SeatLockStatus) error {
	return r.db.Model(&models.SeatLock{}).Where("id = ?", id).
		Update("status", status).Error
}

// ExtendExpiry extends the expiry time of a lock
func (r *SeatLockRepository) ExtendExpiry(id uuid.UUID, newExpiry time.Time) error {
	return r.db.Model(&models.SeatLock{}).Where("id = ?", id).
		Update("expires_at", newExpiry).Error
}

// ReleaseLock releases a lock (sets status to released)
func (r *SeatLockRepository) ReleaseLock(id uuid.UUID) error {
	return r.UpdateStatus(id, models.SeatLockReleased)
}

// CleanupExpired removes expired locks
func (r *SeatLockRepository) CleanupExpired() error {
	return r.db.Model(&models.SeatLock{}).
		Where("status = ? AND expires_at <= ?", models.SeatLockActive, time.Now()).
		Update("status", models.SeatLockExpired).Error
}

// ClearAllLocks clears all active locks
func (r *SeatLockRepository) ClearAllLocks() error {
	return r.db.Model(&models.SeatLock{}).
		Where("status = ?", models.SeatLockActive).
		Update("status", models.SeatLockReleased).Error
}

// GetLockedSeatsByShowtime returns all seat IDs that are currently locked for a showtime
func (r *SeatLockRepository) GetLockedSeatsByShowtime(showtimeID uuid.UUID) ([]uuid.UUID, error) {
	var locks []models.SeatLock
	err := r.db.Where("showtime_id = ? AND status = ? AND expires_at > ?",
		showtimeID, models.SeatLockActive, time.Now()).
		Find(&locks).Error
	if err != nil {
		return nil, err
	}

	// Flatten all seat IDs from locks
	var seatIDs []uuid.UUID
	for _, lock := range locks {
		seatIDs = append(seatIDs, lock.SeatIDs...)
	}
	return seatIDs, nil
}

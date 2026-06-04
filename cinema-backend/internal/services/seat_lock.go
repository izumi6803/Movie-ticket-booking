package services

import (
	"cinema-backend/internal/models"
	"cinema-backend/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// DefaultLockDuration is the default time a seat lock is valid (10 minutes)
	DefaultLockDuration = 10 * time.Minute
	// MaxLockDuration is the maximum time a seat can be locked (15 minutes)
	MaxLockDuration = 15 * time.Minute
	// CleanupInterval is how often we cleanup expired locks
	CleanupInterval = 1 * time.Minute
	// MaxLocksPerHour is the maximum number of locks a user can create per hour
	MaxLocksPerHour = 5
	// LockCooldown is the cooldown period between locks (2 minutes)
	LockCooldown = 2 * time.Minute
)

type SeatLockService struct {
	lockRepo    *repository.SeatLockRepository
	bookingRepo *repository.BookingRepository
	hub         WebSocketHub
}

// WebSocketHub interface để tránh circular dependency
type WebSocketHub interface {
	BroadcastSeatLock(showtimeID string, seatIDs []string, seatLabels []string, userID string)
	BroadcastSeatUnlock(showtimeID string, seatIDs []string, seatLabels []string)
}

func NewSeatLockService(lockRepo *repository.SeatLockRepository, bookingRepo *repository.BookingRepository) *SeatLockService {
	return &SeatLockService{
		lockRepo:    lockRepo,
		bookingRepo: bookingRepo,
	}
}

// SetHub thiết lập WebSocket Hub (gọi sau khi khởi tạo)
func (s *SeatLockService) SetHub(hub WebSocketHub) {
	s.hub = hub
}

// LockSeats attempts to lock seats for a user
func (s *SeatLockService) LockSeats(userID uuid.UUID, showtimeID uuid.UUID, seatIDs []uuid.UUID, seatLabels []string) (*models.SeatLock, error) {
	if len(seatIDs) == 0 {
		return nil, errors.New("no seats provided")
	}

	if len(seatIDs) != len(seatLabels) {
		return nil, errors.New("seat IDs and labels count mismatch")
	}

	// Check if any of the requested seats are already booked
	bookedSeatIDs, err := s.bookingRepo.GetBookedSeatsByShowtime(showtimeID)
	if err != nil {
		return nil, err
	}
	bookedMap := make(map[uuid.UUID]bool)
	for _, id := range bookedSeatIDs {
		bookedMap[id] = true
	}

	// Check if any of the requested seats are already booked
	for _, seatID := range seatIDs {
		if bookedMap[seatID] {
			fmt.Printf("Seat %s is already booked\n", seatID)
			return nil, errors.New("one or more seats are already booked")
		}
	}

	fmt.Printf("No booked seats found. Proceeding to lock...\n")

	// Check rate limit - max 5 locks per hour
	recentLocks, err := s.lockRepo.FindRecentLocksByUser(userID, time.Now().Add(-1*time.Hour))
	if err == nil && len(recentLocks) >= MaxLocksPerHour {
		return nil, errors.New("you have reached the maximum number of seat locks per hour. Please try again later")
	}

	// Check cooldown - 2 minutes between locks
	if len(recentLocks) > 0 {
		lastLock := recentLocks[0]
		if time.Since(lastLock.CreatedAt) < LockCooldown {
			remainingTime := LockCooldown - time.Since(lastLock.CreatedAt)
			return nil, fmt.Errorf("please wait %d seconds before locking seats again", int(remainingTime.Seconds()))
		}
	}

	// Check if user already has an active lock for this showtime
	existingLock, err := s.lockRepo.FindActiveByUser(userID, showtimeID)
	if err != nil {
		fmt.Printf("No existing lock found for user: %v\n", err)
	}
	if err == nil && existingLock != nil {
		fmt.Printf("Found existing lock for user %s, updating...\n", userID)

		// Merge new seats with existing seats (avoid duplicates)
		seatMap := make(map[uuid.UUID]bool)
		labelMap := make(map[uuid.UUID]string)

		// Add existing seats
		for i, id := range existingLock.SeatIDs {
			seatMap[id] = true
			if i < len(existingLock.SeatLabels) {
				labelMap[id] = existingLock.SeatLabels[i]
			}
		}

		// Add new seats
		for i, id := range seatIDs {
			seatMap[id] = true
			if i < len(seatLabels) {
				labelMap[id] = seatLabels[i]
			}
		}

		// Convert back to slices
		mergedSeatIDs := make([]uuid.UUID, 0, len(seatMap))
		mergedSeatLabels := make([]string, 0, len(seatMap))
		for id := range seatMap {
			mergedSeatIDs = append(mergedSeatIDs, id)
			if label, ok := labelMap[id]; ok {
				mergedSeatLabels = append(mergedSeatLabels, label)
			}
		}

		// Update existing lock with merged seats
		existingLock.SeatIDs = mergedSeatIDs
		existingLock.SeatLabels = mergedSeatLabels
		existingLock.ExpiresAt = time.Now().Add(DefaultLockDuration)
		if err := s.lockRepo.Update(existingLock); err != nil {
			fmt.Printf("Failed to update lock: %v\n", err)
			// Don't return error, just log it
		} else {
			// Broadcast update
			if s.hub != nil {
				mergedSeatIDStrs := make([]string, len(mergedSeatIDs))
				for i, id := range mergedSeatIDs {
					mergedSeatIDStrs[i] = id.String()
				}
				s.hub.BroadcastSeatLock(showtimeID.String(), mergedSeatIDStrs, mergedSeatLabels, userID.String())
			}
			fmt.Printf("Lock updated successfully with %d seats\n", len(mergedSeatIDs))
			return existingLock, nil
		}
	}

	fmt.Printf("No existing lock found. Creating new lock...\n")

	// Create new lock
	lock := &models.SeatLock{
		UserID:     userID,
		ShowtimeID: showtimeID,
		SeatIDs:    seatIDs,
		SeatLabels: seatLabels,
		Status:     models.SeatLockActive,
		ExpiresAt:  time.Now().Add(DefaultLockDuration),
	}

	if err := s.lockRepo.Create(lock); err != nil {
		fmt.Printf("Failed to create lock: %v\n", err)
		return nil, err
	}

	fmt.Printf("New lock created successfully\n")

	// Broadcast seat lock to all clients watching this showtime
	if s.hub != nil {
		seatIDStrs := make([]string, len(seatIDs))
		for i, id := range seatIDs {
			seatIDStrs[i] = id.String()
		}
		s.hub.BroadcastSeatLock(showtimeID.String(), seatIDStrs, seatLabels, userID.String())
	}

	return lock, nil
}

// UnlockSeats releases a seat lock
func (s *SeatLockService) UnlockSeats(lockID uuid.UUID, userID uuid.UUID) error {
	lock, err := s.lockRepo.FindByID(lockID)
	if err != nil {
		return errors.New("lock not found")
	}

	if lock.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.lockRepo.ReleaseLock(lockID); err != nil {
		return err
	}

	// Broadcast seat unlock to all clients watching this showtime
	if s.hub != nil {
		seatIDStrs := make([]string, len(lock.SeatIDs))
		for i, id := range lock.SeatIDs {
			seatIDStrs[i] = id.String()
		}
		s.hub.BroadcastSeatUnlock(lock.ShowtimeID.String(), seatIDStrs, lock.SeatLabels)
	}

	return nil
}

// ExtendLock extends the expiry time of a lock
func (s *SeatLockService) ExtendLock(lockID uuid.UUID, userID uuid.UUID) (*models.SeatLock, error) {
	lock, err := s.lockRepo.FindByID(lockID)
	if err != nil {
		return nil, errors.New("lock not found")
	}

	if lock.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if lock.Status != models.SeatLockActive {
		return nil, errors.New("lock is not active")
	}

	// Calculate new expiry, but don't exceed max duration from creation
	maxExpiry := lock.CreatedAt.Add(MaxLockDuration)
	newExpiry := time.Now().Add(DefaultLockDuration)

	if newExpiry.After(maxExpiry) {
		newExpiry = maxExpiry
	}

	if err := s.lockRepo.ExtendExpiry(lockID, newExpiry); err != nil {
		return nil, err
	}

	lock.ExpiresAt = newExpiry
	return lock, nil
}

// GetLockStatus gets the status of a lock
func (s *SeatLockService) GetLockStatus(lockID uuid.UUID) (*models.SeatLock, error) {
	return s.lockRepo.FindByID(lockID)
}

// GetActiveLocksByShowtime gets all active locks for a showtime
func (s *SeatLockService) GetActiveLocksByShowtime(showtimeID uuid.UUID) ([]models.SeatLock, error) {
	return s.lockRepo.FindActiveByShowtime(showtimeID)
}

// CleanupExpiredLocks removes all expired locks
func (s *SeatLockService) CleanupExpiredLocks() error {
	return s.lockRepo.CleanupExpired()
}

// GetLockedSeatsByShowtime returns all seat IDs that are currently locked
func (s *SeatLockService) GetLockedSeatsByShowtime(showtimeID uuid.UUID) ([]uuid.UUID, error) {
	return s.lockRepo.GetLockedSeatsByShowtime(showtimeID)
}

// ClearAllLocks clears all active locks (useful for testing/debugging)
func (s *SeatLockService) ClearAllLocks() error {
	return s.lockRepo.ClearAllLocks()
}

// StartCleanupRoutine starts a background routine to cleanup expired locks
func (s *SeatLockService) StartCleanupRoutine() {
	ticker := time.NewTicker(CleanupInterval)
	go func() {
		for range ticker.C {
			s.CleanupExpiredLocks()
		}
	}()
}

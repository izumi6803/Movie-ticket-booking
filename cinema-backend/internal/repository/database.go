package repository

import (
	"cinema-backend/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(databaseURL string) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}

func Migrate(databaseDB *gorm.DB) error {
	db = databaseDB
	if err := db.AutoMigrate(
		&models.User{},
		&models.PasswordResetToken{},
		&models.Movie{},
		&models.Theater{},
		&models.Screen{},
		&models.Seat{},
		&models.Showtime{},
		&models.Ticket{},
		&models.Concession{},
		&models.Booking{},
		&models.BookingSeat{},
		&models.OrderItem{},
		&models.SeatLock{},
		&models.Payment{},
		&models.SystemSetting{},
	); err != nil {
		return err
	}

	// Keep email uniqueness for active users while allowing a deleted account to register again.
	if err := db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key").Error; err != nil {
		return err
	}
	if err := db.Exec("DROP INDEX IF EXISTS idx_users_email").Error; err != nil {
		return err
	}
	return db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active
        ON users (LOWER(email))
        WHERE deleted_at IS NULL
    `).Error
}

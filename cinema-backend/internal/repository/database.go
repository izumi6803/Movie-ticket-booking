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
	return db.AutoMigrate(
		&models.User{},
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
	)
}

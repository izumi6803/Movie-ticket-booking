package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=Tuananh6803@ dbname=cinema port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop all old constraints
	db.Exec(`ALTER TABLE showtimes DROP CONSTRAINT IF EXISTS showtimes_hall_id_fkey`)
	db.Exec(`ALTER TABLE showtimes DROP CONSTRAINT IF EXISTS fk_cinema_halls_showtimes`)
	db.Exec(`ALTER TABLE showtimes DROP CONSTRAINT IF EXISTS showtimes_screen_id_fkey`)

	// Check if hall_id column exists
	var count int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'showtimes' AND column_name = 'hall_id'
	`).Scan(&count)

	if count > 0 {
		fmt.Println("Found hall_id column, migrating to screen_id...")

		// Drop screen_id if it exists (from failed migration)
		db.Exec(`ALTER TABLE showtimes DROP COLUMN IF EXISTS screen_id`)

		// Rename hall_id to screen_id
		db.Exec(`ALTER TABLE showtimes RENAME COLUMN hall_id TO screen_id`)

		fmt.Println("Migration completed!")
	} else {
		fmt.Println("hall_id column not found, checking screen_id...")

		var screenCount int64
		db.Raw(`
			SELECT COUNT(*) 
			FROM information_schema.columns 
			WHERE table_name = 'showtimes' AND column_name = 'screen_id'
		`).Scan(&screenCount)

		if screenCount == 0 {
			fmt.Println("Adding screen_id column...")
			db.Exec(`ALTER TABLE showtimes ADD COLUMN screen_id UUID REFERENCES screens(id) ON DELETE CASCADE`)
		}
	}

	// Add the correct foreign key constraint
	db.Exec(`ALTER TABLE showtimes ADD CONSTRAINT showtimes_screen_id_fkey 
		FOREIGN KEY (screen_id) REFERENCES screens(id) ON DELETE CASCADE`)

	// Ensure available_seats column exists
	var availableSeatsCount int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'showtimes' AND column_name = 'available_seats'
	`).Scan(&availableSeatsCount)

	if availableSeatsCount == 0 {
		fmt.Println("Adding available_seats column...")
		db.Exec(`ALTER TABLE showtimes ADD COLUMN available_seats INTEGER NOT NULL DEFAULT 0`)
	}

	// Check if ticket_price column exists (old schema)
	var ticketPriceCount int64
	db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_name = 'showtimes' AND column_name = 'ticket_price'
	`).Scan(&ticketPriceCount)

	if ticketPriceCount > 0 {
		fmt.Println("Found ticket_price column, migrating to base_ticket_price...")

		// Drop base_ticket_price if it exists
		db.Exec(`ALTER TABLE showtimes DROP COLUMN IF EXISTS base_ticket_price`)

		// Rename ticket_price to base_ticket_price
		db.Exec(`ALTER TABLE showtimes RENAME COLUMN ticket_price TO base_ticket_price`)

		fmt.Println("ticket_price migration completed!")
	}

	// Fix empty booking codes
	fmt.Println("Checking for empty booking codes...")
	db.Exec(`UPDATE bookings SET booking_code = 'BOOK' || SUBSTRING(MD5(RANDOM()::TEXT), 1, 6) WHERE booking_code = '' OR booking_code IS NULL`)

	// Alter booking_code column to have default value
	fmt.Println("Altering booking_code column to add default value...")
	db.Exec(`ALTER TABLE bookings ALTER COLUMN booking_code SET DEFAULT 'BOOK' || SUBSTRING(MD5(RANDOM()::TEXT), 1, 6)`)

	fmt.Println("Migration completed successfully!")
}

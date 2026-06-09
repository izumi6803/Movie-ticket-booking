package main

import (
	"fmt"
	"log"

	"cinema-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://neondb_owner:npg_T12oWvNRdecQ@ep-quiet-resonance-aoqe46jp.c-2.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// Create admin user
	phone1 := "0123456789"
	admin := models.User{
		Name:     "Admin",
		Email:    "admin@cinema.com",
		Password: string(hashedPassword),
		Role:     "admin",
		Phone:    &phone1,
	}

	result := db.Where("email = ?", admin.Email).FirstOrCreate(&admin)
	if result.Error != nil {
		log.Fatal("Failed to create admin:", result.Error)
	}

	// Create test customer
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("customer123"), bcrypt.DefaultCost)
	phone2 := "0987654321"
	customer := models.User{
		Name:     "Test Customer",
		Email:    "customer@test.com",
		Password: string(hashedPassword2),
		Role:     "customer",
		Phone:    &phone2,
	}

	result = db.Where("email = ?", customer.Email).FirstOrCreate(&customer)
	if result.Error != nil {
		log.Fatal("Failed to create customer:", result.Error)
	}

	fmt.Println("Seed completed successfully!")
	fmt.Println("Admin: admin@cinema.com / admin123")
	fmt.Println("Customer: customer@test.com / customer123")
}

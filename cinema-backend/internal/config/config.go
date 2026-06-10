package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	CloudinaryURL string
	VNPay         VNPayConfig
}

type VNPayConfig struct {
	TmnCode     string
	HashSecret  string
	Endpoint    string
	ReturnURL   string
	Environment string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:          getEnv("PORT", "3001"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cinema?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		CloudinaryURL: getEnv("CLOUDINARY_URL", ""),
		VNPay: VNPayConfig{
			TmnCode:     getEnv("VNPAY_TMN_CODE", "STM5FOWI"),
			HashSecret:  getEnv("VNPAY_HASH_SECRET", "1QJBK7PAT7GB8KYTPKI45X0QGLEE0TGS"),
			Endpoint:    getEnv("VNPAY_ENDPOINT", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"),
			ReturnURL:   getEnv("VNPAY_RETURN_URL", "https://cinema-backend-yc14.onrender.com/api/payments/vnpay/return"),
			Environment: getEnv("VNPAY_ENVIRONMENT", "sandbox"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

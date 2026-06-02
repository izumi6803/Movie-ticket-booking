package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Test if the password hash matches
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqhmM6JGKpS4G3R1G2tJxP7QH0vP2"
	password := "admin123"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does NOT match:", err)

		// Generate a new hash for admin123
		newHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		fmt.Println("New hash for 'admin123':", string(newHash))
	} else {
		fmt.Println("Password matches!")
	}
}

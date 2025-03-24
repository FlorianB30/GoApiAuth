package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash the password before storing it in DB
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("❌ Error hashing password:", err)
	}
	return string(hashed), err
}

// Check if the provided password matches the stored bcrypt hash
func CheckPassword(hashedPwd, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	if err != nil {
		fmt.Println("❌ Password does NOT match:", err)
		return false
	}
	fmt.Println("✅ Password matches!")
	return true
}

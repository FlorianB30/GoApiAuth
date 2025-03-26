package controllers

import (
	"auth-api-go/config"
	"auth-api-go/models"
	"auth-api-go/utils"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
    "encoding/hex"
	"github.com/gin-gonic/gin"
)

// Register New User
func Register(c *gin.Context) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request to struct (only once)
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("❌ Error reading request JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	// Check if password is empty
	// if user.Password == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
	// 	return
	// }
	//check if name , email and password are empty
	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, email, and password are required"})
		return
	}

	// Vérifier si l'email est déjà utilisé
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already in use"})
		return
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user object for database
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	// Save user to database
	result := config.DB.Create(&newUser)
	if result.Error != nil {
		fmt.Println("Error saving user to database:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login User
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if password matches
	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT Token
	token, _ := utils.GenerateToken(user.Email)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Reset Password (with token)
func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	// Find user by reset token
	var user models.User
	if err := config.DB.Where("reset_token = ?", input.Token).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Check if token is expired
	if time.Now().After(user.ResetTokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	config.DB.Save(&user)

	c.JSON(http.StatusNoContent, nil)
}

// Forgot Password (Send Reset Token via Email)
func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Generate reset token
	resetToken := generateResetToken()
	user.ResetToken = resetToken
	user.ResetTokenExpiry = time.Now().Add(15 * time.Minute) // Token valid for 15 min

	// Save token in the database
	config.DB.Save(&user)

	// Send email
	err := sendResetEmail(user.Email, resetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent", "token": resetToken})
}

// Generate a random reset token
// func generateResetToken() string {
// 	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	rand.Seed(time.Now().UnixNano())
// 	b := make([]byte, 32)
// 	for i := range b {
// 		b[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(b)
// }
func generateResetToken() string {
    bytes := make([]byte, 8) // 8 bytes * 2 = 16 caractères hexadécimaux
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func sendResetEmail(email, token string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	to := []string{email}
	msg := []byte("To: " + email + "\r\n" +
		"Subject: Password Reset Request\r\n\r\n" +
		"Here is your password reset token:\r\n" +
		token + "\r\n\r\n" +
		"Or click the link below to reset your password:\r\n" +
		"http://localhost:8081/reset-password?token=" + token + "\r\n")

	fmt.Println("Sending email to:", email)
	fmt.Println("SMTP config:", smtpHost+":"+smtpPort, smtpUser)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, to, msg)
	if err != nil {
		fmt.Println("❌ Email sending failed:", err)
	}
	return err
}

// Logout
func Logout(c *gin.Context) {
	c.JSON(204, gin.H{
		"message": "Logged out successfully. Please delete the token on client side.",
	})
}

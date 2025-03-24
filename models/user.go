package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string    `json:"name"`
	Email            string    `json:"email" gorm:"unique"`
	Password         string    `json:"password"`
	ResetToken       string    `json:"reset_token"`
	ResetTokenExpiry time.Time `json:"reset_token_expiry"`
}

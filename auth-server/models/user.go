package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username      string    `gorm:"unique;not null"`
	Password      string    `gorm:"not null"`
	LastLoginTime time.Time `gorm:"default:null"`
}

func HashPassword(password string) (string, error) {

	hashed_password, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashed_password), nil
}

func CheckPassword(password string, hash string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

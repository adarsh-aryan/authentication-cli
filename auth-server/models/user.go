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

	// constraint:OnDelete:CASCADE tells the DB to delete the Account if this User is deleted
	Account *Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
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

type Account struct {
	gorm.Model

	IsLock               bool `gorm:"default:false"`
	MaxLoginAttempts     int  `gorm:"default:3"`
	CurrentLoginAttempts int

	// binding:"required" ensures your application validates its presence
	// gorm:"not null" forces the database column to reject empty/null values
	UserID uint `gorm:"uniqueIndex;not null" binding:"required"`

	User *User `gorm:"foreignKey:UserID;references:ID"`
}

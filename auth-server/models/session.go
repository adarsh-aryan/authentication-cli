package models

import (
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	SessionID      string `gorm:"uniqueIndex"`
	UserID         uint
	User           User
	ExpirationTime time.Time
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {

	expiration_time, err := strconv.Atoi(os.Getenv("SESSION_TIMEOUT_MINUTES"))
	if err != nil {
		expiration_time = 2 // default
	}
	s.ExpirationTime = time.Now().Add(time.Duration(expiration_time) * time.Minute)
	return nil
}

package models

import (
	"errors"
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username      string    `gorm:"unique;not null"`
	Password      string    `gorm:"not null"`
	LastLoginTime time.Time `gorm:"default:null"`
	TOTPSecret    string    `gorm:"column:totp_secret"`
	Is2FAEnabled  bool      `gorm:"column:is_2fa_enabled;default:false"`

	// constraint:OnDelete:CASCADE tells the DB to delete the Account if this User is deleted
	Account *Account `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func GetUser(username string, db *gorm.DB) (*User, error) {

	var user User
	err := db.Where("username=?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *User) CreateSession(db *gorm.DB) (*Session, error) {
	// save user session with timeout(configurable) (default 2minutes)
	user_session := Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}

	err := db.Create(&user_session).Error
	if err != nil {
		return nil, err
	}

	// update user last login time
	err = user.UpdateLastLogin(db)
	if err != nil {
		return nil, err
	}

	return &user_session, nil
}

func (user *User) IsSessionExpired(sessionId string, db *gorm.DB) (*Session, error) {

	// read the session created at this session id
	var session Session
	err := db.Preload("User").Where("session_id=?", sessionId).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no login session found")
		}
		return nil, err
	}

	// check for is session got expired
	if session.ExpirationTime.Before(time.Now()) {
		return nil, fmt.Errorf("Your session got expired, login again!")
	}

	return &session, nil
}

func (user *User) UpdateLastLogin(db *gorm.DB) error {

	// update the user last login time
	user.LastLoginTime = time.Now()
	err := db.Table("users").Where("id=?", user.ID).Update("last_login_time", user.LastLoginTime).Error
	if err != nil {
		return err
	}

	return nil
}

func (user *User) CheckPassword(password string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
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

func (user *User) SetTOTPSecret(secret string, db *gorm.DB) error {

	user.TOTPSecret = secret
	err := db.Model(&user).Update("totp_secret", user.TOTPSecret).Error

	if err != nil {
		return err
	}

	return nil
}

func (user *User) IsTOTPValid(code string) bool {

	// validate 6digit code input token against the saved secret ket
	valid := totp.Validate(code, user.TOTPSecret)
	return valid
}

func (user *User) Enable2FA(db *gorm.DB) error {

	if user.Is2FAEnabled {
		return nil
	}

	user.Is2FAEnabled = true
	err := db.Table("users").Where("id=?", user.ID).Update("is_2fa_enabled", user.Is2FAEnabled).Error

	if err != nil {
		return err
	}

	return nil
}

func (user *User) Disable2FA(db *gorm.DB) error {

	if !user.Is2FAEnabled {
		return nil
	}
	user.Is2FAEnabled = false
	err := db.Table("users").Where("id=?", user.ID).Update("is_2fa_enabled", user.Is2FAEnabled).Error

	if err != nil {
		return err
	}

	return nil
}

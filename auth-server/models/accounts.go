package models

import (
	"fmt"

	"gorm.io/gorm"
)

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

func (account *Account) ResetLoginAttempts(db *gorm.DB) error {

	if account.CurrentLoginAttempts == 0 {
		return nil
	}

	// reset to current login attempts to zero
	account.CurrentLoginAttempts = 0
	err := db.Table("accounts").Where("id=?", account.ID).Update("current_login_attempts", account.CurrentLoginAttempts).Error
	if err != nil {
		return err
	}

	return nil

}

func (account *Account) HandleLoginAttempts(db *gorm.DB, user *User) error {

	// handle login attempts after each wrong password or wrong totp code
	login_attempts := account.CurrentLoginAttempts
	login_attempts += 1

	account.CurrentLoginAttempts = login_attempts

	if login_attempts >= account.MaxLoginAttempts {
		account.IsLock = true // lock the account

		err := db.Table("accounts").Where("id=?", account.ID).Updates(map[string]interface{}{
			"is_lock":                account.IsLock,
			"current_login_attempts": account.CurrentLoginAttempts,
		}).Error
		if err != nil {
			return err
		}
		return fmt.Errorf("max login attempts reached!, %v account has been locked", user.Username)
	}

	err := db.Table("accounts").Where("id=?", account.ID).Update("current_login_attempts", account.CurrentLoginAttempts).Error
	if err != nil {
		return err
	}

	return nil

}

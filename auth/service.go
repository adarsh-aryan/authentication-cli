package auth

import (
	"errors"
	"fmt"
	"login-sys/auth/models"
	"login-sys/config"
	"login-sys/db"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Register(username string, password string) error {

	// register user

	// username and password validation
	if len(username) < 5 {
		return fmt.Errorf("Username %v is too short. It should be at five characters", username)
	}

	if len(password) < 5 {
		return fmt.Errorf("Password %v is too short. It could be bypass", password)
	}

	// check for is that username is already exist in the users table

	var exists bool
	err := db.DB.Select("1").Table("users").Where("username=?", username).Limit(1).Find(&exists).Error

	if err != nil {
		return err
	}

	// if username is already exists
	if exists {
		return fmt.Errorf("Username %v already exists", username)
	}

	// hash user password
	hashed_password, err := models.HashPassword(password)
	if err != nil {
		return err
	}
	// create user
	user := models.User{
		Username: username,
		Password: hashed_password,
	}
	db.DB.Create(&user)
	return nil
}

func Login(username string, password string) error {
	// login user

	// check for this username registration
	var user models.User
	err := db.DB.Where("username=?", username).First(&user).Error

	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("Username %v is not registered", username)
	}

	// verify password
	// hashed current password
	// compare current hashed password to user hashed password in the database

	err = models.CheckPassword(password, user.Password)
	if err != nil {
		return fmt.Errorf("password verification failed %v", password)
	}

	// save user session with timeout(configurable) (default 2minutes)
	user_session := models.Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}

	err = db.DB.Create(&user_session).Error
	if err != nil {
		return err
	}

	// save the current session in config file
	err = config.Save(user_session.SessionID, user_session.ExpirationTime)
	if err != nil {
		return err
	}

	// update the user last login time
	user.LastLoginTime = time.Now()
	err = db.DB.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func WhoAmI() (string, error) {
	// read session file
	config, err := config.Load()
	if err != nil {
		return "", err
	}

	// read the session created at this session id
	var session models.Session
	err = db.DB.Preload("User").Where("session_id=?", config.SessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("no login session found")
		}
		return "", err
	}

	// check for is session got expired
	if session.ExpirationTime.Before(time.Now()) {
		return "", fmt.Errorf("Session %v got expired", session.SessionID)
	}

	return session.User.Username, nil
	// return username
}

func LogOut() error {
	// delete the session history of user
	cfg, err := config.Load()

	if err != nil {
		return err
	}

	var sessionId string

	// delete the user session from server
	err = db.DB.Select("session_id").Table("sessions").Where("session_id=?", cfg.SessionID).Limit(1).Scan(&sessionId).Error
	if err != nil {
		return err
	}

	if sessionId == "" {
		return fmt.Errorf("no login session found")
	}

	// if we could find the session delete the session from the server
	err = db.DB.Where("session_id=?", sessionId).Delete(&models.Session{}).Error
	if err != nil {
		return err
	}

	// delete the session from the config file
	err = config.Delete()
	if err != nil {
		return err
	}
	return nil
}

package auth

import (
	"errors"
	"fmt"
	"login-sys/auth-server/models"
	"login-sys/shared"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func (s *AuthService) Register(args shared.RegisterArgs, reply *shared.AuthResponse) error {

	// register user

	// username and password validation
	username := args.Username
	password := args.Password

	if len(username) < 5 {
		return fmt.Errorf("Username %v is too short. It should be at five characters", username)
	}

	if len(password) < 5 {
		return fmt.Errorf("Password %v is too short. It could be bypass", password)
	}

	// check for is that username is already exist in the users table

	var exists bool
	err := s.DB.Select("1").Table("users").Where("username=?", username).Limit(1).Find(&exists).Error

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
	// create user with thier account
	user := models.User{
		Username: username,
		Password: hashed_password,
		Account:  &models.Account{},
	}
	err = s.DB.Create(&user).Error
	if err != nil {
		return err
	}

	// update auth response
	reply.SetMessage("Registeration successful!")
	return nil
}

func (s *AuthService) Login(args shared.LoginArgs, reply *shared.LoginResponse) error {
	// login user

	username := args.Username
	password := args.Password

	// check for this username registration
	var user models.User
	err := s.DB.Preload("Account").Where("username=?", username).First(&user).Error

	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("Username %v is not registered", username)
	}

	// check for if account has been locked for this user
	if user.Account != nil && user.Account.IsLock {
		return fmt.Errorf("User %v account has been locked!", username)
	}

	// verify password
	// hashed current password
	// compare current hashed password to user hashed password in the database

	err = models.CheckPassword(password, user.Password)
	if err != nil {
		// if password has not been mached
		// we have to count it as failed attempt
		// update the account info for this user
		login_attempts := user.Account.CurrentLoginAttempts
		login_attempts += 1

		user.Account.CurrentLoginAttempts = login_attempts

		if login_attempts >= user.Account.MaxLoginAttempts {
			user.Account.IsLock = true // lock the account
			err := s.DB.Save(user.Account).Error
			if err != nil {
				return err
			}
			return fmt.Errorf("max login attempts reached!, %v account has been locked", username)
		}

		err := s.DB.Save(user.Account).Error
		if err != nil {
			return err
		}

		return fmt.Errorf("wrong password!")
	}

	// save user session with timeout(configurable) (default 2minutes)
	user_session := models.Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}

	err = s.DB.Create(&user_session).Error
	if err != nil {
		return err
	}

	// update the user last login time
	user.LastLoginTime = time.Now()
	err = s.DB.Save(&user).Error
	if err != nil {
		return err
	}

	// update auth response
	reply.SetUserDetails(&user, &user_session)
	reply.SetSessionId(user_session.SessionID)
	reply.SetMessage("Login successful!")

	return nil
}

func (s *AuthService) WhoAmI(args shared.WhoAmIArgs, reply *shared.AuthResponse) error {

	session_id := args.SessionId

	// read the session created at this session id
	var session models.Session
	err := s.DB.Preload("User").Where("session_id=?", session_id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no login session found")
		}
		return err
	}

	// check for is session got expired
	if session.ExpirationTime.Before(time.Now()) {
		return fmt.Errorf("Your session got expired, login again!")
	}

	// update auth response
	reply.SetMessage(session.User.Username) // set username in response message

	return nil
	// return username
}

func (s *AuthService) LogOut(args shared.LogoutArgs, reply *shared.AuthResponse) error {

	var sessionId string

	// delete the user session from server
	err := s.DB.Select("session_id").Table("sessions").Where("session_id=?", args.SessionId).Limit(1).Scan(&sessionId).Error
	if err != nil {
		return err
	}

	if sessionId == "" {
		return fmt.Errorf("no login session found")
	}

	// if we could find the session delete the session from the server
	err = s.DB.Where("session_id=?", sessionId).Delete(&models.Session{}).Error
	if err != nil {
		return err
	}

	// update the auth response
	reply.SetMessage("Session destroy successfully.")

	return nil
}

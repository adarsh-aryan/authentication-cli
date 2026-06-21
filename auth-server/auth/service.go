package auth

import (
	"fmt"
	"log"
	"login-sys/auth-server/models"
	"login-sys/shared"

	"github.com/pquerna/otp/totp"
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
		return fmt.Errorf("Username %v is too short. It should be at least five characters", username)
	}

	if len(password) < 5 {
		return fmt.Errorf("Password %v is too short. It should be at least five characters", password)
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
		return fmt.Errorf("%v is not registered", username)
	}

	// check for if account has been locked for this user
	if user.Account != nil && user.Account.IsLock {
		return fmt.Errorf("%v account has been locked!", username)
	}

	// verify password
	// hashed current password
	// compare current hashed password to user hashed password in the database

	err = user.CheckPassword(password)
	if err != nil {
		// if password has not been mached
		// we have to count it as failed attempt
		// update the account info for this user
		err := user.Account.HandleLoginAttempts(s.DB, &user)
		if err != nil {
			return err
		}

		return fmt.Errorf("wrong password!, you have %v attempts left", user.Account.MaxLoginAttempts-user.Account.CurrentLoginAttempts)
	}

	if user.Is2FAEnabled {
		// check for is totp code valid
		valid := user.IsTOTPValid(args.OTP)
		if !valid {
			// if otp code is not valid
			// we have to count it as failed login attempt
			// update account info for this user
			err := user.Account.HandleLoginAttempts(s.DB, &user)
			if err != nil {
				return err
			}

			return fmt.Errorf("%s account has MFA enabled, you have %v attempts left", username, user.Account.MaxLoginAttempts-user.Account.CurrentLoginAttempts)
		}
	}

	// if user is successful login reset current login attempts to zero
	user.Account.ResetLoginAttempts(s.DB)

	// create the user session
	user_session, err := user.CreateSession(s.DB)
	if err != nil {
		return err
	}

	// update auth response
	reply.SetUserDetails(&user, user_session)
	reply.SetSessionId(user_session.SessionID)
	reply.SetMessage("Login successful!")

	return nil
}

func (s *AuthService) Request2FASetUp(args shared.SessionArgs, reply *shared.SetUp2FAResponse) error {

	user := &models.User{}
	session, err := user.IsSessionExpired(args.SessionId, s.DB)

	if err != nil {
		return err
	}

	// generate a new unique secret key for the user
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "auth-cli",
		AccountName: session.User.Username,
	})

	user = &session.User
	log.Println("set otp secert")

	// set the user totp secret
	err = user.SetTOTPSecret(key.Secret(), s.DB)

	if err != nil {
		return err
	}

	reply.Secret = key.Secret()
	reply.URL = key.URL()

	return nil
}

func (s *AuthService) Verify2FA(args shared.Verify2FArgs, reply *shared.AuthResponse) error {

	// get user from the session id
	user := &models.User{}
	session, err := user.IsSessionExpired(args.SessionId, s.DB)

	if err != nil {
		return err
	}

	user = &session.User

	// if it is get totp secret saved in user model
	valid := user.IsTOTPValid(args.Code)

	if !valid {
		if err := user.Disable2FA(s.DB); err != nil {
			return err
		}

		reply.SetMessage("Invalid 2FA code. Please try again")
		return nil
	}

	err = user.Enable2FA(s.DB)
	if err != nil {
		return err
	}

	reply.SetMessage("2FA enabled successfully")
	return nil
}

func (s *AuthService) Disable2FA(args shared.SessionArgs, reply *shared.AuthResponse) error {

	// get user from the session id
	user := &models.User{}
	session, err := user.IsSessionExpired(args.SessionId, s.DB)

	if err != nil {
		return err
	}

	user = &session.User
	// disable 2fa for user
	err = user.Disable2FA(s.DB)
	if err != nil {
		return err
	}

	reply.SetMessage("2FA disabled successfully")
	return nil
}

func (s *AuthService) WhoAmI(args shared.SessionArgs, reply *shared.LoginResponse) error {

	user := &models.User{}
	session, err := user.IsSessionExpired(args.SessionId, s.DB)

	if err != nil {
		return err
	}

	user = &session.User

	// update auth response
	reply.SetUserDetails(user, session)
	reply.SetSessionId(session.SessionID)

	return nil
}

func (s *AuthService) LogOut(args shared.SessionArgs, reply *shared.AuthResponse) error {

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

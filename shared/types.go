package shared

import (
	"login-sys/auth-server/models"
	"time"
)

type RegisterArgs struct {
	Username string
	Password string
}

type LoginArgs struct {
	Username string
	Password string
	OTP      string
}

type SessionArgs struct {
	SessionId string
}

type Verify2FArgs struct {
	SessionId string
	Code      string
}

type UserDetails struct {
	Username              string
	RegistrationDate      string
	MFAStatus             bool
	SessionExpirationTime string
	LastLoginTime         string
}

type LoginResponse struct {
	UserDetails UserDetails
	Message     string
	SessionId   string
}

func (lr *LoginResponse) SetUserDetails(user *models.User, session *models.Session) (UserDetails, error) {

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return UserDetails{}, err
	}

	user_details := UserDetails{
		Username:              user.Username,
		RegistrationDate:      user.CreatedAt.In(loc).Format(time.RFC1123),
		MFAStatus:             user.Is2FAEnabled,
		SessionExpirationTime: session.ExpirationTime.In(loc).Format(time.RFC1123),
		LastLoginTime:         user.LastLoginTime.In(loc).Format(time.RFC1123),
	}

	lr.UserDetails = user_details
	return user_details, nil
}

func (lr *LoginResponse) SetSessionId(sessionId string) {
	lr.SessionId = sessionId
}

func (lr *LoginResponse) SetMessage(message string) {
	lr.Message = message
}

func (lr *LoginResponse) GetMessage() string {
	return lr.Message
}

type AuthResponse struct {
	Message   string
	SessionId string
}

func (ar *AuthResponse) SetSessionId(sessionId string) {
	ar.SessionId = sessionId
}

func (ar *AuthResponse) SetMessage(message string) {
	ar.Message = message
}

func (ar *AuthResponse) SetUserDetails(user *models.User, session *models.Session) error {
	return nil
}

func (ar *AuthResponse) GetMessage() string {
	return ar.Message
}

// /  TOTP 2FA RESPONSE
type SetUp2FAResponse struct {
	Secret string
	URL    string
}

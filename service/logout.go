package service

import (
	"errors"

	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
)

type LogoutResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (s *AuthService) Logout(userID string, sessionID string) (*LogoutResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if sessionID == "" {
		return nil, errors.New("session ID is required")
	}

	err := dbhelper.DeleteSession(sessionID, userID)
	if err != nil {
		return nil, errors.New("failed to logout: " + err.Error())
	}

	return &LogoutResponse{
		Message: "Successfully logged out",
		Success: true,
	}, nil
}

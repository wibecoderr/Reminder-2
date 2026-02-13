package service

import (
	"errors"
	"net/http"
	"os/user"

	"github.com/go-playground/validator/v10"
	utils "github.com/wibecoderr/Reminder-2.git"
	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
	"github.com/wibecoderr/Reminder-2.git/model"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string     `json:"token" :"token"`
	User  model.User `json:"user" :"user"`
}
type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (auth *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}
	password, err := dbhelper.GetPasswordByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	userId, err := dbhelper.GetID(req.Email)
	if err != nil {
		return nil, err
	}
	sessionID, err := dbhelper.CreateSession(user.id)
	if err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.id, sessionID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

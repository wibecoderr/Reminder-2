package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	utils2 "github.com/wibecoderr/Reminder-2.git"

	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type RegisterRequest struct {
	Name       string    `json:"name" validate:"required,min=3,max=32"`
	Email      string    `json:"email" validate:"required,email"`
	Phone      string    `json:"phone" validate:"required,numeric,min=10,max=10"`
	Password   string    `json:"password" validate:"required,min=8,max=32"`
	Updated_At time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    int64  `json:"id"`
}

func (s *UserService) CreateUser(req RegisterRequest) (*CreateUserResponse, error) {
	if err := validator.New().Struct(req); err != nil {
		return nil, errors.New("validation error")
	}

	// Check if user exists
	exists, err := dbhelper.IsUserExist(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashP, err := utils2.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	ID, err := dbhelper.CreateUser(req.Name, req.Email, req.Phone, hashP, time.Now())
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		ID:    int64(ID),
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
	"github.com/wibecoderr/Reminder-2.git/model"
)

type ReminderService struct {
	scheduler *ReminderScheduler
}

func NewReminderService(scheduler *ReminderScheduler) *ReminderService {
	return &ReminderService{
		scheduler: scheduler,
	}
}

// CreateReminderRequest represents request to create a reminder
type CreateReminderRequest struct {
	Message   string    `json:"message" validate:"required,min=1,max=500"`
	PopUpTime time.Time `json:"pop_up_time" validate:"required"`
}

// CreateReminderResponse represents response after creating reminder
type CreateReminderResponse struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	PopUpTime time.Time `json:"pop_up_time"`
	UserID    int       `json:"user_id"`
}

// UpdateReminderRequest represents request to update a reminder
type UpdateReminderRequest struct {
	Message   string    `json:"message" validate:"required,min=1,max=500"`
	PopUpTime time.Time `json:"pop_up_time" validate:"required"`
}

func (s *ReminderService) CreateReminder(req CreateReminderRequest, userID int) (*CreateReminderResponse, error) {
	// Validate request
	if err := validator.New().Struct(req); err != nil {
		return nil, errors.New("validation error: " + err.Error())
	}

	// Check if pop up time is in the future
	if req.PopUpTime.Before(time.Now()) {
		return nil, errors.New("pop up time must be in the future")
	}

	// Create reminder in database
	id, err := dbhelper.CreateMessage(req.Message, req.PopUpTime, userID)
	if err != nil {
		return nil, errors.New("failed to create reminder: " + err.Error())
	}

	// Add to scheduler
	reminder := model.Reminder{
		ID:        id,
		Message:   req.Message,
		PopUpTime: req.PopUpTime,
		UserID:    userID,
	}
	s.scheduler.AddReminder(reminder)

	return &CreateReminderResponse{
		ID:        id,
		Message:   req.Message,
		PopUpTime: req.PopUpTime,
		UserID:    userID,
	}, nil
}

func (s *ReminderService) UpdateReminder(reminderID int, req UpdateReminderRequest, userID int) error {
	// Validate request
	if err := validator.New().Struct(req); err != nil {
		return errors.New("validation error: " + err.Error())
	}

	if req.PopUpTime.Before(time.Now()) {
		return errors.New("pop up time must be in the future")
	}

	exists, err := dbhelper.IsReminderOwnedByUser(reminderID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("reminder not found or access denied")
	}

	err = dbhelper.UpdateMessage(reminderID, req.Message, req.PopUpTime)
	if err != nil {
		return errors.New("failed to update reminder: " + err.Error())
	}

	// Update in scheduler (remove old, add new)
	reminder := model.Reminder{
		ID:        reminderID,
		Message:   req.Message,
		PopUpTime: req.PopUpTime,
		UserID:    userID,
	}
	s.scheduler.AddReminder(reminder)

	return nil
}

func (s *ReminderService) GetUserReminders(userID int) (model.Reminder1, error) {
	reminders, err := dbhelper.GetReminderByUserID(userID)
	if err != nil {
		return reminders, err
	}
	return reminders, nil
}

func (s *ReminderService) DeleteReminder(reminderID int, userID int) error {

	exists, err := dbhelper.IsReminderOwnedByUser(reminderID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("reminder not found or access denied")
	}

	err = dbhelper.Delete(reminderID, userID)
	if err != nil {
		return errors.New("failed to delete reminder: " + err.Error())
	}

	return nil
}

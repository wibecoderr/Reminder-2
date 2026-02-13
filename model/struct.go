package model

import "time"

type Reminder struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	PopUpTime time.Time `json:"pop_up_time" db:"pop_up_time"`
	UserID    int       `json:"user_id" db:"user_id"`
}

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	PhoneNo   string    `json:"phone_no" db:"phone_no"` // Changed to string
	Password  string    `json:"-" db:"password"`        // JSON "-" hides from responses
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Error struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
type UserCxt struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}
type Reminder1 struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	PopUpTime time.Time `json:"pop_up_time" db:"pop_up_time"`
	status    string    `json:"status" db:"status"`
}

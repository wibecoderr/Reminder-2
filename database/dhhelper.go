package dbhelper

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/Reminder-2.git"
	"github.com/wibecoderr/Reminder-2.git/model"
	"golang.org/x/crypto/openpgp/packet"
)

var DB *sqlx.DB

func CreateUser(Name, Email, Phone_no, Password string, updated_At time.Time) (int, error) {
	query := `INSERT INTO users (name, email, phone_no, password, updated_at)
              VALUES ($1, LOWER(TRIM($2)), $3, $4, $5)
              returning id`
	var id int
	err := DB.Get(&id, query, Name, Email, Phone_no, Password, updated_At)
	return id, err
}

func IsUserExist(email string) (bool, error) {
	query := `select count(*)>0 from users where email=$1`
	var exists bool
	err := DB.Get(&exists, query, email)
	return exists, err
}

func CreateMessage(message string, popUpTime time.Time, userID int) (int, error) {
	query := `INSERT INTO messages (message, pop_up_time, user_id)
              VALUES ($1, $2, $3)
              RETURNING id`

	var id int
	err := DB.Get(&id, query, message, popUpTime, userID)
	return id, err
}

func UpdateMessage(id int, message string, popUpTime time.Time) error {
	query := `UPDATE messages 
              SET message = $2, pop_up_time = $3
              WHERE id = $1`

	_, err := DB.Exec(query, id, message, popUpTime)
	if err != nil {
		return err
	}
	return err
}
func GetUserIDBySession(sessionID string) (string, error) {
	query := `SELECT user_id FROM sessions WHERE session_id = $1 AND expires_at > NOW()`
	var userID string
	err := DB.Get(&userID, query, sessionID)
	return userID, err
}
func GetPasswordByEmail(email string) (string, error) {
	aql := `SELECT password FROM users WHERE email=$1`
	var password string
	err := DB.Get(&password, aql, email)
	return password, err
}
func CreateSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	query := `INSERT INTO sessions (session_token, user_id, created_at, expires_at)
              VALUES ($1, $2, NOW(), NOW() + INTERVAL '24 hours')
              RETURNING session_token`

	var token string
	err := DB.Get(&token, query, sessionID, userID)
	return token, err
}

func GetID(email string) (string, error) {
	sql := `SELECT id FROM users WHERE email=$1`
	var id string
	err := DB.Get(&id, sql, email)
	return id, err
}

func DeleteSession(sessionID, userId string) error {
	query := `DELETE FROM sessions WHERE session_id=$1 and user_id=$2`
	_, err := DB.Exec(query, sessionID)
	return err
}
func IsReminderOwnedByUser(reminderID, userID int) (bool, error) {
	sql := `select count(*)> 0 from message where id=$1 and user_id=$2`

	var exists bool
	err := DB.Get(&exists, sql, reminderID, userID)
	return exists, err
}
func Delete(id, UserID int) error {
	sql := `DELETE FROM users WHERE id=$1 and user_id=$2`
	_, err := DB.Exec(sql, id, UserID)
	return err
}

func GetReminderByUserID(userID int) (model.Reminder1, error) {
	sql := `SELECT id , message , pop_up_time , status
  FROM users WHERE user_id=$1`
	var reminder model.Reminder1
	err := DB.Get(&reminder, sql, userID)
	return reminder, err

}

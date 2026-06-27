package dbHelper

import (
	"time"

	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

func IsUserExists(email string) (bool, error) {
	SQL := `SELECT count(id) > 0 as is_exist
			  FROM users
			  WHERE email = TRIM($1)
			    AND archived_at IS NULL`

	var exist bool
	err := database.Todo.Get(&exist, SQL, email)
	return exist, err
}

func CreateUser(name, email, password string) error {
	SQL := `INSERT INTO users(name,email,password)
			VALUES (TRIM($1),TRIM($2),$3)`
	_, err := database.Todo.Exec(SQL, name, email, password)
	return err
}

func CreateUserSession(userID string) (string, error) {
	var sessionID string
	SQL := `INSERT INTO user_session(user_id) 
              VALUES ($1) RETURNING id`
	err := database.Todo.Get(&sessionID, SQL, userID)
	return sessionID, err
}

func AuthenticateUser(email, password string) (string, error) {
	SQL := `SELECT id,
       			   password
			  FROM users 
			  WHERE email = TRIM($1)
			    AND archived_at IS NULL`

	var user models.LoginData
	if err := database.Todo.Get(&user, SQL, email); err != nil {
		return "", err
	}
	if passwordErr := utils.CheckPassword(password, user.PasswordHash); passwordErr != nil {
		return "", passwordErr
	}
	return user.ID, nil
}

func GetArchivedAt(sessionID string) (*time.Time, error) {
	var archivedAt *time.Time

	SQL := `SELECT archived_at 
              FROM user_session 
              WHERE id = $1`

	getErr := database.Todo.Get(&archivedAt, SQL, sessionID)
	return archivedAt, getErr
}

func GetUser(userID string) (models.User, error) {
	var user models.User
	SQL := `SELECT id, name, email 
              FROM users 
              WHERE id = $1
                AND archived_at IS NULL`

	err := database.Todo.Get(&user, SQL, userID)
	return user, err
}

func DeleteUserSession(sessionID string) error {
	SQL := `UPDATE user_session
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, sessionID)
	return err
}

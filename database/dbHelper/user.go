package dbHelper

import (
	"github.com/yourusername/my-project/database"
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

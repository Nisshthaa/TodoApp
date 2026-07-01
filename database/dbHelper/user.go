package dbHelper

import (
	"time"

	"github.com/jmoiron/sqlx"
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

func CreateUser(db sqlx.Ext, name, email, password string) (string, error) {

	SQL := `INSERT INTO users(name, email, password) VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id`
	var userID string
	if err := db.QueryRowx(SQL, name, email, password).Scan(&userID); err != nil {
		return "", err
	}

	return userID, nil
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {

	SQL := `INSERT INTO user_session(user_id,session_token) 
              VALUES ($1,$2) `
	_, err := db.Exec(SQL, userID, sessionToken)
	return err

}

func GetUserIDByPassword(email, password string) (string, error) {
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

func DeleteUserSessionTX(tx *sqlx.Tx, sessionID string) error {
	SQL := `UPDATE user_session
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, err := tx.Exec(SQL, sessionID)
	return err
}

func DeleteUser(tx *sqlx.Tx, userID string) error {
	SQL := `UPDATE users
			  SET archived_at = NOW()
			  WHERE id = $1
			    AND archived_at IS NULL`

	_, err := tx.Exec(SQL, userID)
	return err
}

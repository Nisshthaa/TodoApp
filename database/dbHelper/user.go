package dbHelper

import (
	"database/sql"

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

func GetUserBySession(sessionToken string) (*models.User, error) {
	SQL := `SELECT u.id, u.name, u.email, u.created_at 
			FROM users u
			JOIN user_session us on u.id = us.user_id
			WHERE u.archived_at IS NULL AND us.session_token = $1`

	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}

func DeleteUserSessionToken(db sqlx.Ext, userID, token string) error {

	SQL := `UPDATE user_session
			SET archived_at = NOW()
			WHERE user_id = $1
			  AND session_token = $2
			  AND archived_at IS NULL`
	_, err := db.Exec(SQL, userID, token)
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

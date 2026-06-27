package dbHelper

import (
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/models"
)

func IsTodoExists(userID, title string) (bool, error) {

	SQL := `SELECT COUNT(id) > 0
		FROM todos
		WHERE title = TRIM($1)
		  AND user_id=$2
		  AND archived_at IS NULL`
	var exists bool

	err := database.Todo.Get(&exists, SQL, title, userID)

	return exists, err
}

func CreateTodo(body models.TodoRequest) error {
	args := []interface{}{
		body.UserID,
		body.Title,
		body.Description,
	}

	SQL := `INSERT INTO todos(user_id,title,description)
	     VALUES ($1,TRIM($2),TRIM($3))`

	_, err := database.Todo.Exec(SQL, args...)
	return err
}

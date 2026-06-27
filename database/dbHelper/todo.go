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

func GetAllTodos(userID, keyword, completed string) ([]models.Todo, error) {
	SQL := `SELECT id, user_id, title, description, is_completed
				FROM todos
				WHERE user_id = $1
				  AND (
					$2 = '' OR (title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%')
					)
				  AND ($3 = '' OR is_completed = CAST($3 AS BOOLEAN))
				  AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	getErr := database.Todo.Select(&todos, SQL, userID, keyword, completed)
	return todos, getErr
}

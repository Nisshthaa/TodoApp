package dbHelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/models"
)

func IsTodoExists(userID, title string) (bool, error) {
	SQL := `SELECT COUNT(id) > 0
		FROM todos
		WHERE title = TRIM($1)
		  AND user_id = $2
		  AND archived_at IS NULL`

	var exists bool
	err := database.Todo.Get(&exists, SQL, title, userID)
	return exists, err
}

func CreateTodo(body models.TodoRequest, userID string) error {
	SQL := `INSERT INTO todos(user_id,title,description)
	     VALUES ($1,TRIM($2),TRIM($3))`

	_, err := database.Todo.Exec(SQL, userID, body.Title, body.Description)
	return err
}

func GetAllTodos(userID string) ([]models.Todo, error) {
	SQL := `SELECT id, user_id, title, description, is_completed
				FROM todos
				WHERE user_id = $1
				  AND archived_at IS NULL`

	todos := make([]models.Todo, 0)
	err := database.Todo.Select(&todos, SQL, userID)
	return todos, err
}

func GetTodoByID(userID, todoID string) (*models.Todo, error) {
	SQL := `SELECT id, user_id, title, description, is_completed
				FROM todos
				WHERE user_id = $1
				  AND id = $2
				  AND archived_at IS NULL`

	var todo models.Todo
	err := database.Todo.Get(&todo, SQL, userID, todoID)
	return &todo, err
}

func MarkStatusCompleted(todoID, userID string) error {
	SQL := `UPDATE todos	
              SET is_completed = true        
              WHERE id = $1                  
                AND user_id = $2             
                AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, todoID, userID)
	return err
}

func UpdateTodo(userID, todoID string, body models.TodoRequest) error {

	args := []interface{}{
		body.Title,
		body.Description,
		todoID,
		userID,
	}

	SQL := `
		UPDATE todos
		SET
			title = TRIM($1),
			description = TRIM($2),
			updated_at = NOW()
		WHERE
			id = $3
		  AND user_id=$4
			AND archived_at IS NULL
	`

	_, err := database.Todo.Exec(SQL, args...)

	return err
}

func DeleteTodoByID(userID, todoID string) error {
	SQL := `UPDATE todos
			  SET archived_at = NOW()        
			  WHERE id = $1                  
			    AND user_id = $2             
			    AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, todoID, userID)
	return err
}

func DeleteAllTodos(db sqlx.Ext, userID string) error {
	SQL := `UPDATE todos
              SET archived_at = NOW()        
              WHERE user_id = $1             
                AND archived_at IS NULL`

	_, err := db.Exec(SQL, userID)
	return err
}

package models

type TodoRequest struct {
	UserID      string `json:"user_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
}

type Todo struct {
	ID          string `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IsCompleted bool   `json:"isCompleted" db:"is_completed"`
	UserID      string `json:"userId" db:"user_id"`
}

package models

type TodoRequest struct {
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
}

type Todo struct {
	ID          string `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	IsCompleted bool   `json:"isCompleted" db:"is_completed"`
}

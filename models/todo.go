package models

type TodoRequest struct {
	UserID      string `json:"user_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
}

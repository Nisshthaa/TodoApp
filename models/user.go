package models

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"gte=6,lte=15"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,gte=6,lte=15"`
}

type LoginData struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password"`
}

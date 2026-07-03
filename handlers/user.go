package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/database/dbHelper"
	"github.com/yourusername/my-project/middlewares"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body models.RegisterRequest
	if parseErr := utils.ParseBody(r, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	exists, existErr := dbHelper.IsUserExists(body.Email)
	if existErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existErr, "failed to check user existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}

	hashedPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hashErr, "failed to secure password")
		return
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}

		sessionErr := dbHelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}
		return nil
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "User registered successfully",
		Token:   sessionToken,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest

	if parseErr := utils.ParseBody(r, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	v := validator.New()
	if err := v.Struct(body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	userID, userErr := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}

	if userID == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "user not found")
		return
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionErr := dbHelper.CreateUserSession(database.Todo, userID, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "User logged in successfully",
		Token:   sessionToken,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	utils.RespondJSON(w, http.StatusOK, userCtx)
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	userCtx := middlewares.UserContext(r)

	err := dbHelper.DeleteUserSessionToken(database.Todo, userCtx.UserID, token)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"account logout successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	token := r.Header.Get("x-api-key")

	txErr := database.Tx(func(tx *sqlx.Tx) error {

		err := dbHelper.DeleteAllTodos(tx, userID)
		if err != nil {
			return err
		}

		err = dbHelper.DeleteUser(tx, userID)
		if err != nil {
			return err
		}

		return dbHelper.DeleteUserSessionToken(tx, userID, token)
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to delete user account")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"account deleted successfully"})
}

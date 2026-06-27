package handlers

import (
	"net/http"

	"github.com/yourusername/my-project/database/dbHelper"
	"github.com/yourusername/my-project/middlewares"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.TodoRequest
	userCtx := middlewares.UserContext(r)
	todo.UserID = userCtx.UserID

	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	exists, err := dbHelper.IsTodoExists(todo.UserID, todo.Title)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to check todo existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo already exists")
		return
	}

	if err := dbHelper.CreateTodo(todo); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"Todo created successfully"})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	completed := r.URL.Query().Get("completed")

	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetAllTodos(userID, keyword, completed)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

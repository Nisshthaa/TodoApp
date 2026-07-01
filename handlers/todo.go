package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/yourusername/my-project/database/dbHelper"
	"github.com/yourusername/my-project/middlewares"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.TodoRequest
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if err := utils.ParseBody(r, &todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}

	v := validator.New()
	if err := v.Struct(todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	exists, err := dbHelper.IsTodoExists(userID, todo.Title)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to check todo existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo already exists")
		return
	}

	if err := dbHelper.CreateTodo(todo, userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"Todo created successfully"})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, err := dbHelper.GetAllTodos(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	todoID := chi.URLParam(r, "todoId")

	todos, err := dbHelper.GetTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func MarkCompleted(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if err := dbHelper.MarkTodoAsCompleted(todoID, userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to mark todo as completed")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo marked as completed successfully"})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	var body models.TodoRequest

	if err := utils.ParseBody(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Failed to parse request body")
		return
	}

	if err := dbHelper.UpdateTodo(userID, todoID, body); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Failed to update todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Todo updated successfully",
	})
}

func DeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	err := dbHelper.DeleteTodoByID(userID, todoID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo deleted successfully"})
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	err := dbHelper.DeleteAllTodos(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}

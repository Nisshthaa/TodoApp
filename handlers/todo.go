package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/yourusername/my-project/database"
	"github.com/yourusername/my-project/database/dbHelper"
	"github.com/yourusername/my-project/middlewares"
	"github.com/yourusername/my-project/models"
	"github.com/yourusername/my-project/utils"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.TodoRequest
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if parseErr := utils.ParseBody(r, &todo); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse body")
		return
	}

	v := validator.New()
	if err := v.Struct(todo); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	exists, existErr := dbHelper.IsTodoExists(userID, todo.Title)
	if existErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existErr, "failed to check todo existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "todo already exists")
		return
	}

	if createError := dbHelper.CreateTodo(todo, userID); createError != nil {
		utils.RespondError(w, http.StatusInternalServerError, createError, "failed to create todo")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"Todo created successfully"})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	todos, getErr := dbHelper.GetAllTodos(userID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todos)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID
	todoID := chi.URLParam(r, "todoId")

	todo, getErr := dbHelper.GetTodoByID(userID, todoID)
	if getErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, todo)
}

func MarkStatusCompleted(w http.ResponseWriter, r *http.Request) {
	todoID := chi.URLParam(r, "todoId")
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	if markErr := dbHelper.MarkStatusCompleted(todoID, userID); markErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, markErr, "failed to mark todo as completed")
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

	if parseErr := utils.ParseBody(r, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if updateErr := dbHelper.UpdateTodo(userID, todoID, body); updateErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, updateErr, "failed to update todo")
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

	deleteErr := dbHelper.DeleteTodoByID(userID, todoID)
	if deleteErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, deleteErr, "failed to delete todo")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"todo deleted successfully"})
}

func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	userID := userCtx.UserID

	deleteErr := dbHelper.DeleteAllTodos(database.Todo, userID)
	if deleteErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, deleteErr, "failed to delete todos")
		return
	}

	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"all todos deleted successfully"})
}

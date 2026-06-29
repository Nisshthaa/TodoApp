package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/my-project/handlers"
	"github.com/yourusername/my-project/middlewares"
)

type Server struct {
	chi.Router
	server *http.Server
}

func SetUpRoutes() *Server {
	router := chi.NewRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.Post("/register", handlers.RegisterUser)
		v1.Post("/login", handlers.LoginUser)

		v1.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate)

			r.Route("/user", func(user chi.Router) {
				user.Get("/me", handlers.GetUser)
				user.Post("/logout", handlers.LogoutUser)
			})

			r.Route("/todo", func(todo chi.Router) {
				todo.Post("/", handlers.CreateTodo)
				todo.Get("/", handlers.GetAllTodos)
				todo.Delete("/delete-all", handlers.DeleteAllTodos)

				todo.Route("/{todoId}", func(todoIDRoute chi.Router) {
					todoIDRoute.Put("/edit", handlers.UpdateTodo)
					todoIDRoute.Put("/mark-completed", handlers.MarkCompleted)
					todoIDRoute.Delete("/", handlers.DeleteTodoByID)
				})
			})

		})

	})
	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:    port,
		Handler: svc.Router,
	}
	return svc.server.ListenAndServe()
}

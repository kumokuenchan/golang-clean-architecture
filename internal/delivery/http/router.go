package http

import (
	"net/http"

	"clean-arch-sample/internal/usecase"
)

type Router struct {
	userHandler *UserHandler
}

func NewRouter(userUsecase *usecase.UserUsecase) *Router {
	return &Router{
		userHandler: NewUserHandler(userUsecase),
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.userHandler.GetAllUsers(w, req)
		case http.MethodPost:
			r.userHandler.CreateUser(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.userHandler.GetUser(w, req)
		case http.MethodPut:
			r.userHandler.UpdateUser(w, req)
		case http.MethodDelete:
			r.userHandler.DeleteUser(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
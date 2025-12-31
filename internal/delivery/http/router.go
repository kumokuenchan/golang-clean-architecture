package http

import (
	"net/http"

	"clean-arch-sample/internal/api"
)

type Router struct {
	server api.ServerInterface
}

func NewRouter(server api.ServerInterface) *Router {
	return &Router{
		server: server,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.server.GetAllUsers(w, req)
		case http.MethodPost:
			r.server.CreateUser(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			r.server.GetUser(w, req)
		case http.MethodPut:
			r.server.UpdateUser(w, req)
		case http.MethodDelete:
			r.server.DeleteUser(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clean-arch-sample/internal/usecase"
)

// ServerInterface defines the interface for the API server
type ServerInterface interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

// Server implements the ServerInterface
type Server struct {
	userUsecase *usecase.UserUsecase
}

// NewServer creates a new API server
func NewServer(userUsecase *usecase.UserUsecase) *Server {
	return &Server{
		userUsecase: userUsecase,
	}
}

// CreateUser handles the POST /users endpoint
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	data, err := s.userUsecase.CreateUser(req.Name, req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

// GetUser handles the GET /user endpoint
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "Invalid ID"})
		return
	}

	data, err := s.userUsecase.GetUser(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetAllUsers handles the GET /users endpoint
func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	data, err := s.userUsecase.GetAllUsers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// UpdateUser handles the PUT /user endpoint
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "Invalid ID"})
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	data, err := s.userUsecase.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// DeleteUser handles the DELETE /user endpoint
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MessageResponse{Message: "Invalid ID"})
		return
	}

	data, err := s.userUsecase.DeleteUser(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(MessageResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
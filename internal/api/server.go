package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clean-arch-sample/internal/domain"
	"clean-arch-sample/internal/dto"
	"clean-arch-sample/internal/mapper"
	"clean-arch-sample/internal/presenter"
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
	userRepo   domain.UserRepository
	presenter  presenter.UserPresenter
}

// NewServer creates a new API server
func NewServer(userRepo domain.UserRepository, presenter presenter.UserPresenter) *Server {
	return &Server{
		userRepo:  userRepo,
		presenter: presenter,
	}
}

// CreateUser handles the POST /users endpoint
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	if req.Name == "" || req.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "name and email are required"})
		return
	}

	// Use copier mapper to map request to domain
	domainUser := mapper.MapCreateRequestToDomain(&req)

	if err := s.userRepo.Create(domainUser); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	data, _ := s.presenter.PresentUser(domainUser)
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
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Invalid ID"})
		return
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	data, _ := s.presenter.PresentUser(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetAllUsers handles the GET /users endpoint
func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	data, _ := s.presenter.PresentUsers(users)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// UpdateUser handles the PUT /user endpoint
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Invalid ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	// Use copier mapper to map update request to domain
	domainUser := mapper.MapUpdateRequestToDomain(&req, user)

	if err := s.userRepo.Update(domainUser); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	data, _ := s.presenter.PresentUser(domainUser)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// DeleteUser handles the DELETE /user endpoint
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "ID is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: "Invalid ID"})
		return
	}

	_, err = s.userRepo.GetByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	if err := s.userRepo.Delete(id); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.MessageResponse{Message: err.Error()})
		return
	}

	data, _ := s.presenter.PresentSuccess("user deleted successfully")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
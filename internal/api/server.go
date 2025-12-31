package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"clean-arch-sample/internal/dto"
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
	userUsecase usecase.UserUsecaseInterface
}

// NewServer creates a new API server
func NewServer(userUsecase usecase.UserUsecaseInterface) *Server {
	return &Server{
		userUsecase: userUsecase,
	}
}

// writeJSONResponse writes a JSON response with proper headers
func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

// writeErrorResponse writes an error response
func (s *Server) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.MessageResponse{Message: message})
}

// extractAndValidateID extracts and validates ID from query parameters
func (s *Server) extractAndValidateID(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, ErrIDRequired
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, ErrInvalidID
	}

	if id <= 0 {
		return 0, ErrInvalidID
	}

	return id, nil
}

// decodeJSONBody decodes JSON request body
func (s *Server) decodeJSONBody(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}

// CreateUser handles the POST /users endpoint
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := s.decodeJSONBody(r, &req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := s.userUsecase.CreateUser(req.Name, req.Email)
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusCreated, data)
}

// GetUser handles the GET /user endpoint
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractAndValidateID(r)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == ErrInvalidID {
			statusCode = http.StatusBadRequest
		}
		s.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	data, err := s.userUsecase.GetUser(id)
	if err != nil {
		s.writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, data)
}

// GetAllUsers handles the GET /users endpoint
func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	data, err := s.userUsecase.GetAllUsers()
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, data)
}

// UpdateUser handles the PUT /user endpoint
func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractAndValidateID(r)
	if err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := s.decodeJSONBody(r, &req); err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := s.userUsecase.UpdateUser(id, req.Name, req.Email)
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, data)
}

// DeleteUser handles the DELETE /user endpoint
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractAndValidateID(r)
	if err != nil {
		s.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := s.userUsecase.DeleteUser(id)
	if err != nil {
		s.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSONResponse(w, http.StatusOK, data)
}
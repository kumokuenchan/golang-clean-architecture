package dto

// CreateUserRequest defines the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// UpdateUserRequest defines the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// UserResponse defines the response for user operations
type UserResponse struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john@example.com"`
	CreatedAt string `json:"created_at" example:"2023-01-01 12:00:00"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01 12:00:00"`
}

// UsersResponse defines the response for getting all users
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

// MessageResponse defines a generic response with a message
type MessageResponse struct {
	Message string `json:"message" example:"user deleted successfully"`
}
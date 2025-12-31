package mapper

import (
	"time"

	"clean-arch-sample/internal/domain"
	"clean-arch-sample/internal/dto"
	"github.com/jinzhu/copier"
)

// MapCreateRequestToDomain maps CreateUserRequest to domain User
func MapCreateRequestToDomain(req *dto.CreateUserRequest) *domain.User {
	user := &domain.User{}
	copier.Copy(user, req)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return user
}

// MapUpdateRequestToDomain maps UpdateUserRequest to domain User
func MapUpdateRequestToDomain(req *dto.UpdateUserRequest, existingUser *domain.User) *domain.User {
	// Only update fields that are not empty
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	existingUser.UpdatedAt = time.Now()
	return existingUser
}

// MapDomainToResponse maps domain User to UserResponse
func MapDomainToResponse(user *domain.User) *dto.UserResponse {
	response := &dto.UserResponse{}
	copier.Copy(response, user)
	response.CreatedAt = user.CreatedAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
	return response
}

// MapDomainListToResponse maps domain User slice to UsersResponse
func MapDomainListToResponse(users []*domain.User) *dto.UsersResponse {
	response := &dto.UsersResponse{
		Users: make([]dto.UserResponse, 0),
	}
	for _, user := range users {
		response.Users = append(response.Users, *MapDomainToResponse(user))
	}
	return response
}
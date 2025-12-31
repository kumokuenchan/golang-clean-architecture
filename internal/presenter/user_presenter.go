package presenter

import (
	"encoding/json"

	"clean-arch-sample/internal/domain"
	"clean-arch-sample/internal/dto"
	"clean-arch-sample/internal/mapper"
)

type UserPresenter interface {
	PresentUser(user *domain.User) ([]byte, error)
	PresentUsers(users []*domain.User) ([]byte, error)
	PresentError(message string) ([]byte, error)
	PresentSuccess(message string) ([]byte, error)
}

type HTTPUserPresenter struct{}

func NewHTTPUserPresenter() *HTTPUserPresenter {
	return &HTTPUserPresenter{}
}

func (p *HTTPUserPresenter) PresentUser(user *domain.User) ([]byte, error) {
	if user == nil {
		// Return empty user with zero values
		emptyUser := &dto.UserResponse{
			ID:        0,
			Name:      "",
			Email:     "",
			CreatedAt: "",
			UpdatedAt: "",
		}
		return json.Marshal(emptyUser)
	}
	// Use copier mapper to map domain to response
	response := mapper.MapDomainToResponse(user)
	return json.Marshal(response)
}

func (p *HTTPUserPresenter) PresentUsers(users []*domain.User) ([]byte, error) {
	if users == nil {
		// Return empty users response
		response := &dto.UsersResponse{
			Users: []dto.UserResponse{},
		}
		return json.Marshal(response)
	}
	// Use copier mapper to map domain list to response
	response := mapper.MapDomainListToResponse(users)
	return json.Marshal(response)
}

func (p *HTTPUserPresenter) PresentError(message string) ([]byte, error) {
	return json.Marshal(dto.MessageResponse{Message: message})
}

func (p *HTTPUserPresenter) PresentSuccess(message string) ([]byte, error) {
	return json.Marshal(dto.MessageResponse{Message: message})
}
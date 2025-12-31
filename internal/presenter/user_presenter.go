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
	// Use copier mapper to map domain to response
	response := mapper.MapDomainToResponse(user)
	return json.Marshal(response)
}

func (p *HTTPUserPresenter) PresentUsers(users []*domain.User) ([]byte, error) {
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
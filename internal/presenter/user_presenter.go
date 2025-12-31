package presenter

import (
	"encoding/json"

	"clean-arch-sample/internal/domain"
)

type UserOutputData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UsersOutputData struct {
	Users []UserOutputData `json:"users"`
}

type MessageOutputData struct {
	Message string `json:"message"`
}

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
	output := UserOutputData{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(output)
}

func (p *HTTPUserPresenter) PresentUsers(users []*domain.User) ([]byte, error) {
	var output []UserOutputData
	for _, user := range users {
		output = append(output, UserOutputData{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return json.Marshal(UsersOutputData{Users: output})
}

func (p *HTTPUserPresenter) PresentError(message string) ([]byte, error) {
	return json.Marshal(MessageOutputData{Message: message})
}

func (p *HTTPUserPresenter) PresentSuccess(message string) ([]byte, error) {
	return json.Marshal(MessageOutputData{Message: message})
}
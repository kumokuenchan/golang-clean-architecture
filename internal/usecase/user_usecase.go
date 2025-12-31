package usecase

import (
	"time"

	"github.com/kuen/clean-arch-sample/internal/domain"
	"github.com/kuen/clean-arch-sample/internal/presenter"
)

type UserUsecase struct {
	userRepo  domain.UserRepository
	presenter presenter.UserPresenter
}

func NewUserUsecase(userRepo domain.UserRepository, presenter presenter.UserPresenter) *UserUsecase {
	return &UserUsecase{
		userRepo:  userRepo,
		presenter: presenter,
	}
}

func (u *UserUsecase) CreateUser(name, email string) ([]byte, error) {
	if name == "" || email == "" {
		return u.presenter.PresentError("name and email are required")
	}

	user := &domain.User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUser(user)
}

func (u *UserUsecase) GetUser(id int) ([]byte, error) {
	if id <= 0 {
		return u.presenter.PresentError("invalid user ID")
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUser(user)
}

func (u *UserUsecase) GetAllUsers() ([]byte, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUsers(users)
}

func (u *UserUsecase) UpdateUser(id int, name, email string) ([]byte, error) {
	if id <= 0 {
		return u.presenter.PresentError("invalid user ID")
	}

	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return u.presenter.PresentError(err.Error())
	}

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(user); err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUser(user)
}

func (u *UserUsecase) DeleteUser(id int) ([]byte, error) {
	if id <= 0 {
		return u.presenter.PresentError("invalid user ID")
	}

	_, err := u.userRepo.GetByID(id)
	if err != nil {
		return u.presenter.PresentError(err.Error())
	}

	if err := u.userRepo.Delete(id); err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentSuccess("user deleted successfully")
}
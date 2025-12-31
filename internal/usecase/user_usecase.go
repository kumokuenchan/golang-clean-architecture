package usecase

import (
	"clean-arch-sample/internal/domain"
	"clean-arch-sample/internal/dto"
	"clean-arch-sample/internal/mapper"
	"clean-arch-sample/internal/presenter"
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

	// Use copier mapper to map request to domain
	req := &dto.CreateUserRequest{Name: name, Email: email}
	domainUser := mapper.MapCreateRequestToDomain(req)

	if err := u.userRepo.Create(domainUser); err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUser(domainUser)
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

	// Use copier mapper to map update request to domain
	req := &dto.UpdateUserRequest{Name: name, Email: email}
	domainUser := mapper.MapUpdateRequestToDomain(req, user)

	if err := u.userRepo.Update(domainUser); err != nil {
		return u.presenter.PresentError(err.Error())
	}

	return u.presenter.PresentUser(domainUser)
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
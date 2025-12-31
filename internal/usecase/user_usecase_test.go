package usecase

import (
	"testing"

	"clean-arch-sample/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPresenter implements UserPresenter for testing
type MockPresenter struct {
	presentUserCall   *domain.User
	presentUsersCall  []*domain.User
	presentErrorCall  string
	presentSuccessCall string
	presentUserErr    error
	presentUsersErr   error
	presentErrorErr   error
	presentSuccessErr error
}

func NewMockPresenter() *MockPresenter {
	return &MockPresenter{}
}

func (m *MockPresenter) PresentUser(user *domain.User) ([]byte, error) {
	m.presentUserCall = user
	if m.presentUserErr != nil {
		return nil, m.presentUserErr
	}
	return []byte(`{"id":1,"name":"test","email":"test@example.com"}`), nil
}

func (m *MockPresenter) PresentUsers(users []*domain.User) ([]byte, error) {
	m.presentUsersCall = users
	if m.presentUsersErr != nil {
		return nil, m.presentUsersErr
	}
	return []byte(`{"users":[{"id":1,"name":"test","email":"test@example.com"}]}`), nil
}

func (m *MockPresenter) PresentError(message string) ([]byte, error) {
	m.presentErrorCall = message
	if m.presentErrorErr != nil {
		return nil, m.presentErrorErr
	}
	return []byte(`{"message":"` + message + `"}`), nil
}

func (m *MockPresenter) PresentSuccess(message string) ([]byte, error) {
	m.presentSuccessCall = message
	if m.presentSuccessErr != nil {
		return nil, m.presentSuccessErr
	}
	return []byte(`{"message":"` + message + `"}`), nil
}

func (m *MockPresenter) SetPresentUserError(err error) {
	m.presentUserErr = err
}

func (m *MockPresenter) SetPresentUsersError(err error) {
	m.presentUsersErr = err
}

func (m *MockPresenter) SetPresentErrorError(err error) {
	m.presentErrorErr = err
}

func (m *MockPresenter) SetPresentSuccessError(err error) {
	m.presentSuccessErr = err
}

func (m *MockPresenter) WasPresentUserCalledWith(user *domain.User) bool {
	if m.presentUserCall == nil || user == nil {
		return false
	}
	return m.presentUserCall.ID == user.ID && m.presentUserCall.Name == user.Name && m.presentUserCall.Email == user.Email
}

func (m *MockPresenter) WasPresentUsersCalledWith(users []*domain.User) bool {
	if len(m.presentUsersCall) != len(users) {
		return false
	}
	for i, user := range users {
		if m.presentUsersCall[i].ID != user.ID || m.presentUsersCall[i].Name != user.Name || m.presentUsersCall[i].Email != user.Email {
			return false
		}
	}
	return true
}

func (m *MockPresenter) WasPresentErrorCalledWith(message string) bool {
	return m.presentErrorCall == message
}

func (m *MockPresenter) WasPresentSuccessCalledWith(message string) bool {
	return m.presentSuccessCall == message
}

func TestNewUserUsecase(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()

	usecase := NewUserUsecase(mockRepo, mockPresenter)

	assert.NotNil(t, usecase)
	assert.Equal(t, mockRepo, usecase.userRepo)
	assert.Equal(t, mockPresenter, usecase.presenter)
}

func TestUserUsecase_CreateUser(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.CreateUser("John Doe", "john@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasCreateCalledWith(&domain.User{Name: "John Doe", Email: "john@example.com"}))
	assert.True(t, mockPresenter.WasPresentUserCalledWith(&domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}))
}

func TestUserUsecase_CreateUserWithEmptyName(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.CreateUser("", "john@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("name and email are required"))
}

func TestUserUsecase_CreateUserWithEmptyEmail(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.CreateUser("John Doe", "")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("name and email are required"))
}

func TestUserUsecase_CreateUserWithRepositoryError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockRepo.SetError(domain.ErrUserNotFound)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.CreateUser("John Doe", "john@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("user not found"))
}

func TestUserUsecase_CreateUserWithPresenterError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	mockPresenter.SetPresentUserError(assert.AnError)
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.CreateUser("John Doe", "john@example.com")
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestUserUsecase_GetUser(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetUser(1)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetByIDCalledWith(1))
	assert.True(t, mockPresenter.WasPresentUserCalledWith(user))
}

func TestUserUsecase_GetUserWithInvalidID(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetUser(0)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("invalid user ID"))
}

func TestUserUsecase_GetUserWithNegativeID(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetUser(-1)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("invalid user ID"))
}

func TestUserUsecase_GetUserWithRepositoryError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockRepo.SetError(domain.ErrUserNotFound)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetUser(1)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetByIDCalledWith(1))
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("user not found"))
}

func TestUserUsecase_GetUserWithPresenterError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	mockPresenter.SetPresentUserError(assert.AnError)
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetUser(1)
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user1 := &domain.User{ID: 1, Name: "User 1", Email: "user1@example.com"}
	user2 := &domain.User{ID: 2, Name: "User 2", Email: "user2@example.com"}
	mockRepo.Create(user1)
	mockRepo.Create(user2)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetAllUsers()
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetAllCalled())
	assert.True(t, mockPresenter.WasPresentUsersCalledWith([]*domain.User{user1, user2}))
}

func TestUserUsecase_GetAllUsersWithRepositoryError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockRepo.SetError(domain.ErrUserNotFound)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetAllUsers()
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetAllCalled())
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("user not found"))
}

func TestUserUsecase_GetAllUsersWithPresenterError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "User 1", Email: "user1@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	mockPresenter.SetPresentUsersError(assert.AnError)
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.GetAllUsers()
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.UpdateUser(1, "Jane Doe", "jane@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetByIDCalledWith(1))
	assert.True(t, mockRepo.WasUpdateCalledWith(&domain.User{ID: 1, Name: "Jane Doe", Email: "jane@example.com"}))
	assert.True(t, mockPresenter.WasPresentUserCalledWith(&domain.User{ID: 1, Name: "Jane Doe", Email: "jane@example.com"}))
}

func TestUserUsecase_UpdateUserWithInvalidID(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.UpdateUser(0, "Jane Doe", "jane@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("invalid user ID"))
}

func TestUserUsecase_UpdateUserWithRepositoryError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockRepo.SetError(domain.ErrUserNotFound)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.UpdateUser(1, "Jane Doe", "jane@example.com")
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("user not found"))
}

func TestUserUsecase_UpdateUserWithPresenterError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	mockPresenter.SetPresentUserError(assert.AnError)
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.UpdateUser(1, "Jane Doe", "jane@example.com")
	assert.Error(t, err)
	assert.Nil(t, data)
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.DeleteUser(1)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockRepo.WasGetByIDCalledWith(1))
	assert.True(t, mockRepo.WasDeleteCalledWith(1))
	assert.True(t, mockPresenter.WasPresentSuccessCalledWith("user deleted successfully"))
}

func TestUserUsecase_DeleteUserWithInvalidID(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.DeleteUser(0)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("invalid user ID"))
}

func TestUserUsecase_DeleteUserWithRepositoryError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	mockRepo.SetError(domain.ErrUserNotFound)
	mockPresenter := NewMockPresenter()
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.DeleteUser(1)
	require.NoError(t, err)
	require.NotNil(t, data)
	assert.True(t, mockPresenter.WasPresentErrorCalledWith("user not found"))
}

func TestUserUsecase_DeleteUserWithPresenterError(t *testing.T) {
	mockRepo := domain.NewMockUserRepository()
	user := &domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.Create(user)
	mockPresenter := NewMockPresenter()
	mockPresenter.SetPresentSuccessError(assert.AnError)
	usecase := NewUserUsecase(mockRepo, mockPresenter)

	data, err := usecase.DeleteUser(1)
	assert.Error(t, err)
	assert.Nil(t, data)
}
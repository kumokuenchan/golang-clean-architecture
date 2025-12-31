package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)



// Test User entity
func TestUser_NewUser(t *testing.T) {
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Zero(t, user.ID)
	assert.True(t, user.CreatedAt.IsZero())
	assert.True(t, user.UpdatedAt.IsZero())
}

func TestUser_WithTimestamps(t *testing.T) {
	now := time.Now()
	user := &User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

// Test MockUserRepository
func TestMockUserRepository_Create(t *testing.T) {
	repo := NewMockUserRepository()
	user := &User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	err := repo.Create(user)
	require.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.True(t, repo.WasCreateCalledWith(&User{Name: "John Doe", Email: "john@example.com"}))
}

func TestMockUserRepository_CreateWithError(t *testing.T) {
	repo := NewMockUserRepository()
	repo.SetError(ErrUserNotFound)
	user := &User{Name: "John Doe", Email: "john@example.com"}

	err := repo.Create(user)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestMockUserRepository_GetByID(t *testing.T) {
	repo := NewMockUserRepository()
	user := &User{Name: "John Doe", Email: "john@example.com"}
	repo.Create(user)

	foundUser, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Name, foundUser.Name)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.True(t, repo.WasGetByIDCalledWith(1))
}

func TestMockUserRepository_GetByIDNotFound(t *testing.T) {
	repo := NewMockUserRepository()

	user, err := repo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestMockUserRepository_GetAll(t *testing.T) {
	repo := NewMockUserRepository()
	user1 := &User{Name: "User 1", Email: "user1@example.com"}
	user2 := &User{Name: "User 2", Email: "user2@example.com"}
	repo.Create(user1)
	repo.Create(user2)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.True(t, repo.WasGetAllCalled())
}

func TestMockUserRepository_Update(t *testing.T) {
	repo := NewMockUserRepository()
	user := &User{Name: "John Doe", Email: "john@example.com"}
	repo.Create(user)

	user.Name = "Jane Doe"
	err := repo.Update(user)
	require.NoError(t, err)
	assert.True(t, repo.WasUpdateCalledWith(user))
}

func TestMockUserRepository_UpdateNotFound(t *testing.T) {
	repo := NewMockUserRepository()
	user := &User{ID: 999, Name: "John Doe", Email: "john@example.com"}

	err := repo.Update(user)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestMockUserRepository_Delete(t *testing.T) {
	repo := NewMockUserRepository()
	user := &User{Name: "John Doe", Email: "john@example.com"}
	repo.Create(user)

	err := repo.Delete(1)
	require.NoError(t, err)
	assert.True(t, repo.WasDeleteCalledWith(1))

	// Verify user is deleted
	_, err = repo.GetByID(1)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}

func TestMockUserRepository_DeleteNotFound(t *testing.T) {
	repo := NewMockUserRepository()

	err := repo.Delete(999)
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}
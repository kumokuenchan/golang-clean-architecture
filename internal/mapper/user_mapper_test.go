package mapper

import (
	"testing"
	"time"

	"clean-arch-sample/internal/domain"
	"clean-arch-sample/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapCreateRequestToDomain(t *testing.T) {
	req := &dto.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	user := MapCreateRequestToDomain(req)

	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestMapCreateRequestToDomainWithEmptyRequest(t *testing.T) {
	req := &dto.CreateUserRequest{}

	user := MapCreateRequestToDomain(req)

	assert.Empty(t, user.Name)
	assert.Empty(t, user.Email)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
}

func TestMapUpdateRequestToDomain(t *testing.T) {
	now := time.Now()
	existingUser := &domain.User{
		ID:        1,
		Name:      "Old Name",
		Email:     "old@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	req := &dto.UpdateUserRequest{
		Name:  "New Name",
		Email: "new@example.com",
	}

	updatedUser := MapUpdateRequestToDomain(req, existingUser)

	assert.Equal(t, 1, updatedUser.ID)
	assert.Equal(t, "New Name", updatedUser.Name)
	assert.Equal(t, "new@example.com", updatedUser.Email)
	assert.Equal(t, now, updatedUser.CreatedAt) // Should not change
	assert.True(t, updatedUser.UpdatedAt.After(now)) // Should be updated
}

func TestMapUpdateRequestToDomainWithPartialUpdate(t *testing.T) {
	now := time.Now()
	existingUser := &domain.User{
		ID:        1,
		Name:      "Original Name",
		Email:     "original@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	req := &dto.UpdateUserRequest{
		Name:  "Updated Name Only",
		Email: "", // Empty email should not overwrite
	}

	updatedUser := MapUpdateRequestToDomain(req, existingUser)

	assert.Equal(t, "Updated Name Only", updatedUser.Name)
	assert.Equal(t, "original@example.com", updatedUser.Email) // Should remain unchanged
	assert.True(t, updatedUser.UpdatedAt.After(now))
}

func TestMapDomainToResponse(t *testing.T) {
	now := time.Now()
	user := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := MapDomainToResponse(user)

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), response.CreatedAt)
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), response.UpdatedAt)
}

func TestMapDomainToResponseWithDifferentTimestamps(t *testing.T) {
	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 45, 0, time.UTC)
	user := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	response := MapDomainToResponse(user)

	assert.Equal(t, "2023-01-01 12:00:00", response.CreatedAt)
	assert.Equal(t, "2023-01-02 15:30:45", response.UpdatedAt)
}

func TestMapDomainListToResponse(t *testing.T) {
	now := time.Now()
	users := []*domain.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Name:      "User 2",
			Email:     "user2@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	response := MapDomainListToResponse(users)

	require.NotNil(t, response)
	assert.Len(t, response.Users, 2)

	// Check first user
	assert.Equal(t, 1, response.Users[0].ID)
	assert.Equal(t, "User 1", response.Users[0].Name)
	assert.Equal(t, "user1@example.com", response.Users[0].Email)

	// Check second user
	assert.Equal(t, 2, response.Users[1].ID)
	assert.Equal(t, "User 2", response.Users[1].Name)
	assert.Equal(t, "user2@example.com", response.Users[1].Email)
}

func TestMapDomainListToResponseWithEmptyList(t *testing.T) {
	users := []*domain.User{}

	response := MapDomainListToResponse(users)

	require.NotNil(t, response)
	assert.Len(t, response.Users, 0)
}

func TestMapDomainListToResponseWithNilList(t *testing.T) {
	var users []*domain.User = nil

	response := MapDomainListToResponse(users)

	require.NotNil(t, response)
	assert.Len(t, response.Users, 0)
}

func TestMapDomainListToResponsePreservesOrder(t *testing.T) {
	now := time.Now()
	users := []*domain.User{
		{ID: 3, Name: "User 3", Email: "user3@example.com", CreatedAt: now, UpdatedAt: now},
		{ID: 1, Name: "User 1", Email: "user1@example.com", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "User 2", Email: "user2@example.com", CreatedAt: now, UpdatedAt: now},
	}

	response := MapDomainListToResponse(users)

	require.Len(t, response.Users, 3)
	assert.Equal(t, 3, response.Users[0].ID)
	assert.Equal(t, 1, response.Users[1].ID)
	assert.Equal(t, 2, response.Users[2].ID)
}
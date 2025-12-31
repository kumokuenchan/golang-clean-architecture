package presenter

import (
	"encoding/json"
	"testing"
	"time"

	"clean-arch-sample/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPUserPresenter_PresentUser(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	now := time.Now()
	user := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := presenter.PresentUser(user)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, float64(1), response["id"])
	assert.Equal(t, "John Doe", response["name"])
	assert.Equal(t, "john@example.com", response["email"])
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), response["created_at"])
	assert.Equal(t, now.Format("2006-01-02 15:04:05"), response["updated_at"])
}

func TestHTTPUserPresenter_PresentUserWithNilUser(t *testing.T) {
	presenter := NewHTTPUserPresenter()

	data, err := presenter.PresentUser(nil)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, float64(0), response["id"])
	assert.Empty(t, response["name"])
	assert.Empty(t, response["email"])
}

func TestHTTPUserPresenter_PresentUsers(t *testing.T) {
	presenter := NewHTTPUserPresenter()
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

	data, err := presenter.PresentUsers(users)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "users")
	usersArray, ok := response["users"].([]interface{})
	require.True(t, ok)
	assert.Len(t, usersArray, 2)

	// Check first user
	firstUser := usersArray[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstUser["id"])
	assert.Equal(t, "User 1", firstUser["name"])
	assert.Equal(t, "user1@example.com", firstUser["email"])

	// Check second user
	secondUser := usersArray[1].(map[string]interface{})
	assert.Equal(t, float64(2), secondUser["id"])
	assert.Equal(t, "User 2", secondUser["name"])
	assert.Equal(t, "user2@example.com", secondUser["email"])
}

func TestHTTPUserPresenter_PresentUsersWithEmptyList(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	users := []*domain.User{}

	data, err := presenter.PresentUsers(users)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "users")
	usersArray := response["users"].([]interface{})
	assert.Len(t, usersArray, 0)
}

func TestHTTPUserPresenter_PresentUsersWithNilList(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	var users []*domain.User = nil

	data, err := presenter.PresentUsers(users)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "users")
	usersArray := response["users"].([]interface{})
	assert.Len(t, usersArray, 0)
}

func TestHTTPUserPresenter_PresentError(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	errorMessage := "Something went wrong"

	data, err := presenter.PresentError(errorMessage)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, errorMessage, response["message"])
}

func TestHTTPUserPresenter_PresentErrorWithEmptyMessage(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	errorMessage := ""

	data, err := presenter.PresentError(errorMessage)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, "", response["message"])
}

func TestHTTPUserPresenter_PresentSuccess(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	successMessage := "Operation completed successfully"

	data, err := presenter.PresentSuccess(successMessage)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, successMessage, response["message"])
}

func TestHTTPUserPresenter_PresentSuccessWithEmptyMessage(t *testing.T) {
	presenter := NewHTTPUserPresenter()
	successMessage := ""

	data, err := presenter.PresentSuccess(successMessage)
	require.NoError(t, err)
	require.NotNil(t, data)

	var response map[string]interface{}
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	assert.Equal(t, "", response["message"])
}

func TestHTTPUserPresenter_JSONValidity(t *testing.T) {
	presenter := NewHTTPUserPresenter()

	// Test that all presenter methods produce valid JSON
	testCases := []struct {
		name     string
		method   func() ([]byte, error)
	}{
		{"PresentUser", func() ([]byte, error) {
			user := &domain.User{ID: 1, Name: "Test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
			return presenter.PresentUser(user)
		}},
		{"PresentUsers", func() ([]byte, error) {
			user := &domain.User{ID: 1, Name: "Test", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}
			return presenter.PresentUsers([]*domain.User{user})
		}},
		{"PresentError", func() ([]byte, error) {
			return presenter.PresentError("test error")
		}},
		{"PresentSuccess", func() ([]byte, error) {
			return presenter.PresentSuccess("test success")
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.method()
			require.NoError(t, err)
			require.NotNil(t, data)

			var jsonData interface{}
			err = json.Unmarshal(data, &jsonData)
			require.NoError(t, err, "Should produce valid JSON")
		})
	}
}
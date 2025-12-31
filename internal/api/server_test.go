package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clean-arch-sample/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUserUsecase implements UserUsecase for testing
type MockUserUsecase struct {
	createUserCall   struct {
		name  string
		email string
	}
	getUserCall     int
	getAllUsersCall bool
	updateUserCall struct {
		id    int
		name  string
		email string
	}
	deleteUserCall int

	createUserReturn   []byte
	createUserError    error
	getUserReturn     []byte
	getUserError       error
	getAllUsersReturn  []byte
	getAllUsersError   error
	updateUserReturn   []byte
	updateUserError    error
	deleteUserReturn   []byte
	deleteUserError    error
}

func NewMockUserUsecase() *MockUserUsecase {
	return &MockUserUsecase{}
}

func (m *MockUserUsecase) CreateUser(name, email string) ([]byte, error) {
	m.createUserCall.name = name
	m.createUserCall.email = email
	if m.createUserError != nil {
		return nil, m.createUserError
	}
	return m.createUserReturn, nil
}

func (m *MockUserUsecase) GetUser(id int) ([]byte, error) {
	m.getUserCall = id
	if m.getUserError != nil {
		return nil, m.getUserError
	}
	return m.getUserReturn, nil
}

func (m *MockUserUsecase) GetAllUsers() ([]byte, error) {
	m.getAllUsersCall = true
	if m.getAllUsersError != nil {
		return nil, m.getAllUsersError
	}
	return m.getAllUsersReturn, nil
}

func (m *MockUserUsecase) UpdateUser(id int, name, email string) ([]byte, error) {
	m.updateUserCall.id = id
	m.updateUserCall.name = name
	m.updateUserCall.email = email
	if m.updateUserError != nil {
		return nil, m.updateUserError
	}
	return m.updateUserReturn, nil
}

func (m *MockUserUsecase) DeleteUser(id int) ([]byte, error) {
	m.deleteUserCall = id
	if m.deleteUserError != nil {
		return nil, m.deleteUserError
	}
	return m.deleteUserReturn, nil
}

func (m *MockUserUsecase) SetCreateUserReturn(data []byte, err error) {
	m.createUserReturn = data
	m.createUserError = err
}

func (m *MockUserUsecase) SetGetUserReturn(data []byte, err error) {
	m.getUserReturn = data
	m.getUserError = err
}

func (m *MockUserUsecase) SetGetAllUsersReturn(data []byte, err error) {
	m.getAllUsersReturn = data
	m.getAllUsersError = err
}

func (m *MockUserUsecase) SetUpdateUserReturn(data []byte, err error) {
	m.updateUserReturn = data
	m.updateUserError = err
}

func (m *MockUserUsecase) SetDeleteUserReturn(data []byte, err error) {
	m.deleteUserReturn = data
	m.deleteUserError = err
}

func (m *MockUserUsecase) WasCreateUserCalledWith(name, email string) bool {
	return m.createUserCall.name == name && m.createUserCall.email == email
}

func (m *MockUserUsecase) WasGetUserCalledWith(id int) bool {
	return m.getUserCall == id
}

func (m *MockUserUsecase) WasGetAllUsersCalled() bool {
	return m.getAllUsersCall
}

func (m *MockUserUsecase) WasUpdateUserCalledWith(id int, name, email string) bool {
	return m.updateUserCall.id == id && m.updateUserCall.name == name && m.updateUserCall.email == email
}

func (m *MockUserUsecase) WasDeleteUserCalledWith(id int) bool {
	return m.deleteUserCall == id
}

func TestNewServer(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	assert.NotNil(t, server)
	assert.Equal(t, mockUsecase, server.userUsecase)
}

func TestServer_writeJSONResponse(t *testing.T) {
	server := &Server{}
	w := httptest.NewRecorder()
	data := []byte(`{"test": "data"}`)

	server.writeJSONResponse(w, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, data, w.Body.Bytes())
}

func TestServer_writeErrorResponse(t *testing.T) {
	server := &Server{}
	w := httptest.NewRecorder()
	message := "test error"

	server.writeErrorResponse(w, http.StatusBadRequest, message)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response dto.MessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, message, response.Message)
}

func TestServer_extractAndValidateID(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		expectedID int
		expectErr  bool
	}{
		{"Valid ID", "id=123", 123, false},
		{"Invalid ID", "id=abc", 0, true},
		{"Missing ID", "", 0, true},
		{"Zero ID", "id=0", 0, true},
		{"Negative ID", "id=-1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			req := httptest.NewRequest("GET", "/test?"+tt.query, nil)

			id, err := server.extractAndValidateID(req)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestServer_decodeJSONBody(t *testing.T) {
	server := &Server{}

	t.Run("Valid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test","email":"test@example.com"}`))
		req.Header.Set("Content-Type", "application/json")

		var target dto.CreateUserRequest
		err := server.decodeJSONBody(req, &target)

		require.NoError(t, err)
		assert.Equal(t, "test", target.Name)
		assert.Equal(t, "test@example.com", target.Email)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{invalid json}`))
		req.Header.Set("Content-Type", "application/json")

		var target dto.CreateUserRequest
		err := server.decodeJSONBody(req, &target)

		assert.Error(t, err)
	})
}

func TestServer_CreateUser(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetCreateUserReturn([]byte(`{"id":1,"name":"John","email":"john@example.com"}`), nil)
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.CreateUser(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"John","email":"john@example.com"}`, w.Body.String())
	assert.True(t, mockUsecase.WasCreateUserCalledWith("John", "john@example.com"))
}

func TestServer_CreateUserWithInvalidJSON(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.CreateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_CreateUserWithUsecaseError(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetCreateUserReturn(nil, assert.AnError)
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
	req := httptest.NewRequest("POST", "/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.CreateUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestServer_GetUser(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetGetUserReturn([]byte(`{"id":1,"name":"John","email":"john@example.com"}`), nil)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()

	server.GetUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"John","email":"john@example.com"}`, w.Body.String())
	assert.True(t, mockUsecase.WasGetUserCalledWith(1))
}

func TestServer_GetUserWithMissingID(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()

	server.GetUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_GetUserWithInvalidID(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/user?id=invalid", nil)
	w := httptest.NewRecorder()

	server.GetUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_GetUserWithUsecaseError(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetGetUserReturn(nil, assert.AnError)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/user?id=1", nil)
	w := httptest.NewRecorder()

	server.GetUser(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestServer_GetAllUsers(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetGetAllUsersReturn([]byte(`{"users":[{"id":1,"name":"John","email":"john@example.com"}]}`), nil)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	server.GetAllUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"users":[{"id":1,"name":"John","email":"john@example.com"}]}`, w.Body.String())
	assert.True(t, mockUsecase.WasGetAllUsersCalled())
}

func TestServer_GetAllUsersWithUsecaseError(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetGetAllUsersReturn(nil, assert.AnError)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	server.GetAllUsers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestServer_UpdateUser(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetUpdateUserReturn([]byte(`{"id":1,"name":"Jane","email":"jane@example.com"}`), nil)
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{"name":"Jane","email":"jane@example.com"}`)
	req := httptest.NewRequest("PUT", "/user?id=1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.UpdateUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"id":1,"name":"Jane","email":"jane@example.com"}`, w.Body.String())
	assert.True(t, mockUsecase.WasUpdateUserCalledWith(1, "Jane", "jane@example.com"))
}

func TestServer_UpdateUserWithMissingID(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{"name":"Jane","email":"jane@example.com"}`)
	req := httptest.NewRequest("PUT", "/user", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.UpdateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_UpdateUserWithInvalidJSON(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest("PUT", "/user?id=1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.UpdateUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_UpdateUserWithUsecaseError(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetUpdateUserReturn(nil, assert.AnError)
	server := NewServer(mockUsecase)

	body := bytes.NewBufferString(`{"name":"Jane","email":"jane@example.com"}`)
	req := httptest.NewRequest("PUT", "/user?id=1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.UpdateUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestServer_DeleteUser(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetDeleteUserReturn([]byte(`{"message":"user deleted successfully"}`), nil)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("DELETE", "/user?id=1", nil)
	w := httptest.NewRecorder()

	server.DeleteUser(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, `{"message":"user deleted successfully"}`, w.Body.String())
	assert.True(t, mockUsecase.WasDeleteUserCalledWith(1))
}

func TestServer_DeleteUserWithMissingID(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("DELETE", "/user", nil)
	w := httptest.NewRecorder()

	server.DeleteUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestServer_DeleteUserWithUsecaseError(t *testing.T) {
	mockUsecase := &MockUserUsecase{}
	mockUsecase.SetDeleteUserReturn(nil, assert.AnError)
	server := NewServer(mockUsecase)

	req := httptest.NewRequest("DELETE", "/user?id=1", nil)
	w := httptest.NewRecorder()

	server.DeleteUser(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
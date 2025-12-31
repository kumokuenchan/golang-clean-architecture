package domain

import (
	"time"
)

// MockUserRepository implements UserRepository interface for testing
type MockUserRepository struct {
	users    []*User
	nextID   int
	err      error
	getByID  int
	getAll   bool
	create   *User
	update   *User
	delete   int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make([]*User, 0),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(user *User) error {
	m.create = user
	if m.err != nil {
		return m.err
	}
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users = append(m.users, user)
	m.nextID++
	return nil
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	if m.err != nil {
		m.getByID = id
		return nil, m.err
	}
	m.getByID = id
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) GetAll() ([]*User, error) {
	m.getAll = true
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *MockUserRepository) Update(user *User) error {
	m.update = user
	if m.err != nil {
		return m.err
	}
	for i, existingUser := range m.users {
		if existingUser.ID == user.ID {
			user.UpdatedAt = time.Now()
			m.users[i] = user
			return nil
		}
	}
	return ErrUserNotFound
}

func (m *MockUserRepository) Delete(id int) error {
	m.delete = id
	if m.err != nil {
		return m.err
	}
	for i, user := range m.users {
		if user.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}

// Helper methods for testing
func (m *MockUserRepository) SetError(err error) {
	m.err = err
}

func (m *MockUserRepository) WasCreateCalledWith(user *User) bool {
	if m.create == nil || user == nil {
		return m.create == user
	}
	return m.create.Name == user.Name && m.create.Email == user.Email
}

func (m *MockUserRepository) WasGetByIDCalledWith(id int) bool {
	return m.getByID == id
}

func (m *MockUserRepository) WasGetAllCalled() bool {
	return m.getAll
}

func (m *MockUserRepository) WasUpdateCalledWith(user *User) bool {
	if m.update == nil || user == nil {
		return m.update == user
	}
	return m.update.ID == user.ID && m.update.Name == user.Name && m.update.Email == user.Email
}

func (m *MockUserRepository) WasDeleteCalledWith(id int) bool {
	return m.delete == id
}

func (m *MockUserRepository) GetUsers() []*User {
	return m.users
}

func (m *MockUserRepository) AddUser(user *User) {
	user.ID = m.nextID
	m.nextID++
	m.users = append(m.users, user)
}
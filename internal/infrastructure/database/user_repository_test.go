package database

import (
	"database/sql"
	"testing"
	"time"

	"clean-arch-sample/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMySQLUserRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestMySQLUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	user := &domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock the INSERT query
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	require.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_CreateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	user := &domain.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Mock database error
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	err = repo.Create(user)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	expectedUser := &domain.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.CreatedAt, expectedUser.UpdatedAt)
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock no rows found
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetByIDError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock database error
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WithArgs(1).
		WillReturnError(assert.AnError)

	user, err := repo.GetByID(1)
	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	expectedUsers := []*domain.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "User 2",
			Email:     "user2@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock the SELECT query
	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	}
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WillReturnRows(rows)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].ID, users[0].ID)
	assert.Equal(t, expectedUsers[1].ID, users[1].ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetAllError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock database error
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WillReturnError(assert.AnError)

	users, err := repo.GetAll()
	assert.Error(t, err)
	assert.Nil(t, users)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_GetAllEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock empty result
	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WillReturnRows(rows)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, users, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	user := &domain.User{
		ID:    1,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	// Mock the UPDATE query
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(user)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_UpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)
	user := &domain.User{
		ID:    1,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	// Mock database error
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnError(assert.AnError)

	err = repo.Update(user)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock the DELETE query
	mock.ExpectExec("DELETE FROM users").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_DeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Mock database error
	mock.ExpectExec("DELETE FROM users").
		WithArgs(1).
		WillReturnError(assert.AnError)

	err = repo.Delete(1)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_Integration(t *testing.T) {
	// This test demonstrates how the repository methods work together
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewMySQLUserRepository(db)

	// Test Create
	user := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	require.NoError(t, err)
	assert.Equal(t, 1, user.ID)

	// Test GetByID
	rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WithArgs(user.ID).
		WillReturnRows(rows)

	retrievedUser, err := repo.GetByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, retrievedUser.ID)
	assert.Equal(t, user.Name, retrievedUser.Name)

	// Test Update
	user.Name = "Updated User"
	mock.ExpectExec("UPDATE users").
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(user)
	require.NoError(t, err)

	// Test GetAll
	allRows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
		AddRow(user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	mock.ExpectQuery("SELECT id, name, email, created_at, updated_at FROM users").
		WillReturnRows(allRows)

	allUsers, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, allUsers, 1)

	// Test Delete
	mock.ExpectExec("DELETE FROM users").
		WithArgs(user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(user.ID)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
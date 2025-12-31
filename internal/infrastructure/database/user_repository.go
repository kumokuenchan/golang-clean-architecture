package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kuen/clean-arch-sample/internal/domain"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *MySQLUserRepository) GetByID(id int) (*domain.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *MySQLUserRepository) GetAll() ([]*domain.User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *MySQLUserRepository) Update(user *domain.User) error {
	query := "UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?"
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Name, user.Email, user.UpdatedAt, user.ID)
	return err
}

func (r *MySQLUserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
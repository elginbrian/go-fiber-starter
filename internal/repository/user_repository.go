package repository

import (
	"database/sql"
	"errors"
	"fiber-starter/internal/domain"
	"time"
)

type UserRepository interface {
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id int) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	UpdateUser(id int, user domain.User) (domain.User, error)
	DeleteUser(id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}

func (r *userRepository) GetUserByID(id int) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user domain.User) (domain.User, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())", 
		user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	return user, nil
}

func (r *userRepository) UpdateUser(id int, user domain.User) (domain.User, error) {
	_, err := r.db.Exec("UPDATE users SET name = ?, email = ?, password_hash = ?, updated_at = NOW() WHERE id = ?", 
		user.Name, user.Email, user.PasswordHash, id)
	if err != nil {
		return user, err
	}

	user.ID = id
	user.UpdatedAt = time.Now()
	return user, nil
}

func (r *userRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

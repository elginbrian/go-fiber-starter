package repository

import (
	"database/sql"
	"fiber-starter/internal/domain"
)

type AuthRepository interface {
    GetUserByEmail(email string) (*domain.User, error)
}

type authRepository struct {
    db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
    return &authRepository{db: db}
}

func (r *authRepository) GetUserByEmail(email string) (*domain.User, error) {
    var user domain.User
    query := "SELECT id, name, email FROM users WHERE email = ?"
    row := r.db.QueryRow(query, email)

    if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil 
        }
        return nil, err
    }
    return &user, nil
}

package repository

import (
	"context"
	"errors"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) contract.IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, email, bio, image_url, created_at, updated_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Bio, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(ctx, "SELECT id, name, email, bio, image_url, created_at, updated_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Bio, &user.ImageURL, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, fmt.Errorf("error fetching user: %w", err)
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	if user.ImageURL == "" {
		user.ImageURL = "https://static.vecteezy.com/system/resources/previews/009/292/244/non_2x/default-avatar-icon-of-social-media-user-vector.jpg"
	}
	if user.Bio == "" {
		user.Bio = "Hi there!"
	}

	commandTag, err := r.db.Exec(ctx,
		"INSERT INTO users (name, email, password_hash, image_url, bio, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())",
		user.Name, user.Email, user.PasswordHash, user.ImageURL, user.Bio)
	if err != nil {
		return user, fmt.Errorf("error creating user: %w", err)
	}

	if commandTag.RowsAffected() > 0 {
		err := r.db.QueryRow(ctx, 
			"SELECT id, name, email, image_url, bio, created_at, updated_at FROM users WHERE email = $1", 
			user.Email).Scan(&user.ID, &user.Name, &user.Email, &user.ImageURL, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return user, fmt.Errorf("error fetching user after creation: %w", err)
		}
		return user, nil
	}

	return user, errors.New("failed to create user")
}


func (r *userRepository) UpdateUser(ctx context.Context, id string, user entity.User) (entity.User, error) { 
	commandTag, err := r.db.Exec(ctx, 
		`UPDATE users 
		 SET name = $1, email = $2, image_url = $3, bio = $4, updated_at = NOW() 
		 WHERE id = $5`, 
		user.Name, user.Email, user.ImageURL, user.Bio, id)
		
	if err != nil {
		return user, fmt.Errorf("error updating user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return user, errors.New("user not found")
	}

	user.ID = id
	user.UpdatedAt = time.Now()
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	commandTag, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) SearchUsers(ctx context.Context, query string) ([]entity.User, error) {
	rows, err := r.db.Query(ctx, 
		"SELECT id, name, email, created_at, updated_at FROM users WHERE LOWER(name) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1)", 
		"%"+query+"%")	
    if err != nil {
        return nil, fmt.Errorf("error searching users: %w", err)
    }
    defer rows.Close()

    var users []entity.User
    for rows.Next() {
        var user entity.User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
            return nil, fmt.Errorf("error scanning user row: %w", err)
        }
        users = append(users, user)
    }

    if len(users) == 0 {
        return nil, errors.New("no users found")
    }

    return users, nil
}
package repository

import (
	"context"
	"database/sql"
	"fiber-starter/internal/domain"
)

type PostRepository interface {
	FetchAllPosts(ctx context.Context) ([]domain.Post, error)
	FetchPostByID(ctx context.Context, postID int) (*domain.Post, error)
	CreatePost(ctx context.Context, post domain.Post) (*domain.Post, error)
	UpdatePost(ctx context.Context, postID int, post domain.Post) (*domain.Post, error)
	DeletePost(ctx context.Context, postID int) error
}

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FetchAllPosts(ctx context.Context) ([]domain.Post, error) {
	query := "SELECT id, user_id, caption, image_url, created_at, updated_at FROM posts"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Caption, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepository) FetchPostByID(ctx context.Context, postID int) (*domain.Post, error) {
	query := "SELECT id, user_id, caption, image_url, created_at, updated_at FROM posts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, postID)

	var post domain.Post
	if err := row.Scan(&post.ID, &post.UserID, &post.Caption, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) CreatePost(ctx context.Context, post domain.Post) (*domain.Post, error) {
	query := "INSERT INTO posts (user_id, caption, image_url) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	err := r.db.QueryRowContext(ctx, query, post.UserID, post.Caption, post.ImageURL).Scan(
		&post.ID, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, postID int, post domain.Post) (*domain.Post, error) {
	query := "UPDATE posts SET caption = $1, image_url = $2, updated_at = NOW() WHERE id = $3 RETURNING id, user_id, caption, image_url, created_at, updated_at"
	err := r.db.QueryRowContext(ctx, query, post.Caption, post.ImageURL, postID).Scan(
		&post.ID, &post.UserID, &post.Caption, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) DeletePost(ctx context.Context, postID int) error {
	query := "DELETE FROM posts WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}
	return nil
}
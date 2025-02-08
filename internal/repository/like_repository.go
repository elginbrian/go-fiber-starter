package repository

import (
	"context"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type likeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) contract.ILikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) GetLikesByPostID(ctx context.Context, postID string) ([]entity.Like, error) {
	query := "SELECT id, user_id, post_id, created_at FROM likes WHERE post_id = $1"
	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("error fetching likes for post %s: %w", postID, err)
	}
	defer rows.Close()

	var likes []entity.Like
	for rows.Next() {
		var like entity.Like
		if err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning like row: %w", err)
		}
		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return likes, nil
}

func (r *likeRepository) GetLikesByUserID(ctx context.Context, userID string) ([]entity.Like, error) {
	query := "SELECT id, user_id, post_id, created_at FROM likes WHERE user_id = $1"
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching likes for user %s: %w", userID, err)
	}
	defer rows.Close()

	var likes []entity.Like
	for rows.Next() {
		var like entity.Like
		if err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning like row: %w", err)
		}
		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return likes, nil
}

func (r *likeRepository) AddLike(ctx context.Context, like entity.Like) (*entity.Like, error) {
	query := "INSERT INTO likes (user_id, post_id) VALUES ($1, $2) RETURNING id, created_at"
	err := r.db.QueryRow(ctx, query, like.UserID, like.PostID).Scan(&like.ID, &like.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error adding like: %w", err)
	}
	like.CreatedAt = time.Now()
	return &like, nil
}

func (r *likeRepository) RemoveLike(ctx context.Context, userID, postID string) error {
	query := "DELETE FROM likes WHERE user_id = $1 AND post_id = $2"
	result, err := r.db.Exec(ctx, query, userID, postID)
	if err != nil {
		return fmt.Errorf("error removing like: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("like not found for user %s and post %s", userID, postID)
	}
	return nil
}

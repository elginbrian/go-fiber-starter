package repository

import (
	"context"
	contract "fiber-starter/domain/contract"
	entity "fiber-starter/domain/entity"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type commentRepository struct {
    db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) contract.ICommentRepository {
    return &commentRepository{db: db}
}

func (r *commentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]entity.Comment, error) {
    query := "SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE post_id = $1 ORDER BY created_at DESC"
    rows, err := r.db.Query(ctx, query, postID)
    if err != nil {
        return nil, fmt.Errorf("error fetching comments: %w", err)
    }
    defer rows.Close()

    var comments []entity.Comment
    for rows.Next() {
        var comment entity.Comment
        if err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
            return nil, fmt.Errorf("error scanning comment row: %w", err)
        }
        comments = append(comments, comment)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over rows: %w", err)
    }
    return comments, nil
}

func (r *commentRepository) GetCommentByID(ctx context.Context, commentID string) (*entity.Comment, error) {
    query := "SELECT id, user_id, post_id, content, created_at, updated_at FROM comments WHERE id = $1"
    var comment entity.Comment
    err := r.db.QueryRow(ctx, query, commentID).Scan(
        &comment.ID, &comment.UserID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
    )
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("comment not found")
        }
        return nil, fmt.Errorf("error retrieving comment: %w", err)
    }
    return &comment, nil
}

func (r *commentRepository) CreateComment(ctx context.Context, comment entity.Comment) (*entity.Comment, error) {
    query := "INSERT INTO comments (user_id, post_id, content) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
    err := r.db.QueryRow(ctx, query, comment.UserID, comment.PostID, comment.Content).Scan(
        &comment.ID, &comment.CreatedAt, &comment.UpdatedAt,
    )
    if err != nil {
        return nil, fmt.Errorf("error creating comment: %w", err)
    }
    return &comment, nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, commentID string) error {
    query := "DELETE FROM comments WHERE id = $1"
    result, err := r.db.Exec(ctx, query, commentID)
    if err != nil {
        return fmt.Errorf("error deleting comment: %w", err)
    }

    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return pgx.ErrNoRows
    }
    return nil
}
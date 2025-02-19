package repository

import (
	"context"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
)

const (
	getFeedQuery = `WITH friends AS (
	SELECT uf1.friend_id AS id
	FROM users_friends uf1
	INNER JOIN users_friends uf2 
		ON uf2.id = uf1.friend_id 
		AND uf2.friend_id = uf1.id
	WHERE uf1.id = $1
)
SELECT p.* 
FROM friends AS f
INNER JOIN posts AS p 
	ON p.user_id = f.id
ORDER BY p.id DESC
LIMIT $2`
	listQuery = `"SELECT 
    id, user_id, message, created_at, modified_at 
FROM posts
WHERE id < $3
    AND user_id = $1
ORDER BY id DESC 
LIMIT $2"`
)

type postRepository struct {
	ctx       context.Context
	dbCluster *postgresqldb.DBCluster
}

func NewPostRepository(ctx context.Context, dbcluster *postgresqldb.DBCluster) domain.PostRepository {
	return &postRepository{ctx: ctx, dbCluster: dbcluster}
}

func (r *postRepository) List(userId domain.UserKey, limit int, lastPostId domain.PostKey) ([]*domain.Post, error) {
	var posts []*domain.Post

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("postRepository.List: r.dbCluster.GetDB returned error %w", err)
	}
	rows, err := db.Query(r.ctx, listQuery, userId, limit, lastPostId)
	if err != nil {
		return nil, fmt.Errorf("postRepository.List: r.db.QueryRow returned error %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Text, &post.CreatedAt, &post.ModifiedAt); err != nil {
			return nil, fmt.Errorf("postRepository.List: rows.Scan returned error: %w", err)
		}
		posts = append(posts, &post)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postRepository.List: rows iteration error: %w", err)
	}

	return posts, nil
}

// Добавление нового поста
func (r *postRepository) Create(userId domain.UserKey, message domain.PostText) (*domain.Post, error) {

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return nil, fmt.Errorf("postRepository.CreatePost: r.dbCluster.GetDB returned error %w", err)
	}
	var post domain.Post

	err = db.QueryRow(r.ctx, "INSERT INTO posts (user_id, message) VALUES ($1, $2) RETURNING id, user_id, message, created_at, modified_at",
		userId, message).Scan(&post.Id, &post.UserId, &post.Text, &post.CreatedAt, &post.ModifiedAt)
	if err != nil {
		return nil, fmt.Errorf("postRepository.CreatePost: r.db.QueryRow returned error %w", err)
	}
	return &post, nil
}

// Получение поста по ID
func (r *postRepository) Get(id domain.PostKey) (*domain.Post, error) {
	var post domain.Post

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("postRepository.Get: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT id, user_id, message, created_at, modified_at FROM posts WHERE id = $1", id).Scan(
		&post.Id, &post.UserId, &post.Text, &post.CreatedAt, &post.ModifiedAt)
	if err != nil {
		return nil, fmt.Errorf("postRepository.Get: r.db.QueryRow returned error %w", err)
	}
	return &post, nil
}

// Обновление сообщения поста по ID
func (r *postRepository) UpdatePost(postId domain.PostKey, newMessage domain.PostText) (*domain.Post, error) {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return nil, fmt.Errorf("postRepository.UpdatePost: r.dbCluster.GetDB returned error %w", err)
	}
	var post domain.Post
	err = db.QueryRow(r.ctx,
		"UPDATE posts SET message = $1, modified_at = NOW() WHERE id = $2 RETURNING id, user_id, message, created_at, modified_at", newMessage, postId).Scan(
		&post.Id, &post.UserId, &post.Text, &post.CreatedAt, &post.ModifiedAt)

	if err != nil {
		return nil, fmt.Errorf("postRepository.UpdatePost: db.Exec returned error %w", err)
	}
	return &post, nil
}

// Удаление поста по ID
func (r *postRepository) Delete(postId domain.PostKey) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("postRepository.Delete: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Exec(r.ctx, "DELETE FROM posts WHERE id = $1", postId)
	if err != nil {
		return fmt.Errorf("postRepository.Delete: db.Exec returned error %w", err)
	}

	rows := q.RowsAffected()

	if rows != 1 {
		return domain.ErrObjectNotFound
	}
	return nil
}

func (r *postRepository) GetPostOwner(postId domain.PostKey) (domain.UserKey, error) {
	var post domain.UserKey

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return "", fmt.Errorf("postRepository.GetPostOwner: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT user_id FROM posts WHERE id = $1", postId).Scan(&post)
	if err != nil {
		return "", fmt.Errorf("postRepository.GetPostOwner: r.db.QueryRow returned error %w", err)
	}
	return post, nil
}

func (r *postRepository) GetFeed(userId domain.UserKey, limit int) ([]*domain.Post, error) {
	var posts []*domain.Post

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("postRepository.GetFeed: r.dbCluster.GetDB returned error %w", err)
	}
	rows, err := db.Query(r.ctx, getFeedQuery, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("postRepository.GetFeed: r.db.QueryRow returned error %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.Id, &post.UserId, &post.Text, &post.CreatedAt, &post.ModifiedAt); err != nil {
			return nil, fmt.Errorf("postRepository.GetFeed: rows.Scan returned error: %w", err)
		}
		posts = append(posts, &post)
	}

	// Проверка на наличие ошибок после завершения итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postRepository.GetFeed: rows iteration error: %w", err)
	}

	return posts, nil
}

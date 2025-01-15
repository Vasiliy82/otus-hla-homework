package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/infrastructure/postgresqldb"
	"github.com/lib/pq"
)

const (
	getFriendsIds_query = `WITH friends AS (
    SELECT uf1.friend_id AS id
    FROM users_friends uf1
    INNER JOIN users_friends uf2 
        ON uf2.id = uf1.friend_id 
        AND uf2.friend_id = uf1.id
    WHERE uf1.id = $1
)
SELECT friends.id FROM friends`

	setLastActivity_query = `INSERT INTO users_last_activity (id, last_activity) 
VALUES($1, NOW())
ON CONFLICT (id) DO UPDATE SET last_activity = EXCLUDED.last_activity;`

	getUsersActiveSince_query = "SELECT id FROM users_last_activity WHERE last_activity >= NOW() - INTERVAL '1 second' * $1;"
)

type userRepository struct {
	ctx       context.Context
	dbCluster *postgresqldb.DBCluster
}

func NewUserRepository(ctx context.Context, dbcluster *postgresqldb.DBCluster) domain.UserRepository {
	return &userRepository{ctx: ctx, dbCluster: dbcluster}
}

func (r *userRepository) RegisterUser(user *domain.User) (domain.UserKey, error) {
	var userId domain.UserKey

	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return "", fmt.Errorf("userRepository.RegisterUser: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "INSERT INTO users (first_name, last_name, birthdate, biography, city, username, password_hash) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.FirstName, user.LastName, user.Birthdate, user.Biography, user.City, user.Username, user.PasswordHash).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("userRepository.RegisterUser: r.db.QueryRow returned error %w", err)
	}
	// Счетчик для ДЗ № 3 (необходимо считать количество успешно записанных данных)
	incRecordsCount()

	return userId, nil
}

func (r *userRepository) GetByID(id domain.UserKey) (*domain.User, error) {
	var user domain.User

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByUsername: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow(r.ctx, "SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByUsername: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}

func (r *userRepository) Search(firstName, lastName string) ([]*domain.User, error) {
	var users []*domain.User

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Search: r.dbCluster.GetDB returned error %w", err)
	}

	ptnFirstName := fmt.Sprintf("%s%%", firstName)
	ptnLastName := fmt.Sprintf("%s%%", lastName)

	q, err := db.Query(r.ctx, "SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id", ptnFirstName, ptnLastName)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Search: r.db.Query returned error %w", err)
	}
	defer q.Close()

	for q.Next() {
		user := domain.User{}
		err := q.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Birthdate,
			&user.Biography, &user.City, &user.Username, &user.PasswordHash,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("userRepository.Search: q.Scan returned error %w", err)
		}
		users = append(users, &user)
	}

	// проверял максимальное количество потоков, создаваемое JMeter
	// time.Sleep(1000 * time.Millisecond)
	//

	return users, nil
}

func (r *userRepository) AddFriend(my_id, friend_id domain.UserKey) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("userRepository.AddFriend: r.dbCluster.GetDB returned error %w", err)
	}

	_, err = db.Exec(r.ctx, "INSERT INTO users_friends (id, friend_id) SELECT $1 AS id, $2 AS friend_id", my_id, friend_id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // duplicate key value violates unique constraint
				return domain.ErrObjectAlreadyExists
			}
			if pqErr.Code == "23503" { // insert or update violates foreign key constraint
				return domain.ErrObjectNotFound
			}
		}
		// прочие ошибки
		return fmt.Errorf("userRepository.AddFriend: r.db.Exec returned error %w", err)
	}
	return nil
}

func (r *userRepository) RemoveFriend(my_id, friend_id domain.UserKey) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("userRepository.RemoveFriend: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Exec(r.ctx, "DELETE FROM users_friends WHERE id = $1 AND friend_id = $2", my_id, friend_id)
	if err != nil {
		return fmt.Errorf("userRepository.RemoveFriend: r.db.Exec returned error %w", err)
	}
	rows := q.RowsAffected()
	if rows != 1 {
		return domain.ErrObjectNotFound
	}
	return nil
}

func (r *userRepository) GetFriendsIds(id domain.UserKey) ([]domain.UserKey, error) {
	var result []domain.UserKey

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetFriendsIds: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Query(r.ctx, getFriendsIds_query, id)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetFriendsIds: r.db.Query returned error %w", err)
	}
	defer q.Close()

	for q.Next() {
		var friendId domain.UserKey
		err := q.Scan(&friendId)
		if err != nil {
			return nil, fmt.Errorf("userRepository.GetFriendsIds: q.Scan returned error %w", err)
		}
		result = append(result, friendId)
	}
	return result, nil
}

func (r *userRepository) SetLastActivity(id domain.UserKey) error {
	db, err := r.dbCluster.GetDBPool(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("userRepository.SetLastActivity: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Exec(r.ctx, setLastActivity_query, id)
	if err != nil {
		return fmt.Errorf("userRepository.SetLastActivity: r.db.Exec returned error %w", err)
	}
	rows := q.RowsAffected()
	if rows != 1 {
		return domain.ErrObjectNotFound
	}
	return nil
}

func (r *userRepository) GetUsersActiveSince(period time.Duration) ([]domain.UserKey, error) {
	var result []domain.UserKey

	db, err := r.dbCluster.GetDBPool(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetUsersActiveSince: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Query(r.ctx, getUsersActiveSince_query, period.Seconds())
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetUsersActiveSince: r.db.Query returned error %w", err)
	}
	defer q.Close()

	for q.Next() {
		var friendId domain.UserKey
		err := q.Scan(&friendId)
		if err != nil {
			return nil, fmt.Errorf("userRepository.GetUsersActiveSince: q.Scan returned error %w", err)
		}
		result = append(result, friendId)
	}
	return result, nil
}

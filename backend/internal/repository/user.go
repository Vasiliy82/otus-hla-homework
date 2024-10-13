package repository

import (
	"errors"
	"fmt"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/infrastructure/postgresqldb"
	"github.com/lib/pq"
)

type userRepository struct {
	dbCluster *postgresqldb.DBCluster
}

func NewUserRepository(dbcluster *postgresqldb.DBCluster) domain.UserRepository {
	return &userRepository{dbCluster: dbcluster}
}

func (r *userRepository) RegisterUser(user *domain.User) (string, error) {
	var userId string

	db, err := r.dbCluster.GetDB(postgresqldb.ReadWrite)
	if err != nil {
		return "", fmt.Errorf("userRepository.RegisterUser: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow("INSERT INTO users (first_name, last_name, birthdate, biography, city, username, password_hash) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.FirstName, user.LastName, user.Birthdate, user.Biography, user.City, user.Username, user.PasswordHash).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("userRepository.RegisterUser: r.db.QueryRow returned error %w", err)
	}
	return userId, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User

	db, err := r.dbCluster.GetDB(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow("SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByID: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	db, err := r.dbCluster.GetDB(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByUsername: r.dbCluster.GetDB returned error %w", err)
	}

	err = db.QueryRow("SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Birthdate, &user.Biography, &user.City, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("userRepository.GetByUsername: r.db.QueryRow returned error %w", err)
	}
	return &user, nil
}

func (r *userRepository) Search(firstName, lastName string) ([]*domain.User, error) {
	var users []*domain.User

	db, err := r.dbCluster.GetDB(postgresqldb.Read)
	if err != nil {
		return nil, fmt.Errorf("userRepository.Search: r.dbCluster.GetDB returned error %w", err)
	}

	ptnFirstName := fmt.Sprintf("%s%%", firstName)
	ptnLastName := fmt.Sprintf("%s%%", lastName)

	q, err := db.Query("SELECT id, first_name, last_name, birthdate, biography, city, username, password_hash, created_at, updated_at FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id", ptnFirstName, ptnLastName)
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

func (r *userRepository) AddFriend(my_id, friend_id string) error {
	db, err := r.dbCluster.GetDB(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("userRepository.AddFriend: r.dbCluster.GetDB returned error %w", err)
	}

	_, err = db.Exec("INSERT INTO users_friends (id, friend_id) SELECT $1 AS id, $2 AS friend_id", my_id, friend_id)
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

func (r *userRepository) RemoveFriend(my_id, friend_id string) error {
	db, err := r.dbCluster.GetDB(postgresqldb.ReadWrite)
	if err != nil {
		return fmt.Errorf("userRepository.RemoveFriend: r.dbCluster.GetDB returned error %w", err)
	}

	q, err := db.Exec("DELETE FROM users_friends WHERE id = $1 AND friend_id = $2", my_id, friend_id)
	if err != nil {
		return fmt.Errorf("userRepository.RemoveFriend: r.db.Exec returned error %w", err)
	}
	rows, err := q.RowsAffected()
	if err != nil {
		return fmt.Errorf("userRepository.RemoveFriend: q.RowsAffected returned error %w", err)
	}
	if rows != 1 {
		return domain.ErrObjectNotFound
	}
	return nil
}

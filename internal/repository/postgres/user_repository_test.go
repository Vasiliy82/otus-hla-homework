package postgres_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/repository/postgres"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_RegisterUser_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	userRepo := postgres.NewUserRepository(db)

	// Создаем тестового пользователя
	testUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.SecondName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

	// Вызываем метод репозитория
	userID, err := userRepo.RegisterUser(testUser)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, "123", userID)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_RegisterUser_DuplicateUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	testUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем ошибку дублирования
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.SecondName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnError(&pq.Error{Code: "23505"}) // Код ошибки уникального ограничения

	userID, err := userRepo.RegisterUser(testUser)

	// Проверяем, что вернулась ошибка конфликта и ID не был сгенерирован
	assert.ErrorIs(t, err, domain.ErrConflict)
	assert.Equal(t, "", userID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	testUser := domain.User{
		ID:         "123",
		FirstName:  "John",
		SecondName: "Doe",
		Username:   "johndoe",
	}

	// Эмулируем успешный результат SELECT
	mock.ExpectQuery("SELECT id, first_name, second_name").
		WithArgs("123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "second_name", "birthdate", "biography", "city", "username", "created_at"}).
			AddRow(testUser.ID, testUser.FirstName, testUser.SecondName, time.Now(), "", "", testUser.Username, time.Now()))

	user, err := userRepo.GetUserByID("123")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	// Эмулируем ошибку, что пользователь не найден
	mock.ExpectQuery("SELECT id, first_name, second_name").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetUserByID("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, domain.User{}, user)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_CheckUserPasswordHash_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	// Эмулируем успешный результат SELECT для проверки пароля
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE username = $1 AND password_hash = $2")).
		WithArgs("johndoe", "hashedpassword").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

	userID, err := userRepo.CheckUserPasswordHash("johndoe", "hashedpassword")

	// Проверяем, что ошибок нет и пользователь найден
	assert.NoError(t, err)
	assert.Equal(t, "123", userID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_CheckUserPasswordHash_WrongPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)

	// Эмулируем ошибку, что пароль не совпал
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM users WHERE username = $1 AND password_hash = $2")).
		WithArgs("johndoe", "wrongpassword").
		WillReturnError(sql.ErrNoRows)

	userID, err := userRepo.CheckUserPasswordHash("johndoe", "wrongpassword")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

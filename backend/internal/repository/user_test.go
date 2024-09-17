package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vasiliy82/otus-hla-homework/domain"
	repository "github.com/Vasiliy82/otus-hla-homework/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_RegisterUser_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	userRepo := repository.NewUserRepository(db)

	// Создаем тестового пользователя
	testUser := domain.User{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.LastName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

	// Вызываем метод репозитория
	userID, err := userRepo.RegisterUser(&testUser)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, "123", userID)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_RegisterUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	testUser := domain.User{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем ошибку дублирования
	mock.ExpectQuery("^INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.LastName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnError(errors.New("db error")) // Код ошибки уникального ограничения

	userID, err := userRepo.RegisterUser(&testUser)

	// Проверяем, что вернулась ошибка конфликта и ID не был сгенерирован
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	testUser := domain.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe@gmail.com",
	}

	// Эмулируем успешный результат SELECT
	mock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "second_name", "birthdate", "biography", "city", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.FirstName, testUser.LastName, time.Now(), "", "", testUser.Username, testUser.PasswordHash, time.Now(), time.Now()))

	user, err := userRepo.GetByID("123")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@gmail.com", user.Username)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	// Эмулируем ошибку, что пользователь не найден
	mock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetByID("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, &domain.User{}, user)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	testUser := domain.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe@gmail.com",
	}

	// Эмулируем успешный результат SELECT
	mock.ExpectQuery("^SELECT").
		WithArgs("johndoe@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "second_name", "birthdate", "biography", "city", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.FirstName, testUser.LastName, time.Now(), "", "", testUser.Username, testUser.PasswordHash, time.Now(), time.Now()))

	user, err := userRepo.GetByUsername("johndoe@gmail.com")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@gmail.com", user.Username)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	// Эмулируем ошибку, что пользователь не найден
	mock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetByID("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, &domain.User{}, user)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

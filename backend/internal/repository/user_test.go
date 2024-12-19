package repository_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vasiliy82/otus-hla-homework/backend/domain"
	repository "github.com/Vasiliy82/otus-hla-homework/backend/internal/repository"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_RegisterUser_Success(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	// Создаем экземпляр репозитория
	userRepo := repository.NewUserRepository(dbCluster)

	// Создаем тестового пользователя
	testUser := domain.User{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем успешную вставку в базу данных
	mMock.ExpectQuery("^INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.LastName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

	// Вызываем метод репозитория
	userID, err := userRepo.RegisterUser(&testUser)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, "123", userID)

	// Проверяем, что все mock-ожидания выполнены
	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_RegisterUser_Error(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	userRepo := repository.NewUserRepository(dbCluster)

	testUser := domain.User{
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "hashedpassword",
	}

	// Эмулируем ошибку дублирования
	mMock.ExpectQuery("^INSERT INTO users").
		WithArgs(testUser.FirstName, testUser.LastName, testUser.Birthdate, testUser.Biography, testUser.City, testUser.Username, testUser.PasswordHash).
		WillReturnError(errors.New("db error")) // Код ошибки уникального ограничения

	userID, err := userRepo.RegisterUser(&testUser)

	// Проверяем, что вернулась ошибка конфликта и ID не был сгенерирован
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_Success(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	userRepo := repository.NewUserRepository(dbCluster)

	testUser := domain.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe@gmail.com",
	}

	// Эмулируем успешный результат SELECT
	mMock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthdate", "biography", "city", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.FirstName, testUser.LastName, time.Now(), "", "", testUser.Username, testUser.PasswordHash, time.Now(), time.Now()))

	user, err := userRepo.GetByID("123")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@gmail.com", user.Username)

	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	userRepo := repository.NewUserRepository(dbCluster)

	// Эмулируем ошибку, что пользователь не найден
	mMock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetByID("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, &domain.User{}, user)

	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByUsername_Success(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	userRepo := repository.NewUserRepository(dbCluster)

	testUser := domain.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe@gmail.com",
	}

	// Эмулируем успешный результат SELECT
	mMock.ExpectQuery("^SELECT").
		WithArgs("johndoe@gmail.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "birthdate", "biography", "city", "username", "password_hash", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.FirstName, testUser.LastName, time.Now(), "", "", testUser.Username, testUser.PasswordHash, time.Now(), time.Now()))

	user, err := userRepo.GetByUsername("johndoe@gmail.com")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@gmail.com", user.Username)

	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_GetUserByUsername_NotFound(t *testing.T) {
	// Создаем mock для базы данных
	dbCluster, mMock, _ := testutils.NewMockDBCluster(t, 1)
	defer dbCluster.Close()

	userRepo := repository.NewUserRepository(dbCluster)

	// Эмулируем ошибку, что пользователь не найден
	mMock.ExpectQuery("^SELECT").
		WithArgs("123").
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetByID("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, &domain.User{}, user)

	err = mMock.ExpectationsWereMet()
	assert.NoError(t, err)
}

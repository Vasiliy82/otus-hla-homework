package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	repository "github.com/Vasiliy82/otus-hla-homework/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestBlacklistRepository_AddToBlacklist_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	serial := "12345"

	// Эмулируем успешную вставку в базу данных
	mock.ExpectExec("^INSERT INTO blacklisted").
		WithArgs(serial, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Эмулируем успешную вставку

	// Вызываем метод репозитория
	err = blRepo.AddToBlacklist(serial, time.Now().Add(24*time.Hour))

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_AddToBlacklist_Failure(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	serial := "12345"

	// Эмулируем успешную вставку в базу данных
	mock.ExpectExec("^INSERT INTO blacklisted").
		WithArgs(serial, sqlmock.AnyArg()).WillReturnError(errors.New("DB Error"))

	// Вызываем метод репозитория
	err = blRepo.AddToBlacklist(serial, time.Now().Add(24*time.Hour))

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.Error(t, err)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_IsBlacklisted_Success_False(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	serial := "12345"

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^SELECT 1 FROM blacklisted").
		WithArgs(serial).
		WillReturnRows(sqlmock.NewRows([]string{""}))

	// Вызываем метод репозитория
	result, err := blRepo.IsBlacklisted(serial)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, false, result)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_IsBlacklisted_Success_True(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	serial := "12345"

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^SELECT 1 FROM blacklisted").
		WithArgs(serial).
		WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(1))

	// Вызываем метод репозитория
	result, err := blRepo.IsBlacklisted(serial)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_IsBlacklisted_Fail(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	serial := "12345"

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^SELECT 1 FROM blacklisted").
		WithArgs(serial).
		WillReturnError(errors.New("DB Error"))

	// Вызываем метод репозитория
	result, err := blRepo.IsBlacklisted(serial)

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.Error(t, err)
	assert.Equal(t, false, result)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_NewSerial_Success(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^SELECT nextval").
		WillReturnRows(mock.NewRows([]string{""}).AddRow(400))

	// Вызываем метод репозитория
	result, err := blRepo.NewSerial()

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.NoError(t, err)
	assert.Equal(t, "400", result)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBlacklistRepository_NewSerial_Fail(t *testing.T) {
	// Создаем mock для базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создаем экземпляр репозитория
	blRepo := repository.NewBlacklistRepository(db)

	// Эмулируем успешную вставку в базу данных
	mock.ExpectQuery("^SELECT nextval").
		WillReturnError(errors.New("DB Error"))

	// Вызываем метод репозитория
	result, err := blRepo.NewSerial()

	// Проверяем, что ошибок нет и ID пользователя вернулся
	assert.Error(t, err)
	assert.Equal(t, "", result)

	// Проверяем, что все mock-ожидания выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockUserRepository) RegisterUser(user domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserRepository) CheckUserPasswordHash(username string, passwordHash string) (string, error) {
	args := m.Called(username, passwordHash)
	return args.String(0), args.Error(1)
}

func (m *MockSessionRepository) CreateSession(userID, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)

}

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)
	testUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe",
		PasswordHash: "hashedpassword",
	}

	// Мокаем успешную регистрацию
	mockRepo.On("RegisterUser", testUser).Return("123", nil)

	userID, err := userService.RegisterUser(testUser)

	// Проверяем, что ошибок нет и ID пользователя возвращен
	assert.NoError(t, err)
	assert.Equal(t, "123", userID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_DuplicateUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)
	testUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe",
		PasswordHash: "hashedpassword",
	}

	// Мокаем ошибку дублирования логина
	mockRepo.On("RegisterUser", testUser).Return("", domain.ErrConflict)

	userID, err := userService.RegisterUser(testUser)

	// Проверяем, что вернулась ошибка конфликта и ID не сгенерирован
	assert.ErrorIs(t, err, domain.ErrConflict)
	assert.Equal(t, "", userID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)
	testUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe",
		PasswordHash: "hashedpassword",
	}

	// Мокаем ошибку базы данных
	mockRepo.On("RegisterUser", testUser).Return("", errors.New("db error"))

	userID, err := userService.RegisterUser(testUser)

	// Проверяем, что вернулась ошибка базы данных и ID не сгенерирован
	assert.Error(t, err)
	assert.Equal(t, "", userID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	userService := service.NewUserService(mockRepo, mockSessionRepo)

	// Мокаем успешную проверку пароля
	mockRepo.On("CheckUserPasswordHash", "johndoe", mock.Anything).Return("123", nil)

	// Мокаем успешное создание сессии
	mockSessionRepo.On("CreateSession", "123", mock.Anything, mock.Anything).Return(nil)

	token, err := userService.Login("johndoe", "correctpassword")

	// Проверяем, что ошибок нет и токен сгенерирован
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestUserService_CreateSession_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	userService := service.NewUserService(mockRepo, mockSessionRepo)

	// Мокаем успешную проверку пароля
	mockRepo.On("CheckUserPasswordHash", "johndoe", mock.Anything).Return("123", nil)

	// Мокаем ошибку при создании сессии
	mockSessionRepo.On("CreateSession", "123", mock.Anything, mock.Anything).Return(errors.New("db error"))

	token, err := userService.Login("johndoe", "correctpassword")

	// Проверяем, что вернулась ошибка создания сессии и токен не был сгенерирован
	assert.Error(t, err)
	assert.Equal(t, "", token)

	mockRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestUserService_Login_Failed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	// Мокаем ошибку аутентификации
	mockRepo.On("CheckUserPasswordHash", "johndoe", mock.Anything).Return("", errors.New("Auth error"))

	token, err := userService.Login("johndoe", "wrongpassword")

	// Проверяем, что вернулась ошибка и токен не был сгенерирован
	assert.Error(t, err)
	assert.Equal(t, "", token)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	testUser := domain.User{
		ID:         "123",
		FirstName:  "John",
		SecondName: "Doe",
		Username:   "johndoe",
	}

	// Мокаем успешное получение пользователя
	mockRepo.On("GetUserByID", "123").Return(testUser, nil)

	user, err := userService.GetById("123")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe", user.Username)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	// Мокаем ошибку получения пользователя
	mockRepo.On("GetUserByID", "123").Return(domain.User{}, errors.New("user not found"))

	user, err := userService.GetById("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	assert.Error(t, err)
	assert.Equal(t, domain.User{}, user)

	mockRepo.AssertExpectations(t)
}

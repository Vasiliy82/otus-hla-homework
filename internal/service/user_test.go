package service_test

import (
	"database/sql"
	"errors"
	errors_ "errors"
	"testing"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/service"
	"github.com/lib/pq"
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

func (m *MockUserRepository) GetByID(id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}
func (m *MockUserRepository) GetByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockSessionRepository) CreateSession(userID, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)

}

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)
	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
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
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем ошибку дублирования логина
	mockRepo.On("RegisterUser", testUser).Return("", &pq.Error{Code: "23505"})

	userID, err := userService.RegisterUser(testUser)

	// Проверяем, что вернулась ошибка конфликта и ID не сгенерирован
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 409, apperr.Code)
		assert.Equal(t, "Login already used", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}

	assert.Equal(t, "", userID)

	mockRepo.AssertExpectations(t)
}

func TestUserService_RegisterUser_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)
	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем ошибку базы данных
	mockRepo.On("RegisterUser", testUser).Return("", apperrors.NewInternalServerError("db error", errors_.New("db error")))

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

	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(testUser, nil)
	// Мокаем успешное создание сессии
	mockSessionRepo.On("CreateSession", "123", mock.Anything, mock.Anything).Return(nil)

	token, err := userService.Login("johndoe@gmail.com", "correctpassword")

	// Проверяем, что ошибок нет и токен сгенерирован
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestUserService_Login_GetByUsername_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	userService := service.NewUserService(mockRepo, mockSessionRepo)

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(domain.User{}, errors.New("database error"))

	_, err := userService.Login("johndoe@gmail.com", "correctpassword")

	// Проверяем, что вернулась ошибка создания сессии и токен не был сгенерирован
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 500, apperr.Code)
		assert.Equal(t, "Internal server error", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_CreateSession_DBError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	userService := service.NewUserService(mockRepo, mockSessionRepo)

	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(testUser, nil)
	mockSessionRepo.On("CreateSession", "123", mock.Anything, mock.Anything).Return(errors.New("database error"))

	token, err := userService.Login("johndoe@gmail.com", "correctpassword")

	// Проверяем, что вернулась ошибка создания сессии и токен не был сгенерирован
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 500, apperr.Code)
		assert.Equal(t, "Internal server error", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}
	assert.Equal(t, "", token)
	mockRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestUserService_Login_Failed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(testUser, nil)

	token, err := userService.Login("johndoe@gmail.com", "wrongpassword")

	// Проверяем, что вернулась ошибка и токен не был сгенерирован
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 401, apperr.Code)
		assert.Equal(t, "Wrong password", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}
	assert.Equal(t, "", token)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем успешное получение пользователя
	mockRepo.On("GetByID", "123").Return(testUser, nil)

	user, err := userService.GetById("123")

	// Проверяем, что ошибок нет и пользователь получен
	assert.NoError(t, err)
	assert.Equal(t, "johndoe@gmail.com", user.Username)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo, nil)

	// Мокаем ошибку получения пользователя
	mockRepo.On("GetByID", "123").Return(domain.User{}, sql.ErrNoRows)

	user, err := userService.GetById("123")

	// Проверяем, что вернулась ошибка и пользователь не был найден
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 404, apperr.Code)
		assert.Equal(t, "User not found", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}
	assert.Equal(t, domain.User{}, user)

	mockRepo.AssertExpectations(t)
}

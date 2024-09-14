package user_test

import (
	"database/sql"
	"errors"
	errors_ "errors"
	"testing"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/domain/mocks"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	user "github.com/Vasiliy82/otus-hla-homework/internal/services/user"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_RegisterUser_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)
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
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)
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
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)
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
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

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
	// mockSess.On("CreateSession", "123", mock.Anything).Return("cf2df81324491ec70e6cd2e922b06929", nil)

	mockJwt.On("GenerateToken", "123", mock.Anything).Return(domain.TokenString("jwt_token"), nil)

	token, err := userService.Login("johndoe@gmail.com", "correctpassword")

	// Проверяем, что ошибок нет и токен сгенерирован
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_GetByUsername_DBError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(domain.User{}, errors.New("database error"))

	token, err := userService.Login("johndoe@gmail.com", "correctpassword")

	assert.Equal(t, domain.TokenString(""), token)

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
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

	testUser := domain.User{
		ID:           "123",
		FirstName:    "John",
		SecondName:   "Doe",
		Username:     "johndoe@gmail.com",
		PasswordHash: "e6cd2e922b06929cf2df81324491ec70",
	}

	// Мокаем
	mockRepo.On("GetByUsername", "johndoe@gmail.com").Return(testUser, nil)
	// mockSess.On("CreateSession", "123", mock.Anything, mock.Anything).Return("", errors.New("database error"))
	mockJwt.On("GenerateToken", "123", mock.Anything).Return(domain.TokenString(""), errors.New("database error"))

	token, err := userService.Login("johndoe@gmail.com", "correctpassword")

	// Проверяем, что вернулась ошибка создания сессии и токен не был сгенерирован
	var apperr *apperrors.AppError
	if errors.As(err, &apperr) {
		assert.Equal(t, 500, apperr.Code)
		assert.Equal(t, "Internal server error", apperr.Message)
	} else {
		t.Fatalf("expected error of type *apperrors.AppError, got %T", err)
	}
	assert.Equal(t, domain.TokenString(""), token)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Failed(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

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
	assert.Equal(t, domain.TokenString(""), token)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetById_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

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
	mockRepo := mocks.NewUserRepository(t)
	mockJwt := mocks.NewJWTService(t)
	userService := user.NewUserService(mockRepo, mockJwt)

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

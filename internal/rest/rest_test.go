package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/Vasiliy82/otus-hla-homework/internal/rest"
)

// Mock для UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) Login(username, password string) (string, string, error) {
	args := m.Called(username, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUserService) GetById(id string) (domain.User, error) {
	return domain.User{}, nil
}

// 1. Тест на успешный ответ при регистрации
func TestUserHandler_RegisterUser_Success(t *testing.T) {
	// Создаем инстанс Echo
	e := echo.New()

	// Создаем mock для UserService
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		t.Fatal(err)
	}

	// Тестовые данные для успешной регистрации
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}
	mockUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Birthdate:    birthdate,
		Biography:    "Blah-blah-blah",
		City:         "Silent Hill",
		Username:     "johndoe@gmail.com",
		PasswordHash: "482c811da5d5b4bc6d497ffa98491e38",
	}

	// Мокаем успешную регистрацию
	mockUserService.On("RegisterUser", mockUser).Return("bb49a7d7-3e85-4935-9afd-570ec8ea318b", nil)

	// Формируем HTTP-запрос
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON)) // .WithContext(context.Background())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	err = handler.RegisterUser(c)

	// Проверяем, что ошибок нет, а статус-код ответа 200 OK
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем ответ
	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "bb49a7d7-3e85-4935-9afd-570ec8ea318b", resp["user_id"])

	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoPassword(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (нет поля Password)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoCity(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (нет поля City)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoBirthDate(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	// Формируем HTTP-запрос с неполными данными (нет поля BirthDate)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_WrongBirthDate(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-13-01"

	// Формируем HTTP-запрос с неполными данными (нет поля Password)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoSecondName(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (нет поля SecondName)
	reqBody := dto.RegisterUserRequest{
		FirstName: "John",
		Birthdate: birthdateStr,
		Biography: "Blah-blah-blah",
		City:      "Silent Hill",
		Username:  "johndoe@gmail.com",
		Password:  "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoFirstName(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (нет поля FirstName)
	reqBody := dto.RegisterUserRequest{
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_NoUsername(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (нет поля Username)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 2. Тест на ошибочный запрос (например, неполные данные)
func TestUserHandler_RegisterUser_BadRequest_WrongUsername(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"

	// Формируем HTTP-запрос с неполными данными (Username)
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe",
		Password:   "password123",
	}

	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.RegisterUser(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUserService.AssertExpectations(t)
}

// RegisterUser - duplicate user
func TestUserHandler_RegisterUser_Duplicate(t *testing.T) {
	// Создаем инстанс Echo
	e := echo.New()

	// Создаем mock для UserService
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	birthdateStr := "2020-01-01"
	birthdate, err := time.Parse("2006-01-02", birthdateStr)
	if err != nil {
		t.Fatal(err)
	}

	// Тестовые данные для успешной регистрации
	reqBody := dto.RegisterUserRequest{
		FirstName:  "John",
		SecondName: "Doe",
		Birthdate:  birthdateStr,
		Biography:  "Blah-blah-blah",
		City:       "Silent Hill",
		Username:   "johndoe@gmail.com",
		Password:   "password123",
	}
	mockUser := domain.User{
		FirstName:    "John",
		SecondName:   "Doe",
		Birthdate:    birthdate,
		Biography:    "Blah-blah-blah",
		City:         "Silent Hill",
		Username:     "johndoe@gmail.com",
		PasswordHash: "482c811da5d5b4bc6d497ffa98491e38",
	}
	// Мокаем успешную регистрацию
	mockUserService.On("RegisterUser", mockUser).Return("", apperrors.NewConflictError("Login already used"))

	// Формируем HTTP-запрос
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	err = handler.RegisterUser(c)

	// Проверяем, что вернулась ошибка и статус-код 400 Bad Request
	assert.Equal(t, http.StatusConflict, rec.Code)
	mockUserService.AssertExpectations(t)
}

// 3. Тест на успешный ответ при аутентификации
func TestUserHandler_Login_Success(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	// Тестовые данные для успешного логина
	reqBody := dto.LoginRequest{
		Username: "johndoe@gmail.com",
		Password: "password123",
	}

	// Мокаем успешный логин
	mockUserService.On("Login", "johndoe@gmail.com", "password123").Return("bb49a7d7-3e85-4935-9afd-570ec8ea318b", "token123", nil)

	// Формируем HTTP-запрос
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	err := handler.Login(c)

	// Проверяем, что ошибок нет и статус-код 200 OK
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверяем ответ
	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "bb49a7d7-3e85-4935-9afd-570ec8ea318b", resp["user_id"])
	assert.Equal(t, "token123", resp["token"])

	mockUserService.AssertExpectations(t)
}

// 4. Тест на ошибку при неверных данных (неверный пароль)
func TestUserHandler_Login_Failure(t *testing.T) {
	e := echo.New()
	mockUserService := new(MockUserService)
	handler := rest.NewUserHandler(mockUserService)

	// Тестовые данные для логина с неправильным паролем
	reqBody := dto.LoginRequest{
		Username: "johndoe@gmail.com",
		Password: "wrongpassword",
	}

	// Мокаем ошибку логина
	mockUserService.On("Login", "johndoe@gmail.com", "wrongpassword").Return("", "", errors.New("invalid credentials"))

	// Формируем HTTP-запрос
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Вызываем обработчик
	if err := handler.Login(c); err != nil {
		t.Fatal(err)
	}

	// Проверяем, что вернулась ошибка и статус-код 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

package validators

import (
	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/apperrors"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
)

var (
	validate *validator.Validate
	tran     ut.Translator
)

// Инициализация валидатора и переводчика
func init() {
	validate = validator.New()

	validate.RegisterValidation("sex", sexValidator)

	// Установка переводчика для русского языка
	ruLocale := ru.New()
	uni := ut.New(ruLocale, ruLocale)

	// Получение переводчика для русского языка
	tran, _ = uni.GetTranslator("ru")

	// Регистрация русских переводов ошибок валидации
	ruTranslations.RegisterDefaultTranslations(validate, tran)
}

func ValidateRegisterUserRequest(request dto.RegisterUserRequest) error {
	if err := validate.Struct(request); err != nil {
		return apperrors.NewValidationError("Validation error", err, tran)
	}
	return nil
}

func ValidateLoginRequest(request dto.LoginRequest) error {
	if err := validate.Struct(request); err != nil {
		return apperrors.NewValidationError("Validation error", err, tran)
	}
	return nil
}

func ValidateUserId(id domain.UserKey) error {
	if id == "" {
		return apperrors.NewBadRequestError("id cannot be empty")
	}
	return nil
}

// Функция валидации биолгоического пола
func sexValidator(fl validator.FieldLevel) bool {
	sex := fl.Field().String()
	return sex == "M" || sex == "F"
}

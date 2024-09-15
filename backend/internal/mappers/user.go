package mappers

import (
	"time"

	"github.com/Vasiliy82/otus-hla-homework/domain"
	"github.com/Vasiliy82/otus-hla-homework/internal/dto"
)

func ToUser(request dto.RegisterUserRequest) (domain.User, error) {
	birthdate, err := time.Parse("2006-01-02", request.Birthdate)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		FirstName:  request.FirstName,
		SecondName: request.SecondName,
		Biography:  request.Biography,
		City:       request.City,
		Birthdate:  birthdate,
		Username:   request.Username,
	}
	user.SetPassword(request.Password)

	return user, nil
}

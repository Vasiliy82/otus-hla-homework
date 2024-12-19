package mappers

import (
	"fmt"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/dto"
)

func ToSex(sex string) (domain.Sex, error) {
	if sex == "M" {
		return domain.Male, nil
	}
	if sex == "F" {
		return domain.Female, nil
	}
	return domain.Male, fmt.Errorf("error converting %s into domain.Sex", sex)
}

func ToUser(request dto.RegisterUserRequest) (domain.User, error) {
	birthdate, err := time.Parse("2006-01-02", request.Birthdate)
	if err != nil {
		return domain.User{}, err
	}
	sex, err := ToSex(request.Sex)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Sex:       sex,
		Biography: request.Biography,
		City:      request.City,
		Birthdate: birthdate,
		Username:  request.Username,
	}
	user.SetPassword(request.Password)

	return user, nil
}

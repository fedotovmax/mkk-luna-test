package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/validation"
)

type CreateUser struct {
	UserName string `json:"username" validate:"required" example:"username123"`
	Email    string `json:"email" validate:"required" example:"testemail@mail.ru"`
	Password string `json:"password" validate:"required" example:"Assdfsdf2323!_"`
}

func (i *CreateUser) SetPasswordHash(hash string) {
	i.Password = hash
}

func (i *CreateUser) Validate() error {

	var verrs []domain.ValidationError

	err := validation.IsEmail(i.Email)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "email",
			Message: "Некорректный формат email",
		})
	}

	msg, err := validatePassword(i.Password)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "password",
			Message: msg,
		})
	}

	err = validation.MinLength(i.UserName, 3)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "username",
			Message: "Имя пользователя должно быть не менее 3 символов",
		})
	}

	if len(verrs) > 0 {
		ve := &domain.ValidatationErrors{
			Errors: verrs,
		}
		return ve
	}

	return nil
}

type Login struct {
	Email    string `json:"email" validate:"required" example:"testemail@mail.ru"`
	Password string `json:"password" validate:"required" example:"Assdfsdf2323!_"`
}

func (i *Login) Validate() error {

	var verrs []domain.ValidationError

	err := validation.IsEmail(i.Email)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "email",
			Message: "Некорректный формат email",
		})
	}

	msg, err := validatePassword(i.Password)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "password",
			Message: msg,
		})
	}

	if len(verrs) > 0 {
		ve := &domain.ValidatationErrors{
			Errors: verrs,
		}
		return ve
	}

	return nil
}

package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/validation"
)

type CreateTeam struct {
	Name string `json:"name" validate:"required" example:"team 1"`
}

func (i *CreateTeam) Validate() error {
	var verrs []domain.ValidationError

	err := validation.MinLength(i.Name, 3)
	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "Name",
			Message: "Название команды должно быть не менее 3 символов",
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

type InviteMember struct {
	UserID string `json:"user_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

func (i *InviteMember) Validate() error {
	var verrs []domain.ValidationError

	_, err := validation.IsUUID(i.UserID)
	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "user_id",
			Message: "Некорректный формат user_id",
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

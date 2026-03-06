package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/validation"
)

type UUID struct {
	ID string `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

func (i *UUID) Validate(field string) error {
	var verrs []domain.ValidationError

	_, err := validation.IsUUID(i.ID)

	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   field,
			Message: "Некорректный формат id",
		})
	}

	if len(verrs) > 0 {
		return &domain.ValidatationErrors{
			Errors: verrs,
		}
	}

	return nil
}

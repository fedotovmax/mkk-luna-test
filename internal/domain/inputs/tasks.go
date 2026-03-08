package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/validation"
)

type CreateTask struct {
	AssigneeID  *string `json:"assignee_id" validate:"optional" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	Description *string `json:"description" validate:"optional" example:"Description of the task"`
	TeamID      string  `json:"team_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	Title       string  `json:"title" validate:"required" example:"Task title"`
}

func (i *CreateTask) Validate() error {
	var verrs []domain.ValidationError

	if _, err := validation.IsUUID(i.TeamID); err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "team_id",
			Message: "Некорректный формат team_id",
		})
	}

	if err := validation.MinLength(i.Title, 3); err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "title",
			Message: "Название задачи должно быть не менее 3 символов",
		})
	}

	if i.AssigneeID != nil {
		if _, err := validation.IsUUID(*i.AssigneeID); err != nil {
			verrs = append(verrs, domain.ValidationError{
				Field:   "assignee_id",
				Message: "Некорректный формат assignee_id",
			})
		}
	}

	if i.Description != nil {
		if err := validation.LengthRange(*i.Description, 3, 1000); err != nil {
			verrs = append(verrs, domain.ValidationError{
				Field:   "description",
				Message: "Описание должно быть не менее 3 символов и не более 1000 символов",
			})
		}
	}

	if len(verrs) > 0 {
		return &domain.ValidatationErrors{
			Errors: verrs,
		}
	}

	return nil
}

type UpdateTask struct {
	Status      *domain.Status `json:"status" validate:"optional"`
	Title       *string        `json:"title" validate:"required" example:"Task title"`
	Description *string        `json:"description" validate:"optional" example:"Description of the task"`
}

func (i *UpdateTask) Validate() error {
	var verrs []domain.ValidationError

	if i.Status != nil && !i.Status.IsValid() {
		verrs = append(verrs, domain.ValidationError{
			Field:   "status",
			Message: "Некорректное значение статуса",
		})
	}

	if i.Title != nil {
		if err := validation.MinLength(*i.Title, 1); err != nil {
			verrs = append(verrs, domain.ValidationError{
				Field:   "title",
				Message: "Название задачи не может быть пустым",
			})
		}
	}

	if i.Description != nil {
		if err := validation.LengthRange(*i.Description, 3, 1000); err != nil {
			verrs = append(verrs, domain.ValidationError{
				Field:   "description",
				Message: "Описание должно быть не менее 3 символов и не более 1000 символов",
			})
		}
	}

	if len(verrs) > 0 {
		return &domain.ValidatationErrors{
			Errors: verrs,
		}
	}

	return nil
}

type FindManyTasks struct {
	Status     domain.Status
	TeamID     string
	AssigneeID string
}

func (i *FindManyTasks) Validate() error {

	var verrs []domain.ValidationError

	_, err := validation.IsUUID(i.TeamID)
	if err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "team_id",
			Message: "Некорректный формат team_id",
		})
	}

	if i.AssigneeID != "" {
		_, err := validation.IsUUID(i.AssigneeID)
		if err != nil {
			verrs = append(verrs, domain.ValidationError{
				Field:   "assignee_id",
				Message: "Некорректный формат assignee_id",
			})
		}
	}

	if i.Status != "" {
		if !i.Status.IsValid() {
			verrs = append(verrs, domain.ValidationError{
				Field:   "status",
				Message: "Некорректный статус задачи",
			})
		}
	}

	if len(verrs) > 0 {
		return &domain.ValidatationErrors{
			Errors: verrs,
		}
	}

	return nil
}

type CreateHistory struct {
	OldTask     *domain.Task
	ChangedByID string
}

type CreateComment struct {
	Text string
}

func (i *CreateComment) Validate() error {
	var verrs []domain.ValidationError

	if err := validation.LengthRange(i.Text, 1, 500); err != nil {
		verrs = append(verrs, domain.ValidationError{
			Field:   "text",
			Message: "Текст комментария должен быть от 1 до 500 символов",
		})
	}

	if len(verrs) > 0 {
		return &domain.ValidatationErrors{
			Errors: verrs,
		}
	}

	return nil
}

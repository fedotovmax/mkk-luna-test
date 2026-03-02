package domain

type IDResponse struct {
	ID string `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
}

package ports

import (
	"context"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

type TaskStorage interface {
	FindOne(ctx context.Context, id string) (*domain.Task, error)
	FindMany(ctx context.Context, offset, limit int, in *inputs.FindManyTasks) (*domain.FindTasksResponse, error)
	FindTaskHistory(ctx context.Context, taskID string) ([]*domain.History, error)
	FindTaskComments(ctx context.Context, taskID string) ([]*domain.Comment, error)

	Create(ctx context.Context, ownerID string, in *inputs.CreateTask) (string, error)
	CreateHistory(ctx context.Context, in *inputs.CreateHistory) (string, error)
	CreateComment(ctx context.Context, userID, taskID, text string) (string, error)

	Update(ctx context.Context, id string, in *inputs.UpdateTask) error
}

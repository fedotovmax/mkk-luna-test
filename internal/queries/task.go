package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
)

type Tasks interface {
	FindByID(ctx context.Context, id string) (*domain.Task, error)
	FindMany(ctx context.Context, limit int, offset int, in *inputs.FindManyTasks) (*domain.FindTasksResponse, error)
	FindHistory(ctx context.Context, id string) ([]*domain.History, error)
	FindComments(ctx context.Context, id string) ([]*domain.Comment, error)
}

type tasks struct {
	storage ports.TaskStorage
}

func NewTasks(storage ports.TaskStorage) Tasks {
	return &tasks{storage: storage}
}

func (q *tasks) FindMany(ctx context.Context, limit int, offset int, in *inputs.FindManyTasks) (
	*domain.FindTasksResponse, error) {

	res, err := q.storage.FindMany(ctx, offset, limit, in)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (q *tasks) FindByID(ctx context.Context, id string) (*domain.Task, error) {

	t, err := q.storage.FindOne(ctx, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrTaskNotFound, err)
		}
		return nil, err
	}

	return t, nil
}

func (q *tasks) FindHistory(ctx context.Context, id string) ([]*domain.History, error) {
	return q.storage.FindTaskHistory(ctx, id)
}

func (q *tasks) FindComments(ctx context.Context, id string) ([]*domain.Comment, error) {
	return q.storage.FindTaskComments(ctx, id)
}

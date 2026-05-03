package ports

import (
	"context"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type TodoRepository interface {
	FindAll(ctx context.Context) ([]domain.Todo, error)
	Create(ctx context.Context, todo domain.Todo) (domain.Todo, error)
	MarkCompleted(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

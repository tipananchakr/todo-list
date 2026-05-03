package ports

import (
	"context"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type TodoRepository interface {
	FindAllByUser(ctx context.Context, userID string) ([]domain.Todo, error)
	Create(ctx context.Context, todo domain.Todo) (domain.Todo, error)
	MarkCompleted(ctx context.Context, userID string, id string) error
	Delete(ctx context.Context, userID string, id string) error
}

package ports

import (
	"context"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

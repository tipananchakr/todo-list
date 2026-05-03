package application

import (
	"context"
	"strings"
	"time"

	"github.com/tipananchakr/todo-list/internal/core/domain"
	"github.com/tipananchakr/todo-list/internal/core/ports"
)

type TodoService struct {
	repository ports.TodoRepository
	timeout    time.Duration
}

func NewTodoService(repository ports.TodoRepository) *TodoService {
	return &TodoService{
		repository: repository,
		timeout:    5 * time.Second,
	}
}

func (s *TodoService) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repository.FindAll(ctx)
}

func (s *TodoService) CreateTodo(ctx context.Context, body string) (domain.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	todo := domain.Todo{
		Body:      strings.TrimSpace(body),
		Completed: false,
	}

	return s.repository.Create(ctx, todo)
}

func (s *TodoService) CompleteTodo(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repository.MarkCompleted(ctx, id)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repository.Delete(ctx, id)
}

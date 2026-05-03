package application

import (
	"context"
	"testing"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

func TestTodoServiceCreateTodoTrimsBodyAndDefaultsIncomplete(t *testing.T) {
	repository := &fakeTodoRepository{}
	service := NewTodoService(repository)

	todo, err := service.CreateTodo(context.Background(), "user-1", "  Learn testing  ")
	if err != nil {
		t.Fatalf("CreateTodo returned error: %v", err)
	}

	if todo.Body != "Learn testing" {
		t.Fatalf("expected trimmed body, got %q", todo.Body)
	}

	if todo.Completed {
		t.Fatal("expected created todo to be incomplete")
	}

	if todo.UserID != "user-1" {
		t.Fatalf("expected todo user id, got %q", todo.UserID)
	}
}

func TestTodoServiceDelegatesCompleteTodo(t *testing.T) {
	repository := &fakeTodoRepository{}
	service := NewTodoService(repository)

	if err := service.CompleteTodo(context.Background(), "user-1", "6636d3d046b1b2dd3b46e001"); err != nil {
		t.Fatalf("CompleteTodo returned error: %v", err)
	}

	if repository.completedID != "6636d3d046b1b2dd3b46e001" {
		t.Fatalf("expected repository to complete todo id, got %q", repository.completedID)
	}

	if repository.completedUserID != "user-1" {
		t.Fatalf("expected repository to complete todo for user, got %q", repository.completedUserID)
	}
}

type fakeTodoRepository struct {
	completedID     string
	completedUserID string
}

func (r *fakeTodoRepository) FindAllByUser(ctx context.Context, userID string) ([]domain.Todo, error) {
	return []domain.Todo{}, nil
}

func (r *fakeTodoRepository) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	todo.ID = "6636d3d046b1b2dd3b46e001"
	return todo, nil
}

func (r *fakeTodoRepository) MarkCompleted(ctx context.Context, userID string, id string) error {
	r.completedUserID = userID
	r.completedID = id
	return nil
}

func (r *fakeTodoRepository) Delete(ctx context.Context, userID string, id string) error {
	return nil
}

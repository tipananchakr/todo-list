package application

import (
	"context"
	"strings"
	"testing"

	"github.com/tipananchakr/todo-list/internal/core/domain"
)

func TestAuthServiceRegisterCreatesUserAndToken(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, fakePasswordHasher{}, fakeTokenManager{})

	result, err := service.Register(context.Background(), " USER@example.com ", "secret123")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if result.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", result.User.Email)
	}

	if result.Token != "token-"+result.User.ID {
		t.Fatalf("expected token for user, got %q", result.Token)
	}
}

func TestAuthServiceRegisterRejectsExistingEmail(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, fakePasswordHasher{}, fakeTokenManager{})

	if _, err := service.Register(context.Background(), "user@example.com", "secret123"); err != domain.ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, fakePasswordHasher{}, fakeTokenManager{})

	if _, err := service.Login(context.Background(), "user@example.com", "wrong-password"); err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceCurrentUser(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, fakePasswordHasher{}, fakeTokenManager{})

	user, err := service.CurrentUser(context.Background(), "token-user-1")
	if err != nil {
		t.Fatalf("CurrentUser returned error: %v", err)
	}

	if user.ID != "user-1" {
		t.Fatalf("expected user-1, got %q", user.ID)
	}
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		users: map[string]domain.User{
			"user-1": {
				ID:           "user-1",
				Email:        "user@example.com",
				PasswordHash: "hashed-secret123",
			},
		},
	}
}

type fakeUserRepository struct {
	users map[string]domain.User
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *fakeUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = "user-2"
	r.users[user.ID] = user
	return user, nil
}

type fakePasswordHasher struct{}

func (h fakePasswordHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (h fakePasswordHasher) Compare(hash string, password string) error {
	if hash != "hashed-"+password {
		return domain.ErrInvalidCredentials
	}

	return nil
}

type fakeTokenManager struct{}

func (m fakeTokenManager) Generate(userID string) (string, error) {
	return "token-" + userID, nil
}

func (m fakeTokenManager) Validate(token string) (string, error) {
	userID, ok := strings.CutPrefix(token, "token-")
	if !ok || userID == "" {
		return "", domain.ErrInvalidToken
	}

	return userID, nil
}

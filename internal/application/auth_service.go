package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tipananchakr/todo-list/internal/core/domain"
	"github.com/tipananchakr/todo-list/internal/core/ports"
)

type AuthService struct {
	users   ports.UserRepository
	hasher  ports.PasswordHasher
	tokens  ports.TokenManager
	timeout time.Duration
}

type AuthResult struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

func NewAuthService(users ports.UserRepository, hasher ports.PasswordHasher, tokens ports.TokenManager) *AuthService {
	return &AuthService{
		users:   users,
		hasher:  hasher,
		tokens:  tokens,
		timeout: 5 * time.Second,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	email = normalizeEmail(email)
	if email == "" || len(password) < 6 {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return AuthResult{}, domain.ErrEmailAlreadyExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return AuthResult{}, err
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.Create(ctx, domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		return AuthResult{}, err
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.users.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return AuthResult{}, domain.ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{User: user, Token: token}, nil
}

func (s *AuthService) CurrentUser(ctx context.Context, token string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	userID, err := s.tokens.Validate(token)
	if err != nil {
		return domain.User{}, domain.ErrInvalidToken
	}

	return s.users.FindByID(ctx, userID)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

package http

import (
	"context"
	"encoding/json"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/tipananchakr/todo-list/internal/application"
	"github.com/tipananchakr/todo-list/internal/core/domain"
)

func TestHealthRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "GET", "/", "")
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
	}
}

func TestSwaggerUIRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "GET", "/swagger", "")
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
	}

	if !strings.Contains(string(body), "SwaggerUIBundle") {
		t.Fatal("expected swagger ui html")
	}

	if !strings.Contains(string(body), "/docs/openapi.yaml") {
		t.Fatal("expected swagger ui to load openapi document")
	}
}

func TestOpenAPIRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "GET", "/docs/openapi.yaml", "")
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
	}

	if !strings.Contains(string(body), "openapi: 3.0.3") {
		t.Fatal("expected openapi yaml")
	}
}

func TestCreateTodoRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "POST", "/api/todos", `{"body":" Write tests "}`)
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, response.StatusCode)
	}

	var todo domain.Todo
	if err := json.NewDecoder(response.Body).Decode(&todo); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if todo.Body != "Write tests" {
		t.Fatalf("expected trimmed body, got %q", todo.Body)
	}
}

func TestCompleteTodoRouteReturnsBadRequestForInvalidID(t *testing.T) {
	app := newTestApp(&testTodoRepository{markCompletedError: domain.ErrInvalidTodoID}, newTestAuthService())

	response := testRequest(t, app, "PATCH", "/api/todos/invalid-id", "")
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, response.StatusCode)
	}
}

func TestRegisterRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "POST", "/api/auth/register", `{"email":"Test@Example.com","password":"secret123"}`)
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, response.StatusCode)
	}

	var result authResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if result.User.Email != "test@example.com" {
		t.Fatalf("expected normalized email, got %q", result.User.Email)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}
}

func TestLoginRoute(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "POST", "/api/auth/login", `{"email":"test@example.com","password":"secret123"}`)
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
	}
}

func TestMeRouteRequiresBearerToken(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	response := testRequest(t, app, "GET", "/api/auth/me", "")
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", fiber.StatusUnauthorized, response.StatusCode)
	}
}

func TestMeRouteReturnsCurrentUser(t *testing.T) {
	app := newTestApp(&testTodoRepository{}, newTestAuthService())

	request := httptest.NewRequest("GET", "/api/auth/me", nil)
	request.Header.Set(fiber.HeaderAuthorization, "Bearer token-user-1")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
	}
}

func newTestApp(repository *testTodoRepository, authService *application.AuthService) *fiber.App {
	app := fiber.New()
	RegisterRoutes(app, Services{
		Auth: authService,
		Todo: application.NewTodoService(repository),
	})
	return app
}

func testRequest(t *testing.T, app *fiber.App, method string, target string, body string) *nethttp.Response {
	t.Helper()

	request := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	}

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("test request: %v", err)
	}

	return response
}

type testTodoRepository struct {
	markCompletedError error
}

func (r *testTodoRepository) FindAll(ctx context.Context) ([]domain.Todo, error) {
	return []domain.Todo{
		{ID: "6636d3d046b1b2dd3b46e001", Body: "Write docs", Completed: false},
	}, nil
}

func (r *testTodoRepository) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	todo.ID = "6636d3d046b1b2dd3b46e002"
	return todo, nil
}

func (r *testTodoRepository) MarkCompleted(ctx context.Context, id string) error {
	return r.markCompletedError
}

func (r *testTodoRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func newTestAuthService() *application.AuthService {
	users := &testUserRepository{
		users: map[string]domain.User{
			"user-1": {
				ID:           "user-1",
				Email:        "test@example.com",
				PasswordHash: "hashed-secret123",
			},
		},
	}

	return application.NewAuthService(users, testPasswordHasher{}, testTokenManager{})
}

type testUserRepository struct {
	users map[string]domain.User
}

func (r *testUserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *testUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrUserNotFound
}

func (r *testUserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	user.ID = "user-2"
	r.users[user.ID] = user
	return user, nil
}

type testPasswordHasher struct{}

func (h testPasswordHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (h testPasswordHasher) Compare(hash string, password string) error {
	if hash != "hashed-"+password {
		return domain.ErrInvalidCredentials
	}

	return nil
}

type testTokenManager struct{}

func (m testTokenManager) Generate(userID string) (string, error) {
	return "token-" + userID, nil
}

func (m testTokenManager) Validate(token string) (string, error) {
	userID, ok := strings.CutPrefix(token, "token-")
	if !ok || userID == "" {
		return "", domain.ErrInvalidToken
	}

	return userID, nil
}

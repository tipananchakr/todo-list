package http

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/tipananchakr/todo-list/internal/application"
	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type AuthHandler struct {
	service *application.AuthService
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

func registerAuthRoutes(router fiber.Router, service *application.AuthService) {
	handler := AuthHandler{service: service}

	router.Post("/register", handler.register)
	router.Post("/login", handler.login)
	router.Get("/me", handler.requireAuth, handler.me)
}

func (h AuthHandler) register(c *fiber.Ctx) error {
	var request authRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid body"})
	}

	result, err := h.service.Register(c.Context(), request.Email, request.Password)
	if err != nil {
		return authError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(authResponse{User: result.User, Token: result.Token})
}

func (h AuthHandler) login(c *fiber.Ctx) error {
	var request authRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid body"})
	}

	result, err := h.service.Login(c.Context(), request.Email, request.Password)
	if err != nil {
		return authError(c, err)
	}

	return c.JSON(authResponse{User: result.User, Token: result.Token})
}

func (h AuthHandler) requireAuth(c *fiber.Ctx) error {
	token := bearerToken(c.Get(fiber.HeaderAuthorization))
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{Error: "missing bearer token"})
	}

	user, err := h.service.CurrentUser(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{Error: "invalid bearer token"})
	}

	c.Locals("user", user)
	return c.Next()
}

func (h AuthHandler) me(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(domain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{Error: "invalid bearer token"})
	}

	return c.JSON(user)
}

func authError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse{Error: "invalid credentials"})
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(errorResponse{Error: "email already exists"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "auth service error"})
	}
}

func bearerToken(header string) string {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}

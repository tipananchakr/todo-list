package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/tipananchakr/todo-list/internal/application"
	"github.com/tipananchakr/todo-list/internal/core/domain"
)

type TodoHandler struct {
	service *application.TodoService
}

type createTodoRequest struct {
	Body string `json:"body"`
}

func registerTodoRoutes(router fiber.Router, service *application.TodoService) {
	handler := TodoHandler{service: service}

	router.Get("/", handler.getTodos)
	router.Post("/", handler.createTodo)
	router.Patch("/:id", handler.completeTodo)
	router.Delete("/:id", handler.deleteTodo)
}

func (h TodoHandler) getTodos(c *fiber.Ctx) error {
	todos, err := h.service.GetTodos(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "cannot fetch todos"})
	}

	if todos == nil {
		todos = []domain.Todo{}
	}

	return c.JSON(todos)
}

func (h TodoHandler) createTodo(c *fiber.Ctx) error {
	var request createTodoRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid body"})
	}

	todo, err := h.service.CreateTodo(c.Context(), request.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "cannot insert todo"})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func (h TodoHandler) completeTodo(c *fiber.Ctx) error {
	if err := h.service.CompleteTodo(c.Context(), c.Params("id")); err != nil {
		if errors.Is(err, domain.ErrInvalidTodoID) {
			return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid todo id"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "cannot update todo"})
	}

	return c.JSON(successResponse{Success: true})
}

func (h TodoHandler) deleteTodo(c *fiber.Ctx) error {
	if err := h.service.DeleteTodo(c.Context(), c.Params("id")); err != nil {
		if errors.Is(err, domain.ErrInvalidTodoID) {
			return c.Status(fiber.StatusBadRequest).JSON(errorResponse{Error: "invalid todo id"})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "cannot delete todo"})
	}

	return c.JSON(successResponse{Success: true})
}

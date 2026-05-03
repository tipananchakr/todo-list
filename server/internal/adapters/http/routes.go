package http

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"github.com/tipananchakr/todo-list/internal/application"
)

type Services struct {
	Auth *application.AuthService
	Todo *application.TodoService
}

type successResponse struct {
	Success bool `json:"success"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(app *fiber.App, services Services) {
	app.Get("/", health)
	app.Get("/docs/openapi.yaml", openapi)
	app.Get("/swagger", swagger)
	app.Get("/swagger/", swagger)

	api := app.Group("/api")
	registerAuthRoutes(api.Group("/auth"), services.Auth)
	registerTodoRoutes(api.Group("/todos"), services.Todo, services.Auth)
}

func health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func openapi(c *fiber.Ctx) error {
	path, err := findOpenAPIPath()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{Error: "cannot load openapi document"})
	}

	c.Set(fiber.HeaderContentType, "application/yaml; charset=utf-8")
	return c.SendFile(path)
}

func swagger(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(swaggerHTML)
}

func findOpenAPIPath() (string, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(workingDirectory, "docs", "openapi.yaml")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(workingDirectory)
		if parent == workingDirectory {
			return "", os.ErrNotExist
		}
		workingDirectory = parent
	}
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Todo List API Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body { margin: 0; background: #fafafa; }
    .swagger-ui .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/docs/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/tipananchakr/todo-list/internal/adapters/http"
	mongorepo "github.com/tipananchakr/todo-list/internal/adapters/mongodb"
	"github.com/tipananchakr/todo-list/internal/adapters/security"
	"github.com/tipananchakr/todo-list/internal/application"
	"github.com/tipananchakr/todo-list/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	mongoClient, err := mongorepo.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)

	database := mongoClient.Database(cfg.DatabaseName)

	todoRepository := mongorepo.NewTodoRepository(database.Collection(cfg.TodoCollection))
	userRepository := mongorepo.NewUserRepository(database.Collection(cfg.UserCollection))
	if err := userRepository.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}

	todoService := application.NewTodoService(todoRepository)
	authService := application.NewAuthService(
		userRepository,
		security.NewBcryptPasswordHasher(),
		security.NewHMACTokenManager(cfg.JWTSecret, 24*time.Hour),
	)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	http.RegisterRoutes(app, http.Services{
		Auth: authService,
		Todo: todoService,
	})

	log.Fatal(app.Listen("0.0.0.0:" + cfg.Port))
}

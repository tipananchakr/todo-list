# Todo List API

Todo List API is a Go REST service built with Fiber and MongoDB. The project uses hexagonal architecture and feature-based adapters so new features such as auth, todos, notifications, or payments can be added without mixing HTTP, database, and business logic in one place.

## Features

- Health check endpoint
- User registration with bcrypt password hashing
- User login with HMAC signed JWT access tokens
- Authenticated current-user endpoint
- List all todos
- Create a todo
- Mark a todo as completed
- Delete a todo
- Swagger UI at `/swagger`
- OpenAPI YAML at `/docs/openapi.yaml`

## Architecture

```text
.
├── main.go                         # Application composition root
├── internal
│   ├── application                 # Feature use cases / application services
│   │   ├── auth_service.go
│   │   └── todo_service.go
│   ├── core
│   │   ├── domain                  # Business entities and domain errors
│   │   └── ports                   # Interfaces required by the application core
│   └── adapters
│       ├── http                    # Fiber route registration and feature handlers
│       ├── mongodb                 # MongoDB repositories
│       └── security                # Password hashing and JWT token adapter
└── docs
    └── openapi.yaml                # Swagger/OpenAPI document
```

The dependency direction is:

```text
HTTP / MongoDB / Security adapters -> application services -> core ports and domain
```

To add a new feature, add its domain model and port in `internal/core`, its use case in `internal/application`, and its external implementations in `internal/adapters`.

## Requirements

- Go 1.25.4 or compatible
- MongoDB database or MongoDB Atlas cluster

## Environment

Create a `.env` file in the project root:

```env
PORT=3000
MONGODB_URI=mongodb+srv://username:password@cluster.example.mongodb.net/golang_db?appName=Cluster0
MONGODB_DATABASE=golang_db
MONGODB_TODO_COLLECTION=todos
MONGODB_USER_COLLECTION=users
JWT_SECRET=replace-with-a-long-random-secret
```

Defaults:

- `PORT`: `3000`
- `MONGODB_DATABASE`: `golang_db`
- `MONGODB_TODO_COLLECTION`: `todos`
- `MONGODB_USER_COLLECTION`: `users`
- `JWT_SECRET`: `local-dev-secret`

## Run

Install dependencies:

```bash
go mod download
```

Start the API:

```bash
go run main.go
```

Or use make:

```bash
make run
```

The API runs on `http://localhost:3000` by default.

## Test

Run all tests:

```bash
go test ./...
```

Or use make:

```bash
make test
```

## Swagger

Start the API and open Swagger UI in your browser:

```text
http://localhost:3000/swagger
```

The Swagger UI loads the OpenAPI document from:

```text
http://localhost:3000/docs/openapi.yaml
```

The OpenAPI document is also available in the repository at:

```text
docs/openapi.yaml
```
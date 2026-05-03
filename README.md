# Todo List

Todo List is a full-stack app with a React/Vite client and a Go REST API built with Fiber and MongoDB. The server uses hexagonal architecture and feature-based adapters so auth and todo behavior stay separated from HTTP, database, and security implementations.

## Structure

```text
.
├── client        # React + TypeScript + Vite frontend
├── server        # Go Fiber API, MongoDB repositories, OpenAPI docs
├── makefile      # Root commands for server and client workflows
└── README.md
```

## App Features

- Health check endpoint
- User registration with bcrypt password hashing
- User login with HMAC signed JWT access tokens
- Authenticated current-user endpoint
- List user todos
- Create user todos
- Mark a todo as completed
- Delete a todo
- React client with login, register, logout, and todo CRUD
- Server Swagger UI at `/swagger`
- Server OpenAPI YAML at `/docs/openapi.yaml`

## Server Architecture

```text
server
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

To add a new server feature, add its domain model and port in `server/internal/core`, its use case in `server/internal/application`, and its external implementations in `server/internal/adapters`.

## Requirements

- Go 1.25.4 or compatible
- MongoDB database or MongoDB Atlas cluster
- Node.js 20 or compatible
- npm

## Server Environment

Create `server/.env`:

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

Install server dependencies:

```bash
cd server
go mod download
cd ..
```

Start the API from the project root:

```bash
make run
```

Install client dependencies:

```bash
cd client
npm install
cd ..
```

Start the client from the project root:

```bash
make client-dev
```

Defaults:

- API: `http://localhost:3000`
- Client: `http://localhost:5173`

## Test

Run server tests:

```bash
make test
```

Run client checks:

```bash
make client-lint
make client-build
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
server/docs/openapi.yaml
```

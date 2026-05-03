SERVER_DIR := server
GOCACHE ?= $(CURDIR)/.cache/go-build

run:
	cd $(SERVER_DIR) && GOCACHE=$(GOCACHE) go run main.go

test:
	cd $(SERVER_DIR) && GOCACHE=$(GOCACHE) go test ./...

client-dev:
	cd client && npm run dev

client-build:
	cd client && npm run build

client-lint:
	cd client && npm run lint

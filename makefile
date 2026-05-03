GOCACHE ?= $(CURDIR)/.cache/go-build

run:
	GOCACHE=$(GOCACHE) go run main.go

test:
	GOCACHE=$(GOCACHE) go test ./...

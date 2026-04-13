APP_NAME=task-api

.PHONY: run build tidy fmt test

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

tidy:
	go mod tidy

fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

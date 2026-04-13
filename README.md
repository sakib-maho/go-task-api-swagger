# Golang Beego RestAPI Swagger (Rebuilt)

This repository now contains a working Go REST API starter project with Swagger/OpenAPI documentation and task CRUD endpoints.

## What Is Included

- HTTP server with graceful shutdown
- REST endpoints (`/api/v1/tasks`)
- In-memory datastore with thread-safe access
- Input validation and consistent JSON error responses
- OpenAPI spec (`docs/openapi.json`)
- Swagger UI page (`/swagger`)
- Makefile for common developer tasks

## Tech Stack

- Go 1.22
- [chi](https://github.com/go-chi/chi) router
- [uuid](https://github.com/google/uuid) for IDs
- OpenAPI 3.0 (static JSON spec)

## API Endpoints

- `GET /health`
- `GET /api/v1/tasks`
- `POST /api/v1/tasks`
- `GET /api/v1/tasks/{taskID}`
- `PUT /api/v1/tasks/{taskID}`
- `DELETE /api/v1/tasks/{taskID}`

Task status values:

- `todo`
- `in_progress`
- `done`

## Quick Start

```bash
git clone https://github.com/sakib-maho/Golang-Beggo-RestAPI-Swagger.git
cd Golang-Beggo-RestAPI-Swagger
cp .env.example .env
go mod tidy
go run ./cmd/server
```

Server starts on `http://localhost:8080` by default.

## Swagger / OpenAPI

- Swagger page: `http://localhost:8080/swagger`
- OpenAPI file: `http://localhost:8080/swagger/openapi.json`

## Example Requests

Create task:

```bash
curl -X POST "http://localhost:8080/api/v1/tasks" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Write API docs",
    "description": "Document all endpoints",
    "status": "todo"
  }'
```

List tasks:

```bash
curl "http://localhost:8080/api/v1/tasks"
```

## Project Structure

```text
.
├── cmd/server/main.go
├── internal/
│   ├── api/
│   ├── config/
│   ├── model/
│   └── store/
├── docs/
│   ├── openapi.json
│   └── swagger.html
├── Makefile
└── go.mod
```

## Notes

- This implementation uses an in-memory store for fast local development.
- Data resets whenever the server restarts.
- You can replace `internal/store` with PostgreSQL implementation later.

## License

MIT - see [LICENSE](LICENSE).

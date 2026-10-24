# AI Knowledge Hub

A technical documentation assistant built in Go, using RAG (Retrieval-Augmented Generation).

The project is currently in its first API milestone. It exposes a health check and
an in-memory document collection while the Go foundations are developed.

## Run the API

Go 1.26.5 or later is required.

```bash
cd dev/backend
go run ./cmd/server
```

The server listens on `http://localhost:8081`.

```bash
curl http://localhost:8081/health

curl -X POST http://localhost:8081/documents \
  -H 'Content-Type: application/json' \
  -d '{"title":"Go","content":"Go is a programming language."}'

curl http://localhost:8081/documents
```

Documents are kept in memory and are lost when the server stops.

## Run the tests

```bash
cd dev/backend
go test -race ./...
```

## Learning examples

The standalone examples for Go fundamentals live separately from the API:

```bash
cd dev/backend
go run ./cmd/learn
```

The full project plan is documented in [dev/SPEC.md](dev/SPEC.md).

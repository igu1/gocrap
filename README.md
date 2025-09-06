# gocrap

A simple HTTP service that accepts a JSON-defined browser automation "Flow" and executes it using Playwright for Go.

## Features

- Define sequences of steps: `go_to`, `extract`, `extract_multi`, `save`
- Uses Playwright (Chromium) to navigate and scrape
- Returns the in-memory results map as JSON

## Project Layout

- `cmd/server` — main entrypoint for the HTTP server
- `internal/flow` — Flow types and execution logic
- `internal/server` — HTTP handlers and route registration

## Getting Started

1. Install Go 1.23+
2. Install Playwright browsers (first run downloads automatically)

### Run

```bash
make run
# or
go run ./cmd/server
```

Server runs on http://localhost:8080

### Example Request

```bash
curl -X POST http://localhost:8080/run \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Example",
    "url": "https://example.org",
    "path": [
      {"action": "go_to", "target": "/"},
      {"action": "extract", "selector": "h1", "store_as": "title"}
    ]
  }'
```

## Build

```bash
make build
```

## Docker

```bash
make docker-build
make docker-run
```

## License

MIT © 2025 Eesa Bin Zakariyya

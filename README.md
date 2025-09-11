# Gocrap — A Tiny HTTP Browser-Automation Server (with Big Energy)

A compact HTTP service that accepts a JSON-defined browser automation "Flow" and runs it using Playwright for Go. Think of it as your browser butler: it clicks, fills, scrolls, extracts, and saves — politely and predictably. It won’t make coffee, but it can grab cookies. Literally.

## Why Gocrap?

- **Declarative flows**: Describe what to do in JSON; the server does the rest.
- **Playwright-powered**: Reliable, headless (or headful) automation via Chromium.
- **Fast iteration**: POST a flow and immediately get results as JSON.
- **Just enough features**: A pragmatic set of actions for scraping, testing, and small automation.

> Automation that doesn’t talk back (unless there’s a 500).

---

## Project Layout

- `cmd/server` — Main entrypoint for the HTTP server
- `internal/core` — Core flow types, step validation, handler registry
- `internal/actions` — All available flow actions (see catalog below)
- `internal/server` — HTTP handlers and route registration

---

## Getting Started

1. Install Go 1.23+
2. First run will automatically download Playwright browsers

### Run

```bash
make run
# or
go run ./cmd/server
```

Server runs on http://localhost:8080

---

## HTTP API

### POST /run
Execute a flow.

- **Request**: JSON body with a `url` and a `path` (list of steps)
- **Response**: JSON of the in-memory result map (what your steps stored)

```bash
curl -X POST http://localhost:8080/run \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Example",
    "url": "https://example.org",
    "path": [
      {"action": "go_to", "target": "/"},
      {"action": "extract", "selector": "h1", "store_as": "title"},
      {"action": "save", "filename": "scraped_output.json"}
    ]
  }'
```

---

## Flow Model

Each step is an object with the following fields. Only some are required per action.

- **action**: Name of the action (see catalog)
- **target**: CSS selector, URL path, or other target string (action-specific)
- **selector**: CSS selector used by extraction actions
- **child_selector**: Optional secondary selector (currently unused by built-ins)
- **attribute**: Attribute to extract instead of text (e.g., `href`)
- **store_as**: Key name to store extracted data under
- **filename**: File path for saving content, screenshots, or cookies
- **value**: Free-form value (e.g., text to fill, key to press, or special modes)
- **duration**: Milliseconds to wait
- **description**: A human note you can add for readability

---

## Actions Catalog

Below is the complete set of built-in actions with required fields and quick examples.

### 1) go_to
- **Purpose**: Navigate to a path under the base `url`.
- **Requires**: `target`
- **Example**:
```json
{"action": "go_to", "target": "/search"}
```

### 2) extract
- **Purpose**: Extract text or an attribute from the first element matching `selector`.
- **Requires**: `selector`, `store_as`
- **Optional**: `attribute` (e.g., `href`)
- **Example**:
```json
{"action": "extract", "selector": "h1", "store_as": "title"}
```

### 3) extract_multi
- **Purpose**: Extract text/attribute from all elements matching `selector` into an array.
- **Requires**: `selector`, `store_as`
- **Optional**: `attribute`
- **Example**:
```json
{"action": "extract_multi", "selector": ".result a", "attribute": "href", "store_as": "links"}
```

### 4) save
- **Purpose**: Save the current in-memory map to disk as pretty JSON.
- **Requires**: `filename`
- **Example**:
```json
{"action": "save", "filename": "scraped_output.json"}
```

### 5) wait
- **Purpose**: Pause the flow for a number of milliseconds.
- **Requires**: `duration`
- **Example**:
```json
{"action": "wait", "duration": 1500}
```

### 6) click
- **Purpose**: Click an element.
- **Requires**: `target` (CSS selector)
- **Example**:
```json
{"action": "click", "target": "#submit"}
```

### 7) fill
- **Purpose**: Fill an input/textarea with text.
- **Requires**: `target`, `value`
- **Example**:
```json
{"action": "fill", "target": "input[name=q]", "value": "golang"}
```

### 8) press
- **Purpose**: Send a key to an element (e.g., `Enter`).
- **Requires**: `target`, `value` (key)
- **Example**:
```json
{"action": "press", "target": "input[name=q]", "value": "Enter"}
```

### 9) check
- **Purpose**: Check a checkbox.
- **Requires**: `target`
- **Example**:
```json
{"action": "check", "target": "input[type=checkbox]#agree"}
```

### 10) uncheck
- **Purpose**: Uncheck a checkbox.
- **Requires**: `target`
- **Example**:
```json
{"action": "uncheck", "target": "#subscribe"}
```

### 11) scroll
- **Purpose**: Scroll the page.
- **Requires**: `value` — one of: `"top"`, `"bottom"`, or `"x,y"` coordinates (e.g., `"0,1200"`).
- **Example**:
```json
{"action": "scroll", "value": "bottom"}
```

### 12) reload
- **Purpose**: Reload the current page.
- **Requires**: none
- **Example**:
```json
{"action": "reload"}
```

### 13) back
- **Purpose**: Go back in browser history.
- **Requires**: none
- **Example**:
```json
{"action": "back"}
```

### 14) forward
- **Purpose**: Go forward in browser history.
- **Requires**: none
- **Example**:
```json
{"action": "forward"}
```

### 15) cookies
- **Purpose**: Save or load cookies to/from a file.
- **Requires**: `value` ("save" | "load"), `filename`
- **Examples**:
```json
{"action": "cookies", "value": "save", "filename": "extracted_cookies.json"}
```
```json
{"action": "cookies", "value": "load", "filename": "extracted_cookies.json"}
```

### 16) screenshot
- **Purpose**: Save a full-page screenshot.
- **Requires**: `filename`
- **Example**:
```json
{"action": "screenshot", "filename": "page.png"}
```

### 17) eval
- **Purpose**: Execute a JavaScript snippet in the page context.
- **Requires**: `target` (JS string)
- **Example**:
```json
{"action": "eval", "target": "document.title"}
```

---

## A More Complete Example

```json
{
  "title": "News Links",
  "url": "https://example.org",
  "path": [
    {"action": "go_to", "target": "/"},
    {"action": "wait", "duration": 1000},
    {"action": "extract", "selector": "h1", "store_as": "hero_title"},
    {"action": "click", "target": "a[href='/archive']"},
    {"action": "wait", "duration": 500},
    {"action": "extract_multi", "selector": ".article a", "attribute": "href", "store_as": "article_links"},
    {"action": "screenshot", "filename": "archive.png"},
    {"action": "save", "filename": "articles.json"}
  ]
}
```

---

## Tips & Notes

- Most actions auto-wait for elements where it makes sense; still, **`wait`** is your friend for flaky sites.
- Use **`extract_multi`** for lists; single **`extract`** folds into a string automatically.
- Use **`store_as`** to shape the output map; **`save`** persists it to disk.
- Cookies are powerful for authenticated flows: **`cookies: save`** after login, then **`cookies: load`** before later steps.
- Remember: with great power comes great responsibility — and sometimes rate limits.

---

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

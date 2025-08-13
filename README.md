# Exchange Rate Service

A simple HTTP service providing latest and historical foreign exchange rates, plus currency conversion. Rates are fetched from the Exchangerate-API v6 and cached in-memory. A background scheduler pre-warms the cache.

## Features
- Latest rate lookup for supported base/target currencies
- Historical rate lookup for the last 90 days
- Currency conversion (latest or on a specific date within 90 days)
- In-memory caching with TTL (default 1h)
- Background scheduler to refresh and warm cache (default every 1h)
- Dockerized build using a static distroless runtime image

## Requirements
- Go 1.22+
- An Exchangerate-API key (v6)

## Configuration
The service reads configuration from a YAML file named `config.yaml` in the current working directory.

Example `config.yaml`:
```yaml
server:
  port: 8080

exchange_rate:
  api_key: "YOUR-EXCHANGERATE-API-KEY"
```

Notes:
- The server listens on the configured `server.port`.
- The Exchangerate-API key is required. The service will fail without it.
- When running in Docker, place or mount `config.yaml` at `/app/config.yaml` inside the container (see Docker section).

## Supported Currencies
Defined in `internal/domain/exchange.go`:
- USD, INR, EUR, JPY, GBP

Requests including other currencies will be rejected.

## API
Base URL: `http://localhost:<port>`

### GET /latest
Get the latest exchange rate from one currency to another.

Query parameters:
- `from` (required): base currency (e.g., `USD`)
- `to` (required): target currency (e.g., `INR`)

Response:
```json
{
  "from": "USD",
  "to": "INR",
  "rate": 83.12
}
```

Example:
```bash
curl "http://localhost:8080/latest?from=USD&to=INR"
```

### GET /convert
Convert an amount between currencies using the latest rate or a specific date within the last 90 days.

Query parameters:
- `from` (required): base currency
- `to` (required): target currency
- `amount` (required): amount to convert (non-negative number)
- `date` (optional): `YYYY-MM-DD`; if provided, must be within the last 90 days

Response:
```json
{
  "from": "USD",
  "to": "EUR",
  "amount": 10,
  "converted": 9.14
}
```

Examples:
```bash
# Latest conversion
curl "http://localhost:8080/convert?from=USD&to=EUR&amount=10"

# Historical conversion
curl "http://localhost:8080/convert?from=USD&to=EUR&amount=10&date=2024-06-01"
```

### GET /historical
Get the exchange rate for a specific date within the last 90 days.

Query parameters:
- `from` (required): base currency
- `to` (required): target currency
- `date` (required): `YYYY-MM-DD`; must be within the last 90 days

Response:
```json
{
  "from": "GBP",
  "to": "JPY",
  "rate": 200.45
}
```

Example:
```bash
curl "http://localhost:8080/historical?from=GBP&to=JPY&date=2024-06-01"
```

### Errors
- 400 Bad Request: missing or invalid parameters, unsupported currency, or date outside of last 90 days
- 500 Internal Server Error: upstream API errors or unexpected conditions

## Running Locally
1. Create `config.yaml` as above.
2. Start the server:
```bash
go run ./cmd/server
```
The server will log: `Server starting on port <port>..`

### Build
```bash
go build -o bin/server ./cmd/server
./bin/server
```

### Test
```bash
go test ./...
```

## Docker
The repository includes a multi-stage Dockerfile that builds a static binary and runs it in a distroless image.

### Build Image
```bash
docker build -t exchange-rate-service:latest .
```

### Run Container
- Expose the configured port
- Provide `config.yaml` at `/app/config.yaml` inside the container

```bash
# Assuming your local config.yaml sets port 8080
docker run --rm \
  -p 8080:8080 \
  -v "$(pwd)/config.yaml:/app/config.yaml:ro" \
  --name exchange-rate-service \
  exchange-rate-service:latest
```

## Architecture Overview
- `cmd/server/main.go`: bootstraps configuration, cache, service, scheduler, and HTTP routes
- `internal/config`: loads `config.yaml` into an app-wide config struct
- `internal/domain`: domain primitives (supported currencies)
- `internal/infrastructure`:
  - `api_client.go`: HTTP client for Exchangerate-API (latest and historical endpoints)
  - `cache.go`: in-memory rates cache with TTL
- `internal/service`:
  - `exchange_service.go`: business logic, input validation, cache usage, and API orchestration
  - `scheduler.go`: background cache warm-up for all supported base currencies at a fixed interval
- `internal/handler`:
  - `exchange_handler.go`: HTTP handlers and route registration (`/latest`, `/convert`, `/historical`)
- `internal/utils`: small helpers (date/time, string checks)

## Behavior and Limits
- Rates are cached per base currency for 1 hour (configurable only via code at construction time in `main.go`).
- Historical lookups are limited to the last 90 days.
- Only a small set of currencies is supported by default; extend `internal/domain/exchange.go` as needed.
- Network timeouts to the upstream API are set to 5 seconds.

## Upstream Data Source
- Exchangerate-API v6 (`https://v6.exchangerate-api.com`) — obtain an API key from `https://www.exchangerate-api.com/`.

## License
MIT (or your preferred license).

# --- Builder stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Cache module downloads
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# --- Runtime stage ---
FROM gcr.io/distroless/static:nonroot

# Set workdir where config.yaml should be mounted/copied at runtime
WORKDIR /app

# Copy compiled binary
COPY --from=builder /out/server /server

# Run as non-root user
USER nonroot:nonroot

# Optionally expose a port if your config.yaml uses one, e.g., 8080
# EXPOSE 8080

# Expect a config.yaml present in /app at runtime (mount or bake into image)
ENTRYPOINT ["/server"]

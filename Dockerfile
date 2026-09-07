# ─── Base Stage ────────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS base
WORKDIR /app

# Install dependencies only
COPY go.mod go.sum ./
RUN go mod download

# ─── Development Stage ─────────────────────────────────────────────────────────
FROM base AS dev
# Install make, air for hot reload, and migrate for database migrations
RUN apk add --no-cache make && \
    go install github.com/air-verse/air@latest && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# ─── Builder Stage ─────────────────────────────────────────────────────────────
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api/main.go

# ─── Production Stage ──────────────────────────────────────────────────────────
FROM alpine:3.20 AS prod
WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezone support
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./api"]

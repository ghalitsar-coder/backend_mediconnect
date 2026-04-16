FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies (no gcc needed since we use pure go sqlite)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Explicitly set CGO_ENABLED=0 since we are using glebarez pure-go sqlite
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Minimalist runtime
FROM alpine:latest

WORKDIR /app

# Setup timezone
RUN apk --no-cache add ca-certificates tzdata

# Copy built binary
COPY --from=builder /app/main .

# Copy environment template if needed or config structure
COPY .env .env

# Expose port
EXPOSE 8080

CMD ["./main"]

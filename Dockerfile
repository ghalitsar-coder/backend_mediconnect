# Stage 1: Build the Go application
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files to download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the entire project 
COPY . .

# Build the executable
RUN CGO_ENABLED=0 GOOS=linux go build -o mediconnect-backend .

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/mediconnect-backend .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./mediconnect-backend"]

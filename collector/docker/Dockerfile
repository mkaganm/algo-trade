# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o collector ./cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/collector .
COPY .env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./collector"]
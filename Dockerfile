# Build stage: use the official Go image.
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Install git (required for downloading Go module dependencies).
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum files to download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the project files.
COPY . .

# Build the Go binary.
RUN go build -o pr-review-bot

# Final stage: use a minimal image.
FROM alpine:latest

# Install tzdata to support time zones.
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from the builder stage.
COPY --from=builder /app/pr-review-bot .

# Optionally, copy the .env file if you want it baked into the image.
COPY .env .

CMD ["./pr-review-bot"]

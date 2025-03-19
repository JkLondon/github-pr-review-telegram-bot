# Build stage: use the official Go image.
FROM golang:1.18-alpine AS builder

# Set the working directory inside the container.
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

WORKDIR /app

# Copy the binary from the builder stage.
COPY --from=builder /app/pr-review-bot .

# Optionally, copy the .env file if you want it baked into the image.
# If you prefer to mount it externally, you can remove this line.
COPY .env .

# Command to run the bot.
CMD ["./pr-review-bot"]

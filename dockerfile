# Start from an image that includes both Go and Node.js/npm
FROM node:latest AS npm_builder

# Set the working directory for npm
WORKDIR /app

# Copy npm files
COPY package-lock.json package.json postcss.config.js tailwind.config.js ./

# Install npm dependencies
RUN npm install

# Start from the official Golang image
FROM golang:1.22 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

WORKDIR /app

# Copy the source code into the container
COPY . .

# Copy npm dependencies from the npm_builder stage
COPY --from=npm_builder /app/node_modules ./node_modules

# Build the Go Application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/web

# Start a new stage from scratch
# This stage will contain the compiled Go binary only
FROM alpine:latest

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/tls ./tls

# Expose the port the application runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

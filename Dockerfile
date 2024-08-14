# Build stage using the official Golang image
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application code to the container
COPY . .

# Build the Go application statically (important for Distroless)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapp

# Final stage: Distroless container
FROM gcr.io/distroless/base-debian11

# Set working directory inside the container
WORKDIR /app

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/goapp .

# Set environment variables (can be overridden at runtime)
ENV CRON_SCHEDULE="*/30 * * * * *"
ENV API_URL="https://google.com"

# Command to run the Go application
ENTRYPOINT ["./goapp"]

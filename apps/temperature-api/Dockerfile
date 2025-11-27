# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /apps

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o temperature-api

# Runtime stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /apps

# Copy the binary from the builder stage
COPY --from=builder /apps/temperature-api .

# Expose the port the app runs on
EXPOSE 8081

# Command to run the application
CMD ["./temperature-api"]
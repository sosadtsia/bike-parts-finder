# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk --no-cache add ca-certificates git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /scraper ./cmd/scraper

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the built executable from the build stage
COPY --from=builder /scraper .

# Use the non-root user
USER appuser

# Expose health check port
EXPOSE 8081

# Run the application
CMD ["./scraper"]

FROM golang:1.25.2-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api

# Final stage
FROM alpine:3.22

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api /app/api

# Copy web assets
COPY --from=builder /app/web/frontend/dist /app/web/dist

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/api"]

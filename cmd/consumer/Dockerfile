FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o /app/consumer ./cmd/consumer

# Final stage
FROM alpine:3.22

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/consumer /app/consumer

# Run the application
CMD ["/app/consumer"]

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy root module files first
COPY go.mod go.sum ./
COPY ./models ./models
COPY ./database ./database

# Copy service specific files
COPY ./services/auth/go.mod ./services/auth/go.sum ./
COPY ./services/auth .

# Download dependencies for root module first
RUN go mod download

# Build the application
RUN go build -o auth-service main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/auth-service .

# Expose port
EXPOSE 8084

# Run the binary
CMD ["./auth-service"]
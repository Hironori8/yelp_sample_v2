FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy shared models
COPY ./shared/models ./shared/models

# Copy service specific files
COPY ./services/business/go.mod ./services/business/go.sum ./
COPY ./services/business .

# Download dependencies for root module first
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]
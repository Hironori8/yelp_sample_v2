FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy service specific go.mod and go.sum
COPY ./services/logging/go.mod ./services/logging/go.sum ./

# Copy source code
COPY ./services/logging .

RUN go mod download

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8083

CMD ["./main"]
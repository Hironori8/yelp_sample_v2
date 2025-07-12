FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy service specific go.mod and go.sum
COPY ./services/gateway/go.mod ./services/gateway/go.sum ./

# Copy source code
COPY ./services/gateway .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]